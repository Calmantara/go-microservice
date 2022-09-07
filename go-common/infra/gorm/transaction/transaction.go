package transaction

import (
	"context"
	"errors"

	"github.com/Calmantara/go-common/logger"
	"gorm.io/gorm"

	gormconf "github.com/Calmantara/go-common/infra/gorm"
)

type TransactionTerm string

func (t TransactionTerm) String() string {
	return string(t)
}

const (
	TRANSACTION_KEY       TransactionTerm = "GORM_TRANSACTION"
	TRANSACTION_ERROR_KEY TransactionTerm = "GORM_ERROR_TRANSACTION"
)

type Transaction interface {
	GormBeginTransaction(ctx context.Context) *gorm.DB
	GormEndTransaction(ctx context.Context) (err error)
}

type TransactionImpl struct {
	sugar        logger.CustomLogger
	postgresConf gormconf.PostgresConfig
}

func NewTransaction(sugar logger.CustomLogger, postgresConf gormconf.PostgresConfig) Transaction {
	return &TransactionImpl{
		sugar:        sugar,
		postgresConf: postgresConf,
	}
}

func (tr *TransactionImpl) GormBeginTransaction(ctx context.Context) *gorm.DB {
	txI := ctx.Value(TRANSACTION_KEY)
	if txI == nil {
		txI = ctx.Value(TRANSACTION_KEY.String())
	}
	if txI == nil {
		tr.sugar.WithContext(ctx).Info("creating new transaction")
		return tr.postgresConf.GetClient().Begin().WithContext(context.Background())
	}
	return txI.(*gorm.DB)
}
func (tr *TransactionImpl) GormEndTransaction(ctx context.Context) (err error) {
	txI := ctx.Value(TRANSACTION_KEY)
	if txI == nil {
		txI = ctx.Value(TRANSACTION_KEY.String())
	}
	if txI == nil {
		err = errors.New("cannot find transaction")
		return err
	}
	tx := txI.(*gorm.DB)

	if ctx.Value(TRANSACTION_ERROR_KEY) != nil ||
		ctx.Value(TRANSACTION_ERROR_KEY.String()) != nil {
		return tx.Rollback().Error
	}
	return tx.Commit().Error
}
