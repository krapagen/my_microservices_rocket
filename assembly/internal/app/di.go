package app

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/IBM/sarama"

	"github.com/krapagen/my_microservices_rocket/assembly/internal/config"
	orderPaidConsumer "github.com/krapagen/my_microservices_rocket/assembly/internal/consumer/order_paid"
	shipAssembledProducer "github.com/krapagen/my_microservices_rocket/assembly/internal/producer/ship_assembled"
	assemblyService "github.com/krapagen/my_microservices_rocket/assembly/internal/service/assembly"
	"github.com/krapagen/my_microservices_rocket/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/krapagen/my_microservices_rocket/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/krapagen/my_microservices_rocket/platform/pkg/kafka/producer"
	kafkaMiddleware "github.com/krapagen/my_microservices_rocket/platform/pkg/middleware/kafka"
)

type diContainer struct {
	// Инфраструктура
	syncProducer  sarama.SyncProducer
	consumerGroup sarama.ConsumerGroup

	// Kafka-обёртки
	orderPaidConsumer     *wrappedKafkaConsumer.Consumer
	shipAssembledProducer *wrappedKafkaProducer.Producer

	// Сервисы
	assemblerSvc         assemblyService.Assembler
	shipAssembledSvc     shipAssembledProducer.ProducerService
	orderPaidConsumerSvc orderPaidConsumer.ConsumerService
}

// SyncProducer возвращает синхронный Kafka-продюсер.
func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.KafkaBrokers(),
			newKafkaProducerConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать sync producer", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka sync producer", func(_ context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

// ConsumerGroup возвращает Kafka consumer group.
func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.KafkaBrokers(),
			config.AppConfig().OrderPaid.OrderPaidGroup(),
			newKafkaConsumerConfig(),
		)
		if err != nil {
			slog.Error("не удалось создать consumer group", "error", err)
			os.Exit(1)
		}

		closer.Add("Kafka consumer group", func(_ context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}

// OrderPaidConsumer возвращает обёртку Kafka-потребителя для события OrderPaid.
func (d *diContainer) OrderPaidConsumer() *wrappedKafkaConsumer.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{config.AppConfig().OrderPaid.OrderPaidTopicName()},
			wrappedKafkaConsumer.WithMiddlewares(kafkaMiddleware.ConsumerLogging()),
		)
	}

	return d.orderPaidConsumer
}

// ShipAssembledProducer возвращает обёртку Kafka-продюсера для события ShipAssembled.
func (d *diContainer) ShipAssembledProducer() *wrappedKafkaProducer.Producer {
	if d.shipAssembledProducer == nil {
		d.shipAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().ShipAssembled.ShipAssembledProducerTopic(),
		)
	}

	return d.shipAssembledProducer
}

// AssemblyService возвращает сервис бизнес-логики сборки.
func (d *diContainer) AssemblyService() assemblyService.Assembler {
	if d.assemblerSvc == nil {
		buildTimeSec := config.AppConfig().Assembler.AssembleLimitTimeSec
		if buildTimeSec < 0 {
			slog.Error("некорректный assembler.limit, должен быть >= 0", "limit", buildTimeSec)
			os.Exit(1)
		}

		buildTime := time.Duration(buildTimeSec) * time.Second
		d.assemblerSvc = assemblyService.NewService(buildTime)
	}

	return d.assemblerSvc
}

// ShipAssembledService возвращает сервис публикации события ShipAssembled.
func (d *diContainer) ShipAssembledService() shipAssembledProducer.ProducerService {
	if d.shipAssembledSvc == nil {
		d.shipAssembledSvc = shipAssembledProducer.NewService(d.ShipAssembledProducer())
	}

	return d.shipAssembledSvc
}

// OrderPaidConsumerService возвращает сервис обработки события OrderPaid.
func (d *diContainer) OrderPaidConsumerService() orderPaidConsumer.ConsumerService {
	if d.orderPaidConsumerSvc == nil {
		d.orderPaidConsumerSvc = orderPaidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.AssemblyService(),
			d.ShipAssembledService(),
		)
	}

	return d.orderPaidConsumerSvc
}

func newKafkaProducerConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	return cfg
}

func newKafkaConsumerConfig() *sarama.Config {
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest
	cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}
	return cfg
}
