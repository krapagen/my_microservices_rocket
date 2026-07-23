package order

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"

	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	repoConverter "github.com/krapagen/my_microservices_rocket/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/Create"
	log := slog.With("op", op)
	// Create атомарно сохраняет заказ и его строки
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.createOrder(txCtx, order); err != nil {
			log.ErrorContext(txCtx, "ошибка создания заказа", "error", err)
			return err
		}
		if err := r.createOrderItems(txCtx, order); err != nil {
			log.ErrorContext(txCtx, "ошибка создания списка деталей заказа", "error", err)
			return err
		}
		log.InfoContext(ctx, "Заказ создан", "orderUUID", order.UUID)
		return nil
	})
}

func (r *repository) createOrder(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/createOrder"
	log := slog.With("op", op)
	orderRecord, _ := repoConverter.OrderToRecord(order)
	query := squirrel.Insert("orders").
		Columns("uuid", "transaction_uuid", "payment_method", "status", "created_at", "updated_at").
		PlaceholderFormat(squirrel.Dollar).
		Values(orderRecord.UUID, orderRecord.TransactionUUID, orderRecord.PaymentMethod, orderRecord.Status, orderRecord.CreatedAt, orderRecord.UpdatedAt)

	sql, args, err := query.ToSql()
	if err != nil {
		log.ErrorContext(ctx, "ошибка запроса создания заказа", "error", err)
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		log.ErrorContext(ctx, "ошибка создания заказа", "error", err)
		return fmt.Errorf("insert order: %w", err)
	}

	log.InfoContext(ctx, "Заказ создан", "orderUUID", order.UUID)

	return nil
}

// order передаётся целиком: order_uuid берётся из order.UUID, остальные
// поля — из order.Items (в model.OrderItem ссылки на родителя нет).
func (r *repository) createOrderItems(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/createOrderItems"
	log := slog.With("op", op)
	if len(order.Items) == 0 {
		log.InfoContext(ctx, "Список деталей заказа пуст")
		return nil
	}

	_, items := repoConverter.OrderToRecord(order)

	query := squirrel.Insert("order_items").
		Columns("order_uuid", "part_uuid", "part_type", "price").
		PlaceholderFormat(squirrel.Dollar)

	for _, item := range items {
		query = query.Values(item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		log.ErrorContext(ctx, "ошибка создания запроса списка деталей заказа", "error", err)
		return fmt.Errorf("build query: %w", err)
	}

	_, err = r.getter.DefaultTrOrDB(ctx, r.pool).Exec(ctx, sql, args...)
	if err != nil {
		log.ErrorContext(ctx, "ошибка создания списка деталей заказа", "error", err)
		return fmt.Errorf("insert order items: %w", err)
	}

	log.InfoContext(ctx, "Список деталей заказа создан", "orderUUID", order.UUID)

	return nil
}
