package walletserver

import (
	"context"

	"github.com/Calmantara/go-common/infra/gorm/transaction"
	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/model"
	"github.com/Calmantara/go-common/pb"
	"github.com/Calmantara/go-wallet/entity"
	"github.com/Calmantara/go-wallet/service/wallet"

	grpcserver "github.com/Calmantara/go-common/infra/grpc/server"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"
)

type WalletServerImpl struct {
	sugar     logger.CustomLogger
	readTrx   transaction.Transaction
	server    grpcserver.GRPCServer
	util      serviceutil.UtilService
	assert    serviceassert.Assert
	walletSvc wallet.WalletSvc
	pb.UnimplementedWalletServiceServer
}

func NewWalletServer(sugar logger.CustomLogger,
	readTrx transaction.Transaction,
	server grpcserver.GRPCServer,
	util serviceutil.UtilService,
	assert serviceassert.Assert,
	walletSvc wallet.WalletSvc) WalletServer {

	walletServer := &WalletServerImpl{
		sugar:     sugar,
		readTrx:   readTrx,
		util:      util,
		assert:    assert,
		walletSvc: walletSvc,
		server:    server,
	}
	pb.RegisterWalletServiceServer(server.GetServer(), walletServer)
	return walletServer
}

func (w *WalletServerImpl) GetWallet(ctx context.Context, wallet *pb.Wallet) (walletResp *pb.WalletResponse, err error) {
	// init response
	walletResp = &pb.WalletResponse{}
	// set transaction
	ctx = context.WithValue(ctx, transaction.TRANSACTION_KEY, w.readTrx.GormBeginTransaction(ctx))
	w.sugar.WithContext(ctx).Info("%T-GetWallet is invoked", w)
	defer func() {
		// close transaction
		if errTx := w.readTrx.GormEndTransaction(ctx); errTx != nil {
			w.sugar.WithContext(ctx).Errorf("error when process payload:%v and transaction:%v", err, errTx)
		}
		w.sugar.WithContext(ctx).Info("%T-GetWallet executed", w)
	}()
	ctx = w.util.GetCorrelationIdFromGrpc(ctx)

	// processing payload
	var walletDetail entity.Wallet
	w.sugar.WithContext(ctx).Infof("transforming payload from proto to entity: %v", wallet.GetId())
	if err = w.util.ObjectMapper(wallet, &walletDetail); err != nil {
		walletResp = &pb.WalletResponse{
			WalletDetail: &pb.Wallet{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: model.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when transforming input to entity:%v", err.Error())
		return walletResp, err
	}
	// calling service
	w.sugar.WithContext(ctx).Infof("processing service: %v", wallet.GetId())
	if errMsg := w.walletSvc.GetWalletDetail(ctx, &walletDetail); errMsg.Error != nil {
		walletResp = &pb.WalletResponse{
			WalletDetail: &pb.Wallet{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     errMsg.Error.Error(),
				ErrorType: errMsg.ErrorType.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when calling service:%v", errMsg.Error.Error())
		return walletResp, err
	}

	w.sugar.WithContext(ctx).Infof("transforming payload from entity to proto: %v", wallet.GetId())
	if err = w.util.ObjectMapper(&walletDetail, wallet); err != nil {
		walletResp = &pb.WalletResponse{
			WalletDetail: &pb.Wallet{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: model.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when transforming input to entity:%v", err.Error())
		return walletResp, err
	}
	walletResp = &pb.WalletResponse{
		WalletDetail: wallet,
		ErrorMessage: &pb.ErrorMessage{},
	}
	return walletResp, err
}
