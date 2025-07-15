package repository

import "context"

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
