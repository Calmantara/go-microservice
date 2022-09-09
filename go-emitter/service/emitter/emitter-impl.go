package emitter

import (
	"context"
	"errors"

	"github.com/Calmantara/go-common/infra/kafka/producer"
	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/model"
	"github.com/Calmantara/go-emitter/entity"
	"github.com/Calmantara/go-emitter/repository/emitter"

	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"
)

type EmitterSvcImpl struct {
	sugar       logger.CustomLogger
	assert      serviceassert.Assert
	util        serviceutil.UtilService
	emitterRepo emitter.EmitterRepo
	producer    producer.Producer
}

func NewEmitterSvc(sugar logger.CustomLogger, assert serviceassert.Assert, util serviceutil.UtilService, emitterRepo emitter.EmitterRepo, producer producer.Producer) EmitterSvc {
	return &EmitterSvcImpl{sugar: sugar, assert: assert, emitterRepo: emitterRepo, producer: producer, util: util}
}

func (e *EmitterSvcImpl) ProcessEmitterPayload(ctx context.Context, emitterPayload *entity.EmitterPayload) (errMsg model.ErrorModel) {
	e.sugar.WithContext(ctx).Infof("%T-ProcessEmitterPayload is invoked", e)
	defer e.sugar.WithContext(ctx).Infof("%T-ProcessEmitterPayload executed", e)

	// check payload
	if e.assert.IsEmpty(emitterPayload.Issuer) || e.assert.IsEmpty(emitterPayload.Topic.String()) || e.assert.IsEmpty(emitterPayload.Message) {
		errMsg = model.ErrorModel{
			Error:     errors.New("payload is invalid"),
			ErrorType: model.ERR_BAD_REQUEST_TYPE,
		}
		e.sugar.WithContext(ctx).Errorf("invalid payload %v", *emitterPayload)
		return errMsg
	}
	// insert to database to get id with status true
	emitterPayload.Status = true
	if err := e.emitterRepo.InsertEmitter(ctx, emitterPayload); err != nil {
		errMsg = model.ErrorModel{
			Error:     errors.New("error processing payload"),
			ErrorType: model.ERR_INTERNAL_TYPE,
		}
		e.sugar.WithContext(ctx).Errorf("error inserting payload %v", err.Error())
		return errMsg
	}
	// send to kafka
	ctxBg := e.util.ContextBackground(ctx)
	go func(payload entity.EmitterPayload) {
		if err := e.producer.Publish(ctxBg, payload.Topic, payload); err != nil {
			e.sugar.WithContext(ctx).Errorf("error when emitting:%v", err.Error())
			//  update failed status
			payload.Attempt += 1
			payload.Status = false
			e.emitterRepo.InsertEmitter(ctxBg, &payload)
			return
		}
		e.sugar.WithContext(ctx).Infof("success emitting payload with id:%v", payload.Id)
	}(*emitterPayload)
	return errMsg
}
