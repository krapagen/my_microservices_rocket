package order

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	errs "github.com/krapagen/my_microservices_rocket/order/internal/errors"
	"github.com/krapagen/my_microservices_rocket/order/internal/model"
	repoConverter "github.com/krapagen/my_microservices_rocket/order/internal/repository/converter"
	"github.com/krapagen/my_microservices_rocket/order/internal/repository/record"
)

// TxManager определяет контракт для управления транзакциями
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type repository struct {
	pool      *pgxpool.Pool
	getter    *trmpgx.CtxGetter
	txManager TxManager
}

// New создаёт новый репозиторий заказов.
func New(pool *pgxpool.Pool, txManager TxManager) *repository {
	return &repository{
		pool:      pool,
		getter:    trmpgx.DefaultCtxGetter,
		txManager: txManager,
	}
}

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
