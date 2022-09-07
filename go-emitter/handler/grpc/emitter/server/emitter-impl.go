package emitterserver

import (
	"context"

	"github.com/Calmantara/go-common/infra/gorm/transaction"
	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/model"
	"github.com/Calmantara/go-common/pb"
	"github.com/Calmantara/go-emitter/entity"
	"github.com/Calmantara/go-emitter/service/emitter"

	grpcserver "github.com/Calmantara/go-common/infra/grpc/server"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"
)

type EmitterServerImpl struct {
	sugar      logger.CustomLogger
	writeTrx   transaction.Transaction
	server     grpcserver.GRPCServer
	util       serviceutil.UtilService
	assert     serviceassert.Assert
	emitterSvc emitter.EmitterSvc
	pb.UnimplementedEmitterServiceServer
}

func NewEmitterServer(sugar logger.CustomLogger,
	writeTrx transaction.Transaction,
	server grpcserver.GRPCServer,
	util serviceutil.UtilService,
	assert serviceassert.Assert,
	emitterSvc emitter.EmitterSvc) EmitterServer {

	emitterServer := &EmitterServerImpl{
		sugar:      sugar,
		writeTrx:   writeTrx,
		util:       util,
		assert:     assert,
		emitterSvc: emitterSvc,
		server:     server,
	}
	pb.RegisterEmitterServiceServer(server.GetServer(), emitterServer)
	return emitterServer
}

func (w *EmitterServerImpl) SendEmitterPayload(ctx context.Context, emitter *pb.Emitter) (emitterResp *pb.EmitterResponse, err error) {
	// init response
	emitterResp = &pb.EmitterResponse{}
	// set transaction
	ctx = context.WithValue(ctx, transaction.TRANSACTION_KEY, w.writeTrx.GormBeginTransaction(ctx))
	w.sugar.WithContext(ctx).Info("%T-GetEmitter is invoked", w)
	defer func() {
		// close transaction
		if errTx := w.writeTrx.GormEndTransaction(ctx); errTx != nil {
			w.sugar.WithContext(ctx).Errorf("error when process payload:%v and transaction:%v", err, errTx)
		}
		w.sugar.WithContext(ctx).Info("%T-GetEmitter executed", w)
	}()
	ctx = w.util.GetCorrelationIdFromGrpc(ctx)

	// processing payload
	var emitterDetail entity.EmitterPayload
	w.sugar.WithContext(ctx).Infof("transforming payload from proto to entity from: %v", emitter.GetIssuer())
	if err = w.util.ObjectMapper(emitter, &emitterDetail); err != nil {
		emitterResp = &pb.EmitterResponse{
			EmitterDetail: &pb.Emitter{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: model.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when transforming input to entity:%v", err.Error())
		return emitterResp, err
	}
	// calling service
	w.sugar.WithContext(ctx).Infof("processing service from: %v with topic:%v", emitter.GetIssuer(), emitter.GetTopic())
	if errMsg := w.emitterSvc.ProcessEmitterPayload(ctx, &emitterDetail); errMsg.Error != nil {
		emitterResp = &pb.EmitterResponse{
			EmitterDetail: &pb.Emitter{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: model.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when calling service:%v", errMsg.Error.Error())
		return emitterResp, errMsg.Error
	}

	w.sugar.WithContext(ctx).Infof("transforming payload from entity to proto from: %v with topic:%v", emitter.GetIssuer(), emitter.GetTopic())
	if err = w.util.ObjectMapper(&emitterDetail, emitter); err != nil {
		emitterResp = &pb.EmitterResponse{
			EmitterDetail: &pb.Emitter{},
			ErrorMessage: &pb.ErrorMessage{
				Error:     err.Error(),
				ErrorType: model.ERR_INTERNAL_TYPE.String()},
		}
		w.sugar.WithContext(ctx).Errorf("error when transforming input to entity:%v", err.Error())
		return emitterResp, err
	}
	emitterResp = &pb.EmitterResponse{
		EmitterDetail: emitter,
		ErrorMessage:  &pb.ErrorMessage{},
	}
	return emitterResp, err
}
