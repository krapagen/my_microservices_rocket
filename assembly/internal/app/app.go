package app

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/config"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/logger"
)

const shutdownTimeout = 5 * time.Second

// App — корневая структура приложения, управляющая жизненным циклом компонентов.
type App struct {
	diContainer *diContainer
}

// New создаёт и инициализирует приложение.
func New(_ context.Context) *App {
	a := &App{}
	a.initDeps()
	return a
}

// Run управляет жизненным циклом Kafka-consumer:
// запускает потребление, ждёт сигнал завершения и выполняет graceful shutdown.
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.runConsumer(ctx)
	}()

	var runErr error
	select {
	case runErr = <-errCh:
		// consumer завершился сам (с ошибкой или nil)
	case <-ctx.Done():
		slog.Info("получен сигнал завершения, начинаем graceful shutdown")
	}
	cancel()

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

func (a *App) initDeps() {
	a.diContainer = &diContainer{}
	logger.Init(config.AppConfig().Logger.Level)
}

func (a *App) runConsumer(ctx context.Context) error {
	slog.Info(
		"🚀 Kafka-потребитель OrderPaid запущен",
		"topic", config.AppConfig().OrderPaid.OrderPaidTopicName(),
		"group_id", config.AppConfig().OrderPaid.OrderPaidGroup(),
	)

	return a.diContainer.OrderPaidConsumerService().RunConsumer(ctx)
}
