package user

import (
	"context"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TxManager определяет контракт для управления транзакциями
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type repository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
	txMgr  TxManager
}

func New(pool *pgxpool.Pool, txMgr TxManager) *repository {
	return &repository{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
		txMgr:  txMgr,
	}
}
