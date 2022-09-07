package wallet

import (
	"context"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-wallet/entity"

	configgorm "github.com/Calmantara/go-common/infra/gorm"
)

type WalletRepoImpl struct {
	sugar logger.CustomLogger
	// conf config
	readCln  configgorm.PostgresConfig
	writeCln configgorm.PostgresConfig
}

func NewWalletRepo(sugar logger.CustomLogger, readCln configgorm.PostgresConfig, writeCln configgorm.PostgresConfig) WalletRepo {
	// read config and decide migrator
	if readCln.GetParam().Automigrate {
		sugar.Logger().Info("automigrate invoked for wallet")
		cln := readCln.GetClient()
		cln.AutoMigrate(entity.Wallet{})
	}

	sugar.Logger().Info("init wallet repo")
	return &WalletRepoImpl{sugar: sugar, readCln: readCln, writeCln: writeCln}
}

// read
func (w *WalletRepoImpl) ReadWallet(ctx context.Context, wallet *entity.Wallet) (err error) {
	w.sugar.WithContext(ctx).Infof("%T-ReadWallet is invoked", w)
	// generate transaction
	txn := w.readCln.GenerateTransaction(ctx)
	txn.Model(entity.Wallet{}).
		Select("id, created_at, record_flag").
		Where(configgorm.ActiveRecordQuery()).
		Where("id = ?", wallet.Id).Find(wallet)
	if err = txn.Error; err != nil {
		w.sugar.WithContext(ctx).Errorf("error execute ReadWallet:%v", err.Error())
	}
	w.sugar.WithContext(ctx).Infof("%T-ReadWallet executed", w)
	return err
}

// write
func (w *WalletRepoImpl) UpsertWallet(ctx context.Context, wallet *entity.Wallet) (err error) {
	w.sugar.WithContext(ctx).Infof("%T-InsertWallet is invoked", w)
	// generate transaction
	txn := w.writeCln.GenerateTransaction(ctx)
	txn.Model(entity.Wallet{}).
		Create(wallet)

	if err = txn.Error; err != nil {
		w.sugar.WithContext(ctx).Errorf("error execute InsertWallet:%v", err.Error())
	}
	w.sugar.WithContext(ctx).Infof("%T-InsertWallet executed", w)
	return err
}
