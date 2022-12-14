package balance

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/pb"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-wallet/entity"
	"github.com/Calmantara/go-wallet/model"
	"github.com/Calmantara/go-wallet/repository/balance"
	"github.com/Calmantara/go-wallet/service/wallet"
	"github.com/lovoo/goka"

	cmodel "github.com/Calmantara/go-common/model"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	redisservice "github.com/Calmantara/go-common/service/redis"
	serviceutil "github.com/Calmantara/go-common/service/util"
)

type BalanceSvcImpl struct {
	conf model.BalanceConf

	sugar       logger.CustomLogger
	config      config.ConfigSetup
	balanceRepo balance.BalanceRepo
	util        serviceutil.UtilService
	assert      serviceassert.Assert
	redis       redisservice.RedisService
	walletSvc   wallet.WalletSvc
}

func NewBalanceSvc(sugar logger.CustomLogger, config config.ConfigSetup,
	balanceRepo balance.BalanceRepo, util serviceutil.UtilService,
	assert serviceassert.Assert, redis redisservice.RedisService, walletSvc wallet.WalletSvc) BalanceSvc {
	conf := model.BalanceConf{}
	config.GetConfig("balanceconf", &conf)

	return &BalanceSvcImpl{
		conf: conf, config: config, balanceRepo: balanceRepo,
		util: util, assert: assert, redis: redis, walletSvc: walletSvc, sugar: sugar,
	}
}

func (b *BalanceSvcImpl) GetBalanceDetail(ctx context.Context, balance *model.WalletBalance) (errModel cmodel.ErrorModel) {
	b.sugar.WithContext(ctx).Infof("%T-GetBalanceDetail is invoked", b)
	defer b.sugar.WithContext(ctx).Infof("%T-GetBalanceDetail executed", b)

	// checking wallet exist or not
	wallet := entity.Wallet{Id: balance.WalletId}
	if errModel = b.walletSvc.GetWalletDetail(ctx, &wallet); errModel.Error != nil {
		return errModel
	}

	b.sugar.WithContext(ctx).Info("getting balance in database")
	if err := b.balanceRepo.ReadBalance(ctx, balance); err != nil {
		b.sugar.WithContext(ctx).Errorf("error getting balance in database:%v", err)
		errModel = cmodel.ErrorModel{
			Error:     err,
			ErrorType: cmodel.ERR_INTERNAL_TYPE,
		}
	}
	return errModel
}
func (b *BalanceSvcImpl) InsertBalance(ctx context.Context, balance *entity.Balance) (errModel cmodel.ErrorModel) {
	b.sugar.WithContext(ctx).Infof("%T-UpsertBalance is invoked", b)
	defer b.sugar.WithContext(ctx).Infof("%T-UpsertBalance executed", b)

	b.sugar.WithContext(ctx).Info("insert balance to database")
	if err := b.balanceRepo.InsertBalance(ctx, balance); err != nil {
		b.sugar.WithContext(ctx).Errorf("error getting balance in database:%v", err)
		errModel = cmodel.ErrorModel{
			Error:     err,
			ErrorType: cmodel.ERR_INTERNAL_TYPE,
		}
		return errModel
	}
	// checking balance
	if balance.CreatedAt == nil {
		err := errors.New("balance is not inserted")
		errModel = cmodel.ErrorModel{
			Error:     err,
			ErrorType: cmodel.ERR_NO_WALLET_TYPE,
		}
		b.sugar.WithContext(ctx).Errorf("error upserting balance:%v", err)
	}

	// set redis cache for thresholding
	ctxBg := b.util.ContextBackground(ctx)
	go func(blc *entity.Balance) {
		var walletBalance model.WalletBalance
		// get redis data first
		b.sugar.WithContext(ctxBg).Infof("checking balance in redis:%v", model.BALANCE_KEY.Append(blc.WalletId))
		if err := b.redis.Get(ctxBg, model.BALANCE_KEY.Append(blc.WalletId), &walletBalance); err != nil {
			b.sugar.WithContext(ctxBg).Infof("redis is empty for:%v", model.BALANCE_KEY.Append(blc.WalletId))
			// set redis key
			walletBalance = model.WalletBalance{
				WalletId:  blc.WalletId,
				Amount:    blc.Amount,
				LastTopUp: blc.CreatedAt,
			}
			if err = b.redis.Set(ctxBg, model.BALANCE_KEY.Append(blc.WalletId), &walletBalance, time.Duration(b.conf.RedisTtl*int(time.Minute))); err != nil {
				b.sugar.WithContext(ctxBg).Errorf("error when set in redis:%v", err)
			}
			return
		}
		// if exist
		timeTwoMinutes := walletBalance.LastTopUp.Add(time.Duration(b.conf.RedisTtl * int(time.Minute)))
		timeNow := time.Now()

		if timeNow.After(timeTwoMinutes) {
			// no need to increase per window
			return
		}
		walletBalance.Amount += blc.Amount
		b.redis.Set(ctxBg, model.BALANCE_KEY.Append(blc.WalletId), &walletBalance, timeTwoMinutes.Sub(timeNow))
	}(balance)

	return errModel
}

