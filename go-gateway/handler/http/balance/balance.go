package balance

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-common/topic"
	"github.com/Calmantara/go-gateway/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"

	"github.com/Calmantara/go-common/pb"
	pbe "github.com/Calmantara/go-common/pb"
	"github.com/Calmantara/go-wallet/entity"

	grpcclient "github.com/Calmantara/go-common/infra/grpc/client"
	cmodel "github.com/Calmantara/go-common/model"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"
	emitterclient "github.com/Calmantara/go-emitter/handler/grpc/emitter/client"
	balanceclient "github.com/Calmantara/go-wallet/handler/grpc/balance/client"
	wmodel "github.com/Calmantara/go-wallet/model"
)

type BalanceHdl interface {
	GetBalanceDetail(ctx *gin.Context)
	PostBalance(ctx *gin.Context)
}

type BalanceHdlImpl struct {
	sugar         logger.CustomLogger
	util          serviceutil.UtilService
	assert        serviceassert.Assert
	confWallet    model.GoWallet
	confEmitter   model.GoEmitter
	confService   model.Service
	balanceClient balanceclient.BalanceClient
	emitterCln    emitterclient.EmitterClient
}

func NewBalanceHdl(sugar logger.CustomLogger, config config.ConfigSetup,
	util serviceutil.UtilService, assert serviceassert.Assert) BalanceHdl {

	// get config
	var confWallet model.GoWallet
	config.GetConfig("gowallet", &confWallet)
	var confEmitter model.GoEmitter
	config.GetConfig("goemitter", &confEmitter)
	var confService model.Service
	config.GetConfig("service", &confService)

	// init client
	grpcWalletCln := grpcclient.NewGRPCClientConnection(sugar, confWallet.Host)
	grpcEmitterCln := grpcclient.NewGRPCClientConnection(sugar, confEmitter.Host)
	// get clients
	balanceCln := balanceclient.NewBalanceClient(sugar, grpcWalletCln)
	emitterCln := emitterclient.NewEmitterClient(sugar, grpcEmitterCln)

	return &BalanceHdlImpl{sugar: sugar,
		confWallet: confWallet, confService: confService, confEmitter: confEmitter,
		util: util, assert: assert,
		balanceClient: balanceCln, emitterCln: emitterCln}
}

func (w *BalanceHdlImpl) GetBalanceDetail(ctx *gin.Context) {
	// set correlation id
	w.util.SetCorrelationIdFromHeader(ctx)
	// check payload
	w.sugar.WithContext(ctx).Info("checking wallet id")
	walletIdStr, ok := ctx.GetQuery("wallet_id")
	if !ok {
		w.sugar.WithContext(ctx).Error("wallet id is empty")
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_BAD_REQUEST_TYPE, "wallet id is not found")
		return
	}
	// validate balance
	walletId, _ := strconv.Atoi(walletIdStr)
	if w.assert.IsZero(walletId) {
		w.sugar.WithContext(ctx).Error("wallet id is zero")
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_BAD_REQUEST_TYPE, "wallet id cannot zero")
		return
	}
	// init retries and timeout
	var resp *pb.BalanceResponse
	var err error

	// routing based on api version
	url := ctx.Request.URL.Path

	method := w.balanceClient.GetClient().GetBalance
	if strings.Contains(url, "v2") {
		method = w.balanceClient.GetClient().GetBalanceByTtl
	}

	w.sugar.WithContext(ctx).Info("fetching go-wallet")
	for i := 0; i < w.confWallet.MaxRetries; i++ {
		ctxCln, cancel := context.WithTimeout(ctx, time.Second*time.Duration(w.confWallet.Timeout))
		defer cancel()
		// fetching
		ctxCln = w.util.InsertCorrelationIdToGrpc(ctxCln)
		resp, err = method(ctxCln, &pb.Wallet{Id: int64(walletId)})
		if err == nil {
			// indicate does not error happen
			break
		}
		w.sugar.WithContext(ctx).Info("failed to fetch go-wallet, retries num:%v with err:%v", i, err.Error())
		// wait every 10 seconds
		time.Sleep(time.Second * time.Duration(w.confWallet.TimeToWait))
	}

	if err != nil {
		w.sugar.WithContext(ctx).Error("error fetching to go-wallet:%v", err.Error())
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_INTERNAL_TYPE, "connection error")
		return
	}

	// transform error message and check
	errMsg := cmodel.ErrorModel{
		Error:     errors.New(resp.GetErrorMessage().Error),
		ErrorType: cmodel.ResponseType(resp.ErrorMessage.GetErrorType()),
	}
	if !w.assert.IsEmpty(errMsg.Error.Error()) || !w.assert.IsEmpty(errMsg.ErrorType.String()) {
		w.sugar.WithContext(ctx).Error("error fetching processing payload to go-wallet:%v", errMsg.Error.Error())
		w.util.UtilErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error.Error())
		return
	}
	// valid request
	var balance wmodel.WalletBalance
	if err = w.util.ObjectMapper(resp.GetBalanceDetail(), &balance); err != nil {
		w.sugar.WithContext(ctx).Error("error transforming response go-emitter:%v", err.Error())
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_INTERNAL_TYPE, "error transforming payload")
		return
	}
	//success
	w.util.UtilResponseSwitcher(ctx, cmodel.SUCCESS_OK_TYPE, balance)
}

