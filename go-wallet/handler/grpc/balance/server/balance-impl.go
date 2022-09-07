package balanceserver

import (
	"context"

	"github.com/Calmantara/go-common/infra/gorm/transaction"
	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/pb"
	"github.com/Calmantara/go-wallet/model"
	"github.com/Calmantara/go-wallet/service/balance"

	grpcserver "github.com/Calmantara/go-common/infra/grpc/server"
	cmodel "github.com/Calmantara/go-common/model"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"
)

type BalanceServerImpl struct {
	sugar      logger.CustomLogger
	readTrx    transaction.Transaction
	server     grpcserver.GRPCServer
	util       serviceutil.UtilService
	assert     serviceassert.Assert
	balanceSvc balance.BalanceSvc
	pb.UnimplementedBalanceServiceServer
}

func NewBalanceServer(sugar logger.CustomLogger,
	readTrx transaction.Transaction,
	server grpcserver.GRPCServer,
	util serviceutil.UtilService,
	assert serviceassert.Assert,
	balanceSvc balance.BalanceSvc) BalanceServer {

	balanceServer := &BalanceServerImpl{
		sugar:      sugar,
		readTrx:    readTrx,
		util:       util,
		assert:     assert,
		balanceSvc: balanceSvc,
	}
	pb.RegisterBalanceServiceServer(server.GetServer(), balanceServer)
	return balanceServer
}

func (w *BalanceServerImpl) GetBalance(ctx context.Context, wallet *pb.Wallet) (balanceResp *pb.BalanceResponse, err error) {
	// init response
	balanceResp = &pb.BalanceResponse{}
	// set transaction
	ctx = context.WithValue(ctx, transaction.TRANSACTION_KEY, w.readTrx.GormBeginTransaction(ctx))
	w.sugar.WithContext(ctx).Info("%T-GetBalance is invoked", w)
	defer func() {
		// close transaction
		if errTx := w.readTrx.GormEndTransaction(ctx); errTx != nil {
			w.sugar.WithContext(ctx).Errorf("error when process payload:%v and transaction:%v", err, errTx)
		}
		w.sugar.WithContext(ctx).Info("%T-GetBalance executed", w)
	}()
	ctx = w.util.GetCorrelationIdFromGrpc(ctx)

	// processing payload
	balanceDetail := model.WalletBalance{WalletId: uint64(wallet.GetId())}
	w.sugar.WithContext(ctx).Infof("processing service: %v", wallet.GetId())
	if errMsg := w.balanceSvc.GetBalanceDetail(ctx, &balanceDetail); errMsg.Error != nil {
		balanceResp = &pb.BalanceResponse{
			BalanceDetail: &pb.Balance{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: cmodel.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when calling service:%v", errMsg.Error.Error())
		return balanceResp, errMsg.Error
	}

	w.sugar.WithContext(ctx).Infof("transforming payload from entity to proto: %v", wallet.GetId())
	balance := &pb.Balance{}
	if err = w.util.ObjectMapper(&balanceDetail, balance); err != nil {
		balanceResp = &pb.BalanceResponse{
			BalanceDetail: &pb.Balance{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: cmodel.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when transforming input to entity:%v", err.Error())
		return balanceResp, err
	}
	balanceResp = &pb.BalanceResponse{
		BalanceDetail: balance,
		ErrorMessage:  &pb.ErrorMessage{},
	}
	return balanceResp, err
}

func (w *BalanceServerImpl) GetBalanceByTtl(ctx context.Context, wallet *pb.Wallet) (balanceResp *pb.BalanceResponse, err error) {
	// init response
	balanceResp = &pb.BalanceResponse{}
	// set transaction
	ctx = context.WithValue(ctx, transaction.TRANSACTION_KEY, w.readTrx.GormBeginTransaction(ctx))
	w.sugar.WithContext(ctx).Info("%T-GetBalanceByTtl is invoked", w)
	defer func() {
		// close transaction
		if errTx := w.readTrx.GormEndTransaction(ctx); errTx != nil {
			w.sugar.WithContext(ctx).Errorf("error when process payload:%v and transaction:%v", err, errTx)
		}
		w.sugar.WithContext(ctx).Info("%T-GetBalanceByTtl executed", w)
	}()
	ctx = w.util.GetCorrelationIdFromGrpc(ctx)

	// processing payload
	balanceDetail := model.WalletBalance{WalletId: uint64(wallet.GetId())}
	w.sugar.WithContext(ctx).Infof("processing service: %v", wallet.GetId())
	if errMsg := w.balanceSvc.GetBalanceDetailByTtl(ctx, &balanceDetail); errMsg.Error != nil {
		balanceResp = &pb.BalanceResponse{
			BalanceDetail: &pb.Balance{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: cmodel.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when calling service:%v", errMsg.Error.Error())
		return balanceResp, errMsg.Error
	}

	w.sugar.WithContext(ctx).Infof("transforming payload from entity to proto: %v", wallet.GetId())
	balance := &pb.Balance{}
	if err = w.util.ObjectMapper(&balanceDetail, balance); err != nil {
		balanceResp = &pb.BalanceResponse{
			BalanceDetail: &pb.Balance{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: cmodel.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when transforming input to entity:%v", err.Error())
		return balanceResp, err
	}
	balanceResp = &pb.BalanceResponse{
		BalanceDetail: balance,
		ErrorMessage:  &pb.ErrorMessage{},
	}
	return balanceResp, err
}
