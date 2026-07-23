package order

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	repoConverter "github.com/krapagen/my_microservices_rocket/order/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/order/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	op := "order/internal/repository/order/Get"
	log := slog.With("op", op)

	type orderRow struct {
		UUID            uuid.UUID  `db:"uuid"`
		TransactionUUID *uuid.UUID `db:"transaction_uuid"`
		PaymentMethod   *string    `db:"payment_method"`
		Status          string     `db:"status"`
		CreatedAt       time.Time  `db:"created_at"`
		UpdatedAt       *time.Time `db:"updated_at"`
		PartUUID        *uuid.UUID `db:"part_uuid"`
		PartType        *string    `db:"part_type"`
		Price           *int64     `db:"price"`
	}

	query := squirrel.Select("orders.uuid", "orders.transaction_uuid", "orders.payment_method", "orders.status", "orders.created_at", "orders.updated_at", "order_items.part_uuid", "order_items.part_type", "order_items.price").
		From("orders").
		LeftJoin("order_items ON orders.uuid = order_items.order_uuid").
		Where(squirrel.Eq{"orders.uuid": orderUUID}).
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		log.ErrorContext(ctx, "ошибка запроса получения заказа", "error", err)
		return model.Order{}, fmt.Errorf("build query: %w", err)
	}

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(ctx, sql, args...)
	if err != nil {
		log.ErrorContext(ctx, "ошибка получения заказа", "error", err)
		return model.Order{}, fmt.Errorf("query order: %w", err)
	}
	defer rows.Close()

	orderRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[orderRow])
	if err != nil {
		log.ErrorContext(ctx, "ошибка чтения строк заказа", "error", err)
		return model.Order{}, fmt.Errorf("collect order rows: %w", err)
	}

	if len(orderRows) == 0 {
		log.ErrorContext(ctx, "Заказ не найден", "orderUUID", orderUUID)
		return model.Order{}, errs.ErrOrderNotFound
	}

	first := orderRows[0]
	order := record.Order{
		UUID:            first.UUID,
		TransactionUUID: first.TransactionUUID,
		PaymentMethod:   first.PaymentMethod,
		Status:          first.Status,
		CreatedAt:       first.CreatedAt,
		UpdatedAt:       first.UpdatedAt,
	}

	items := make([]record.OrderItem, 0, len(orderRows))
	for _, row := range orderRows {
		if row.PartUUID != nil && row.PartType != nil && row.Price != nil {
			items = append(items, record.OrderItem{
				OrderUUID: order.UUID,
				PartUUID:  *row.PartUUID,
				PartType:  *row.PartType,
				Price:     *row.Price,
			})
		}
	}

	log.InfoContext(ctx, "Заказ получен", "orderUUID", order.UUID)

	return repoConverter.OrderToModel(order, items), nil
}
