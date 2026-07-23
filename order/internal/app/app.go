package app

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderv1API "github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1"
	"github.com/krapagen/my_microservices_rocket/order/internal/config"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/logger"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
)

const (
	// HTTP таймауты
	httpReadHeaderTimeout = 5 * time.Second  // Защита от Slowloris атаки
	httpReadTimeout       = 15 * time.Second // Лимит на чтение всего запроса
	httpWriteTimeout      = 15 * time.Second // Лимит на запись ответа
	httpIdleTimeout       = 60 * time.Second // Таймаут keep-alive соединений
	httpShutdownTimeout   = 5 * time.Second  // Таймаут graceful shutdown

	// gRPC клиент keepalive параметры
	grpcKeepaliveTime    = 5 * time.Minute // Интервал ping'ов для обнаружения мёртвого сервера
	grpcKeepaliveTimeout = 1 * time.Second // Таймаут ожидания pong
)

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов
type App struct {
	diContainer *diContainer
	httpServer  *http.Server
	listener    net.Listener
}

// New создаёт и инициализирует приложение
func New(ctx context.Context) *App {
	a := &App{}

	a.initDeps(ctx)

	return a
}

func initGRPCClient(address, service string) *grpc.ClientConn {
	// Создать gRPC соединение
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                grpcKeepaliveTime,
			Timeout:             grpcKeepaliveTimeout,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		slog.Error("не удалось подключиться к "+service, "error", err)
		os.Exit(1)
	}
	return conn
}

// Run управляет жизненным циклом приложения: запускает http-сервер,
// обрабатывает сигналы ОС и выполняет graceful shutdown
//
// Сервер запускается в отдельной горутине, а main-горутина синхронно ждёт
// либо сигнал SIGINT/SIGTERM, либо падение сервера. После этого
// closer.CloseAll вызывается синхронно — main-горутина гарантированно
// дожидается завершения всех закрытий перед выходом из Run
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runHTTPServer()
	}()

	var runErr error
	select {
	case runErr = <-errCh:
		// сервер сам упал (например, bind: address already in use)
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}
	cancel() // снимаем перехват сигналов, повторный Ctrl+C завершит процесс принудительно

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer shutdownCancel()

	if err := closer.CloseAll(shutdownCtx); err != nil {
		slog.Error("ошибка при завершении работы", "error", err)
		if runErr == nil {
			runErr = err
		}
	}

	return runErr
}

// initDeps последовательно инициализирует все зависимости приложения
func (a *App) initDeps(ctx context.Context) {
	inits := []func(context.Context){
		a.initDI,
		a.initLogger,
		a.initListener,
		a.initInventoryClient,
		a.initPaymentClient,
		a.initHTTPServer,
	}

	for _, f := range inits {
		f(ctx)
	}
}

// initDI создаёт DI-контейнер
func (a *App) initDI(_ context.Context) {
	a.diContainer = &diContainer{}
}

// initLogger настраивает глобальный slog с уровнем из конфига
func (a *App) initLogger(_ context.Context) {
	logger.Init(config.AppConfig().Logger.Level)
}

// initListener создаёт TCP-листенер для http-сервера
func (a *App) initListener(_ context.Context) {
	listener, err := net.Listen("tcp", config.AppConfig().HTTP.Address()) //nolint:noctx // net.Listen не требует контекст, адрес из конфига
	if err != nil {
		slog.Error("не удалось создать TCP-листенер", "error", err)
		os.Exit(1)
	}

	a.listener = listener
}

func (a *App) initInventoryClient(_ context.Context) {
	// Создать gRPC соединение с InventoryService
	a.diContainer.inventoryConn = initGRPCClient(config.AppConfig().InventoryClient.Address, "InventoryService")
}

func (a *App) initPaymentClient(_ context.Context) {
	// Создать gRPC соединение с PaymentService
	a.diContainer.paymentConn = initGRPCClient(config.AppConfig().PaymentClient.Address, "PaymentService")
}

// initHTTPServer создаёт и настраивает HTTP-сервер, регистрирует обработчики
func (a *App) initHTTPServer(ctx context.Context) {
	// Получаем API-обработчик до регистрации closer'а: ленивая инициализация
	// зацепит за собой создание пула БД и зарегистрирует его в closer'е.
	// Closer работает по LIFO, поэтому пул должен попасть туда раньше http-сервера —
	// тогда при shutdown сначала остановится приём запросов, а уже потом закроется БД
	api := a.diContainer.OrderV1API(ctx)
	// Создаём HTTP обработчик с новой архитектурой
	server, err := orderv1.NewServer(api, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
	if err != nil {
		slog.Error("ошибка создания зависимостей приложения", "error", err)
		return
	}
	a.httpServer = &http.Server{
		Handler:           server,
		ReadHeaderTimeout: httpReadHeaderTimeout, // Защита от Slowloris атаки
		ReadTimeout:       httpReadTimeout,       // Лимит на чтение всего запроса
		WriteTimeout:      httpWriteTimeout,      // Лимит на запись ответа
		IdleTimeout:       httpIdleTimeout,       // Таймаут keep-alive соединений
	}
	closer.Add("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
}

// runHTTPServer запускает HTTP-сервер и блокирует до его остановки
func (a *App) runHTTPServer() error {
	slog.Info("🚀 http-сервер запущен", "address", config.AppConfig().HTTP.Address())

	return a.httpServer.Serve(a.listener)
}
