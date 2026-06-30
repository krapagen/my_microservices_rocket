package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	paymentv1 "github.com/krapagen/my_microservices_rocket/shared/pkg/proto/payment/v1"
)

// Server реализует gRPC сервис оплаты
type Server struct {
	paymentv1.UnimplementedPaymentServiceServer
}

// NewServer создаёт новый экземпляр сервера оплаты
func NewServer() *Server {
	return &Server{}
}

// PayOrder обрабатывает оплату заказа
func (s *Server) PayOrder(
	ctx context.Context,
	req *paymentv1.PayOrderRequest,
) (*paymentv1.PayOrderResponse, error) {
	op := "Функция payment/pkg/service/PayOrder"
	log := slog.With("op", op)
	// 1. Проверить, что order_uuid не пустой → INVALID_ARGUMENT
	if req.GetOrderUuid() == "" {
		log.ErrorContext(ctx, "order_uuid не может быть пустым")
		return nil, status.Error(codes.InvalidArgument, "order_uuid не может быть пустым")
	}
	// 2. Проверить, что payment_method != UNSPECIFIED → INVALID_ARGUMENT
	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		log.ErrorContext(ctx, "payment_method не может быть UNSPECIFIED")
		return nil, status.Error(codes.InvalidArgument, "payment_method не может быть UNSPECIFIED")
	}
	// 3. Проверить формат UUID → INVALID_ARGUMENT
	if _, err := uuid.Parse(req.GetOrderUuid()); err != nil {
		log.ErrorContext(ctx, "order_uuid имеет неверный формат UUID", "error", err)
		return nil, status.Error(codes.InvalidArgument, "order_uuid имеет неверный формат UUID")
	}

	// 4. Сгенерировать transaction_uuid (UUID v4)

	transactionUuid := uuid.New()

	log.InfoContext(ctx, "Сгенерирован transaction_uuid")
	// 5. Вывести в лог: "оплата прошла успешно, order_uuid: X, transaction_uuid: Y"
	log.InfoContext(
		ctx, "Проверка входных данных прошла успешно, оплата прошла успешно",
		"order_uuid", req.GetOrderUuid(),
		"payment_method", req.GetPaymentMethod(),
	)
	// 6. Вернуть transaction_uuid
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUuid.String(),
	}, nil
}
