package wallet

import (
	"errors"
	"strconv"
	"time"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-gateway/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"

	"github.com/Calmantara/go-common/pb"
	"github.com/Calmantara/go-wallet/entity"

	grpcclient "github.com/Calmantara/go-common/infra/grpc/client"
	cmodel "github.com/Calmantara/go-common/model"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"
	walletclient "github.com/Calmantara/go-wallet/handler/grpc/wallet/client"
)

type WalletHdl interface {
	GetWalletDetail(ctx *gin.Context)
	PostWallet(ctx *gin.Context)
}

type WalletHdlImpl struct {
	sugar        logger.CustomLogger
	util         serviceutil.UtilService
	assert       serviceassert.Assert
	config       model.GoWallet
	walletClient walletclient.WalletClient
}

func NewWalletHdl(sugar logger.CustomLogger, config config.ConfigSetup,
	util serviceutil.UtilService, assert serviceassert.Assert) WalletHdl {

	// get config
	var conf model.GoWallet
	config.GetConfig("gowallet", &conf)

	// init client
	grpcCln := grpcclient.NewGRPCClientConnection(sugar, conf.Host)
	// get wallet client
	walletCln := walletclient.NewWalletClient(sugar, grpcCln)

	return &WalletHdlImpl{sugar: sugar, config: conf, util: util, assert: assert, walletClient: walletCln}
}

func (w *WalletHdlImpl) GetWalletDetail(ctx *gin.Context) {
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
	// validate wallet
	walletId, _ := strconv.Atoi(walletIdStr)
	if w.assert.IsZero(walletId) {
		w.sugar.WithContext(ctx).Error("wallet id is zero")
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_BAD_REQUEST_TYPE, "wallet id cannot zero")
		return
	}
	// init retries and timeout
	var resp *pb.WalletResponse
	var err error

	w.sugar.WithContext(ctx).Info("fetching go-wallet")
	for i := 0; i < w.config.MaxRetries; i++ {
		ctxCln, cancel := context.WithTimeout(ctx, time.Minute*time.Duration(w.config.Timeout))
		defer cancel()
		// fetching
		ctxCln = w.util.InsertCorrelationIdToGrpc(ctxCln)
		resp, err = w.walletClient.GetClient().GetWallet(ctxCln, &pb.Wallet{Id: int64(walletId)})
		if err == nil {
			// indicate does not error happen
			break
		}
		w.sugar.WithContext(ctx).Infof("failed to fetch, retries num:%v with err:%v", i, err.Error())
		// wait every 10 seconds
		time.Sleep(time.Second * time.Duration(w.config.TimeToWait))
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
	var wallet entity.Wallet
	if err = w.util.ObjectMapper(resp.GetWalletDetail(), &wallet); err != nil {
		w.util.UtilErrorResponseSwitcher(ctx, cmodel.ERR_INTERNAL_TYPE, "error transforming payload")
		return
	}
	//success
	w.util.UtilResponseSwitcher(ctx, cmodel.SUCCESS_OK_TYPE, wallet)
}

func (w *WalletHdlImpl) PostWallet(ctx *gin.Context) {}
