package postgres

import (
	"context"

	"github.com/mohamedfawas/quboolkallyanam.xyz/pkg/database/postgres"
	"github.com/mohamedfawas/quboolkallyanam.xyz/services/payment/internal/domain/repository"
	"gorm.io/gorm"
)

type txManager struct {
	db *postgres.Client
}

func NewTxManager(db *postgres.Client) repository.TxManager {
	return &txManager{db: db}
}

type txKey struct{} // used to carry *gorm.DB in context

func (t *txManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return t.db.GormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

// Helper function to get transaction from context or fallback to main DB
func GetDBFromContext(ctx context.Context, mainDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return mainDB
}
