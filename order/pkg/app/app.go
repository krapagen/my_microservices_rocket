package app

// server, err := orderv1.NewServer(apiHandler, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	orderv1API "github.com/krapagen/my_microservices_rocket/order/internal/api/order/v1"
	iamClient "github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/iam/v1"
	inventoryClient "github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/krapagen/my_microservices_rocket/order/internal/client/grpc/payment/v1"
	"github.com/krapagen/my_microservices_rocket/order/internal/middleware"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	orderProducer "github.com/krapagen/my_microservices_rocket/order/internal/producer/order_producer"
	orderRepository "github.com/krapagen/my_microservices_rocket/order/internal/repository/order"
	service "github.com/krapagen/my_microservices_rocket/order/internal/service/order"
	orderv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/openapi/order/v1"
	authv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/auth/v1"
	inventoryv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

// NewHTTPHandler creates HTTP handler from gRPC clients (for tests)
func NewHTTPHandler(pool *pgxpool.Pool, txManager orderRepository.TxManager, inventoryGRPCClient inventoryv1.InventoryServiceClient, paymentGRPCClient paymentv1.PaymentServiceClient, authClient ...authv1.AuthServiceClient) (http.Handler, error) {
	var client authv1.AuthServiceClient
	if len(authClient) > 0 {
		client = authClient[0]
	}

	return NewHTTPHandlerWithProducer(
		pool,
		txManager,
		inventoryGRPCClient,
		paymentGRPCClient,
		client,
		noopOrderPaidProducer{},
	)
}

// NewHTTPHandlerWithProducer creates HTTP handler and uses provided OrderPaid producer.
func NewHTTPHandlerWithProducer(
	pool *pgxpool.Pool,
	txManager orderRepository.TxManager,
	inventoryGRPCClient inventoryv1.InventoryServiceClient,
	paymentGRPCClient paymentv1.PaymentServiceClient,
	authClient authv1.AuthServiceClient,
	orderPaidProducer orderProducer.ProducerService,
) (http.Handler, error) {
	// Repository layer
	orderRepo := orderRepository.New(pool, txManager)

	// Create client adapters
	invClient := inventoryClient.New(inventoryGRPCClient)
	payClient := paymentClient.New(paymentGRPCClient)

	// Service layer
	orderService := service.New(orderRepo, invClient, payClient, orderPaidProducer, txManager)

	// API layer
	api := orderv1API.NewAPI(orderService)

	// Create OpenAPI server with error handler
	server, err := orderv1.NewServer(api, orderv1.WithErrorHandler(orderv1API.ErrorHandler))
	if err != nil {
		return nil, err
	}

	if authClient != nil {
		return middleware.HTTPAuth(iamClient.New(authClient))(server), nil
	}

	return server, nil
}

type noopOrderPaidProducer struct{}

func (noopOrderPaidProducer) ProduceOrderPaid(_ context.Context, _ model.OrderPaid) error {
	return nil
}
