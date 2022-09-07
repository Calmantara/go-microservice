package balance

import (
	"context"
	"time"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-wallet/entity"
	"github.com/Calmantara/go-wallet/model"

	configgorm "github.com/Calmantara/go-common/infra/gorm"
)

type BalanceRepoImpl struct {
	sugar    logger.CustomLogger
	readCln  configgorm.PostgresConfig
	writeCln configgorm.PostgresConfig
}

func NewBalanceRepo(sugar logger.CustomLogger, readCln configgorm.PostgresConfig, writeCln configgorm.PostgresConfig) BalanceRepo {
	// read config and decide migrator
	balance := &BalanceRepoImpl{sugar: sugar, readCln: readCln, writeCln: writeCln}
	if readCln.GetParam().Automigrate {
		sugar.Logger().Info("automigrate invoked for balance")
		cln := readCln.GetClient()
		cln.AutoMigrate(entity.Balance{})
		// generate view
		sugar.Logger().Info("automigrate invoked generating view")
		db, _ := cln.DB()
		db.Exec(balance.GenerateWalletBalanceView())
	}
	sugar.Logger().Info("init balance repo")
	return balance
}

// read
func (b *BalanceRepoImpl) ReadBalance(ctx context.Context, walletBalance *model.WalletBalance) (err error) {
	b.sugar.WithContext(ctx).Infof("%T-ReadBalance is invoked", b)
	// generate transaction
	txn := b.readCln.GenerateTransaction(ctx)
	// abc := map[string]any{}
	txn.Raw(b.GetWalletBalanceView(), time.Now(), walletBalance.WalletId).
		Scan(&walletBalance)
	if err = txn.Error; err != nil {
		b.sugar.WithContext(ctx).Errorf("error execute ReadBalance:%v", err.Error())
	}
	b.sugar.WithContext(ctx).Infof("%T-ReadBalance executed", b)
	return err
}

// write
func (b *BalanceRepoImpl) InsertBalance(ctx context.Context, balance *entity.Balance) (err error) {
	b.sugar.WithContext(ctx).Infof("%T-InsertBalance is invoked", b)
	// generate transaction
	txn := b.writeCln.GenerateTransaction(ctx)
	txn.Model(entity.Balance{}).
		Create(balance)
	if err = txn.Error; err != nil {
		b.sugar.WithContext(ctx).Errorf("error execute InsertBalance:%v", err.Error())
	}
	b.sugar.WithContext(ctx).Infof("%T-InsertBalance executed", b)
	return err
}
