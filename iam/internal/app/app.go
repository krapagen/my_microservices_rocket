package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	iamv1API "github.com/krapagen/my_microservices_rocket/iam/internal/api/http/iam/v1"
	"github.com/krapagen/my_microservices_rocket/iam/internal/config"
	"github.com/krapagen/my_microservices_rocket/iam/internal/interceptor"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/grpc/health"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/logger"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/ratelimit"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/tracing"
	iamv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/iam/v1"
	iamauthv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
	iamuserv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/user/v1"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
	shutdownTimeout           = 5 * time.Second

	httpReadHeaderTimeout = 5 * time.Second  // Защита от Slowloris атаки
	httpReadTimeout       = 15 * time.Second // Лимит на чтение всего запроса
	httpWriteTimeout      = 15 * time.Second // Лимит на запись ответа
	httpIdleTimeout       = 60 * time.Second // Таймаут keep-alive соединений
)

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов
type App struct {
	diContainer  *diContainer
	grpcServer   *grpc.Server
	httpServer   *http.Server
	listener     net.Listener
	httpListener net.Listener
}

// New создаёт и инициализирует приложение
func New(ctx context.Context) *App {
	a := &App{}

	a.initDeps(ctx)

	return a
}

// Run управляет жизненным циклом приложения: запускает gRPC-сервер,
// обрабатывает сигналы ОС и выполняет graceful shutdown
//
// Сервер запускается в отдельной горутине, а main-горутина синхронно ждёт
// либо сигнал SIGINT/SIGTERM, либо падение сервера. После этого
// closer.CloseAll вызывается синхронно — main-горутина гарантированно
// дожидается завершения всех закрытий перед выходом из Run
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 2)
	go func() {
		errCh <- a.runGRPCServer()
	}()
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

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
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
		a.initTracer,
		a.initListener,
		a.initGRPCServer,
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
	logger.Init(logger.Config{
		Level:             config.AppConfig().Logger.Level,
		ServiceName:       config.AppConfig().OtelConfig.ServiceName,
		Environment:       config.AppConfig().Stage,
		EnableOTLP:        true,
		CollectorEndpoint: config.AppConfig().OtelConfig.Endpoint,
	})
	closer.Add("logger", func(ctx context.Context) error {
		return logger.Close()
	})
}

// initTracer создаёт OTel TracerProvider с OTLP-экспортером и регистрирует его глобально.
// Вместе с ним устанавливается propagator (W3C TraceContext + Baggage) —
// без него trace context не передавался бы между сервисами через gRPC metadata
func (a *App) initTracer(ctx context.Context) {
	shutdown, err := tracing.InitTracer(ctx, tracing.Config{
		CollectorEndpoint: config.AppConfig().OtelConfig.Endpoint,
		ServiceName:       config.AppConfig().OtelConfig.ServiceName,
		Environment:       config.AppConfig().Stage,
	})
	if err != nil {
		slog.Error("не удалось инициализировать трейсер", "error", err)
		return
	}
	closer.Add("tracer", func(ctx context.Context) error {
		return shutdown(ctx)
	})
}

// initListener создаёт TCP-листенер для gRPC-сервера
func (a *App) initListener(_ context.Context) {
	listener, err := net.Listen("tcp", config.AppConfig().GRPC.Address()) //nolint:noctx // net.Listen не требует контекст, адрес из конфига
	if err != nil {
		slog.Error("не удалось создать TCP-листенер", "error", err)
		os.Exit(1)
	}

	a.listener = listener

	// HTTP-листенер для REST API (register/login/logout/getUser)
	httpListener, err := net.Listen("tcp", config.AppConfig().HTTP.Address()) //nolint:noctx // net.Listen не требует контекст, адрес из конфига
	if err != nil {
		slog.Error("не удалось создать HTTP-листенер", "error", err)
		os.Exit(1)
	}

	a.httpListener = httpListener
}

// initGRPCServer создаёт и настраивает gRPC-сервер, регистрирует обработчики
func (a *App) initGRPCServer(ctx context.Context) {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: false,
		}),
		grpc.UnaryInterceptor(interceptor.ErrorInterceptor),
	)

	// Получаем API-обработчик до регистрации closer'а: ленивая инициализация
	// зацепит за собой создание пула БД и зарегистрирует его в closer'е.
	// Closer работает по LIFO, поэтому пул должен попасть туда раньше gRPC-сервера —
	// тогда при shutdown сначала остановится приём запросов, а уже потом закроется БД
	apiAuth := a.diContainer.IamAuthV1API(ctx)
	apiUser := a.diContainer.IamUserV1API(ctx)

	closer.Add("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	iamuserv1.RegisterUserServiceServer(a.grpcServer, apiUser)
	iamauthv1.RegisterAuthServiceServer(a.grpcServer, apiAuth)
}

// runGRPCServer запускает gRPC-сервер и блокирует до его остановки
func (a *App) runGRPCServer() error {
	slog.Info("🚀 gRPC-сервер запущен", "address", config.AppConfig().GRPC.Address())

	return a.grpcServer.Serve(a.listener)
}

// initHTTPServer создаёт и настраивает HTTP-сервер, регистрирует обработчики
func (a *App) initHTTPServer(ctx context.Context) {
	// Получаем API-обработчик до регистрации closer'а: ленивая инициализация
	// зацепит за собой создание пула БД и зарегистрирует его в closer'е.
	// Closer работает по LIFO, поэтому пул должен попасть туда раньше http-сервера —
	// тогда при shutdown сначала остановится приём запросов, а уже потом закроется БД
	api := a.diContainer.IamV1API(ctx)

	server, err := iamv1.NewServer(api, iamv1.WithErrorHandler(iamv1API.ErrorHandler))
	if err != nil {
		slog.Error("ошибка создания HTTP-сервера", "error", err)
		return
	}

	// otelhttp — автоматические span'ы и метрики для каждого HTTP-запроса
	tracedHandler := otelhttp.NewHandler(server, config.AppConfig().OtelConfig.ServiceName)

	// Rate limiter — первым в цепочке, до otel: нет смысла тратить ресурсы
	// на трейсинг и обработку, если лимит уже превышен.
	// При недоступном Redis работает fail-open и пропускает запросы.
	rlCfg := config.AppConfig().RateLimit
	limitedHandler := ratelimit.Middleware(
		a.diContainer.RateLimiter(),
		rlCfg.Rate(),
		rlCfg.Burst(),
	)(tracedHandler)

	a.httpServer = &http.Server{
		Handler:           limitedHandler,
		ReadHeaderTimeout: httpReadHeaderTimeout,
		ReadTimeout:       httpReadTimeout,
		WriteTimeout:      httpWriteTimeout,
		IdleTimeout:       httpIdleTimeout,
	}
	closer.Add("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Shutdown(ctx)
	})
}

// runHTTPServer запускает HTTP-сервер и блокирует до его остановки
func (a *App) runHTTPServer() error {
	slog.Info("🚀 http-сервер запущен", "address", config.AppConfig().HTTP.Address())

	if err := a.httpServer.Serve(a.httpListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
