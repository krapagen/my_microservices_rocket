package order

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	repoConverter "github.com/krapagen/my_microservices_rocket/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/Update"
	log := slog.With("op", op)
	return r.txManager.Do(ctx, func(txCtx context.Context) error {
		if err := r.updateOrder(txCtx, order); err != nil {
			log.ErrorContext(txCtx, "ошибка обновления заказа", "error", err)
			return err
		}
		if err := r.updateOrderItems(txCtx, order); err != nil {
			log.ErrorContext(txCtx, "ошибка обновления списка деталей заказа", "error", err)
			return err
		}
		log.InfoContext(ctx, "Заказ обновлен", "orderUUID", order.UUID)
		return nil
	})
}

func (r *repository) updateOrder(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/updateOrder"
	log := slog.With("op", op)
	orderRecord, _ := repoConverter.OrderToRecord(order)

	query := `
		UPDATE orders
		SET updated_at       = $1,
		    transaction_uuid = COALESCE($2, transaction_uuid),
		    payment_method   = COALESCE($3, payment_method),
		    status           = COALESCE(NULLIF($4, ''), status)
		WHERE uuid = $5`

	res, err := r.getter.DefaultTrOrDB(ctx, r.pool).Exec(
		ctx, query,
		time.Now(),
		orderRecord.TransactionUUID,
		orderRecord.PaymentMethod,
		orderRecord.Status,
		orderRecord.UUID,
	)
	if err != nil {
		log.ErrorContext(ctx, "ошибка обновления заказа", "error", err)
		return fmt.Errorf("update order: %w", err)
	}

	if res.RowsAffected() == 0 {
		log.ErrorContext(ctx, "Заказ не найден для обновления", "orderUUID", order.UUID)
		return errs.ErrOrderNotFound
	}

	log.InfoContext(ctx, "Заказ обновлен", "orderUUID", order.UUID)
	return nil
}

func (r *repository) updateOrderItems(ctx context.Context, order model.Order) error {
	op := "order/internal/repository/order/updateOrderItems"
	log := slog.With("op", op)

	if len(order.Items) == 0 {
		log.InfoContext(ctx, "Список деталей заказа пуст, обновление не требуется")
		return nil
	}

	_, items := repoConverter.OrderToRecord(order)

	batch := &pgx.Batch{}

	query := `
    UPDATE order_items
    SET part_type = COALESCE($3, part_type),
        price     = COALESCE($4, price)
    WHERE order_uuid = $1 AND part_uuid = $2`

	for _, item := range items {
		batch.Queue(query, item.OrderUUID, item.PartUUID, item.PartType, item.Price)
	}

	br := r.getter.DefaultTrOrDB(ctx, r.pool).SendBatch(ctx, batch)
	defer br.Close()

	for _, item := range items {
		tag, err := br.Exec()
		if err != nil {
			log.ErrorContext(ctx, "ошибка обновления детали заказа", "error", err, "partUUID", item.PartUUID)
			return fmt.Errorf("update order item %s: %w", item.PartUUID, err)
		}
		if tag.RowsAffected() == 0 {
			log.ErrorContext(ctx, "Деталь заказа не найдена для обновления", "orderUUID", order.UUID, "partUUID", item.PartUUID)
			return fmt.Errorf("order item not found: orderUUID=%s, partUUID=%s: %w", order.UUID, item.PartUUID, errs.ErrOrderItemNotFound)
		}
	}

	log.InfoContext(ctx, "Список деталей заказа обновлен", "orderUUID", order.UUID)
	return nil
}
