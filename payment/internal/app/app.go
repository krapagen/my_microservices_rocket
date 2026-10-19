package app

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/krapagen/my_microservices_rocket/payment/internal/config"
	"github.com/krapagen/my_microservices_rocket/payment/internal/interceptor"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/grpc/health"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/logger"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/tracing"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

const (
	grpcMaxConnectionIdle     = 15 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Second
	grpcKeepaliveTime         = 5 * time.Minute
	grpcKeepaliveTimeout      = 1 * time.Second
	grpcMinPingInterval       = 5 * time.Minute
	shutdownTimeout           = 5 * time.Second
)

// App — корневая структура приложения, управляющая жизненным циклом всех компонентов
type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
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

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runGRPCServer()
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

	api := a.diContainer.PaymentV1API(ctx)

	closer.Add("gRPC server", func(_ context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	// Регистрируем health service для проверки работоспособности
	health.RegisterService(a.grpcServer)

	paymentv1.RegisterPaymentServiceServer(a.grpcServer, api)
}

// runGRPCServer запускает gRPC-сервер и блокирует до его остановки
func (a *App) runGRPCServer() error {
	slog.Info("🚀 gRPC-сервер запущен", "address", config.AppConfig().GRPC.Address())

	return a.grpcServer.Serve(a.listener)
}