func (w *BalanceHdlImpl) PostBalance(ctx *gin.Context) {
	// set correlation id
	w.util.SetCorrelationIdFromHeader(ctx)
	// check payload
	//generate balance payload
	balancePayload := entity.Balance{}
	if err := ctx.ShouldBind(&balancePayload); err != nil {
		w.sugar.WithContext(ctx).Error("failed to binding payload is zero:%v", err.Error())
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_BAD_REQUEST_TYPE, "payload is invalid")
		return
	}

	w.sugar.WithContext(ctx).Info("checking wallet id")
	// validate balance
	if w.assert.IsZero(int(balancePayload.WalletId)) {
		w.sugar.WithContext(ctx).Error("wallet id is zero")
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_BAD_REQUEST_TYPE, "wallet id cannot zero")
		return
	}

	// transform balance payload to string
	b, _ := json.Marshal(&balancePayload)
	balancePayloadStr := string(b)

	// generate emitter proto
	emitterProto := &pbe.Emitter{
		Topic:   topic.BALANCE_TRANSACTION_TOPIC.String(),
		Issuer:  w.confService.Name,
		Message: balancePayloadStr,
	}

	// init retries and timeout
	var resp *pbe.EmitterResponse
	var err error

	w.sugar.WithContext(ctx).Info("fetching go-emitter")
	for i := 0; i < w.confEmitter.MaxRetries; i++ {
		ctxCln, cancel := context.WithTimeout(ctx, time.Second*time.Duration(w.confEmitter.Timeout))
		defer cancel()
		// fetching
		ctxCln = w.util.InsertCorrelationIdToGrpc(ctxCln)
		resp, err = w.emitterCln.GetClient().SendEmitterPayload(ctxCln, emitterProto)
		if err == nil {
			// indicate does not error happen
			break
		}
		w.sugar.WithContext(ctx).Infof("failed to fetch go-emitter, retries num:%v with err:%v", i, err.Error())
		// wait every 10 seconds
		time.Sleep(time.Second * time.Duration(w.confEmitter.TimeToWait))
	}

	if err != nil {
		w.sugar.WithContext(ctx).Error("error fetching to go-emitter:%v", err.Error())
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_INTERNAL_TYPE, "connection error")
		return
	}

	// transform error message and check
	errMsg := cmodel.ErrorModel{
		Error:     errors.New(resp.GetErrorMessage().Error),
		ErrorType: cmodel.ResponseType(resp.ErrorMessage.GetErrorType()),
	}
	if !w.assert.IsEmpty(errMsg.Error.Error()) || !w.assert.IsEmpty(errMsg.ErrorType.String()) {
		w.sugar.WithContext(ctx).Error("error fetching processing payload to go-emitter:%v", errMsg.Error.Error())
		w.util.UtilErrorResponseSwitcher(ctx, errMsg.ErrorType, errMsg.Error.Error())
		return
	}
	// valid request
	var balance entity.Balance
	if err = json.Unmarshal([]byte(resp.GetEmitterDetail().GetMessage()), &balance); err != nil {
		w.sugar.WithContext(ctx).Error("error transforming response go-emitter:%v", err.Error())
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_INTERNAL_TYPE, "error transforming payload")
		return
	}
	//success
	w.util.UtilResponseSwitcher(ctx, cmodel.SUCCESS_OK_TYPE, map[string]any{
		"message": "balance transaction succes",
		"payload": balance})
}