func (b *BalanceSvcImpl) GetBalanceDetailByTtl(ctx context.Context, balance *model.WalletBalance) (errModel cmodel.ErrorModel) {
	// checking wallet exist or not
	wallet := entity.Wallet{Id: balance.WalletId}
	if errModel = b.walletSvc.GetWalletDetail(ctx, &wallet); errModel.Error != nil {
		return errModel
	}

	balanceChan := make(chan model.WalletBalance)
	go func(walletId uint64) {
		var walletBalance model.WalletBalance
		if err := b.redis.Get(ctx, model.BALANCE_KEY.Append(walletId), &walletBalance); err != nil {
			b.sugar.WithContext(ctx).Info("cache is empty")
		}
		balanceChan <- walletBalance
	}(balance.WalletId)

	// fetching from database aswell
	if err := b.balanceRepo.ReadSumBalance(ctx, balance); err != nil {
		b.sugar.WithContext(ctx).Errorf("error fetching balance:%v", err)
		err = errors.New("error connection")
		errModel = cmodel.ErrorModel{
			Error:     err,
			ErrorType: cmodel.ERR_INTERNAL_TYPE,
		}
		return errModel
	}
	// checking window value
	blc := <-balanceChan
	if blc.Amount > int64(b.conf.Threshold) {
		balance.AboveThreshold = true
	}
	return errModel
}

// consumer
func (b *BalanceSvcImpl) ConsumeKafkaPayload() goka.ProcessCallback {
	return func(ctx goka.Context, msg interface{}) {
		// get header first
		header := ctx.Headers()[logger.CorrelationKey.String()]
		ctxWithValue := context.WithValue(ctx.Context(), logger.CorrelationKey.String(), string(header))

		// check deduplication key
		var proceed bool
		b.redis.Get(ctxWithValue, redisservice.RedisKey(ctx.Key()), &proceed)
		if proceed {
			b.sugar.WithContext(ctxWithValue).Infof("payload already proceed with key:%v", ctx.Key())
			return
		}

		// transform to balance entity
		emitterPayload := &pb.Emitter{}
		json.Unmarshal([]byte(msg.(string)), &emitterPayload)
		// transform to balance
		var balance entity.Balance
		json.Unmarshal([]byte(emitterPayload.GetMessage()), &balance)
		b.sugar.WithContext(ctxWithValue).Info("inserting payload for wallet:%v", balance.WalletId)
		if errMsg := b.InsertBalance(ctxWithValue, &balance); errMsg.Error != nil {
			b.sugar.WithContext(ctxWithValue).Errorf("error when inserting:%v", errMsg.Error)
			return
		}
		// set deduplication code in redis with ttl 7 days
		if err := b.redis.Set(ctxWithValue, redisservice.RedisKey(ctx.Key()), true, 24*7*time.Hour); err != nil {
			b.sugar.WithContext(ctxWithValue).Errorf("error when set dudiplication key:%v", err.Error())
		}
	}
}
