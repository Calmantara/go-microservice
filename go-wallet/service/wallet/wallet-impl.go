package wallet

import (
	"context"
	"errors"
	"time"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-wallet/entity"
	"github.com/Calmantara/go-wallet/model"
	"github.com/Calmantara/go-wallet/repository/wallet"

	cmodel "github.com/Calmantara/go-common/model"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	redisservice "github.com/Calmantara/go-common/service/redis"
	serviceutil "github.com/Calmantara/go-common/service/util"
)

type WalletSvcImpl struct {
	conf model.WalletConf

	sugar      logger.CustomLogger
	config     config.ConfigSetup
	walletRepo wallet.WalletRepo
	util       serviceutil.UtilService
	assert     serviceassert.Assert
	redis      redisservice.RedisService
}

func NewWalletSvc(sugar logger.CustomLogger, config config.ConfigSetup,
	walletRepo wallet.WalletRepo, util serviceutil.UtilService,
	assert serviceassert.Assert, redis redisservice.RedisService) WalletSvc {
	conf := model.WalletConf{}
	config.GetConfig("walletconf", &conf)

	return &WalletSvcImpl{
		conf: conf, config: config, walletRepo: walletRepo,
		util: util, assert: assert, redis: redis, sugar: sugar,
	}
}

func (w *WalletSvcImpl) GetWalletDetail(ctx context.Context, wallet *entity.Wallet) (errModel cmodel.ErrorModel) {
	w.sugar.WithContext(ctx).Infof("%T-GetWalletDetail is invoked", w)
	defer w.sugar.WithContext(ctx).Infof("%T-GetWalletDetail executed", w)
	// get redis data first
	w.sugar.WithContext(ctx).Info("checking wallet in redis")
	if err := w.redis.Get(ctx, model.WALLET_KEY.Append(wallet.Id), wallet); err != nil {
		w.sugar.WithContext(ctx).Infof("redis is empty for:%v", model.WALLET_KEY.Append(wallet.Id))

		w.sugar.WithContext(ctx).Info("getting wallet in database")
		if err = w.walletRepo.ReadWallet(ctx, wallet); err != nil {
			w.sugar.WithContext(ctx).Errorf("error getting wallet in database:%v", err)
			errModel = cmodel.ErrorModel{
				Error:     err,
				ErrorType: cmodel.ERR_INTERNAL_TYPE,
			}
			return errModel
		}
	}
	// checking wallet
	if wallet.CreatedAt == nil {
		err := errors.New("wallet is not created")
		errModel = cmodel.ErrorModel{
			Error:     err,
			ErrorType: cmodel.ERR_NO_WALLET_TYPE,
		}
		w.sugar.WithContext(ctx).Errorf("error getting wallet:%v", err)
		return
	}

	// set to redis
	ctxBg := w.util.ContextBackground(ctx)
	go func() {
		w.redis.Set(ctxBg, model.WALLET_KEY.Append(wallet.Id), wallet, time.Duration(w.conf.RedisTtl*int(time.Minute)))
	}()

	return errModel
}
func (w *WalletSvcImpl) UpsertWallet(ctx context.Context, wallet *entity.Wallet) (errModel cmodel.ErrorModel) {
	w.sugar.WithContext(ctx).Infof("%T-UpsertWallet is invoked", w)
	defer w.sugar.WithContext(ctx).Infof("%T-UpsertWallet executed", w)

	//delete redis cache
	ctxBg := w.util.ContextBackground(ctx)
	go w.redis.Delete(ctxBg, model.WALLET_KEY.Append(wallet.Id))

	w.sugar.WithContext(ctx).Info("upsert wallet to database")
	if err := w.walletRepo.UpsertWallet(ctx, wallet); err != nil {
		w.sugar.WithContext(ctx).Errorf("error getting wallet in database:%v", err)
		errModel = cmodel.ErrorModel{
			Error:     err,
			ErrorType: cmodel.ERR_INTERNAL_TYPE,
		}
		return errModel
	}
	// checking wallet
	if wallet.CreatedAt == nil {
		err := errors.New("wallet is not inserted")
		errModel = cmodel.ErrorModel{
			Error:     err,
			ErrorType: cmodel.ERR_NO_WALLET_TYPE,
		}
		w.sugar.WithContext(ctx).Errorf("error upserting wallet:%v", err)
	}
	return errModel
}
