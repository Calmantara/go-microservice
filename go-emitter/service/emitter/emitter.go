package emitter

import (
	"context"

	"github.com/Calmantara/go-common/model"
	"github.com/Calmantara/go-emitter/entity"
)

type EmitterSvc interface {
	ProcessEmitterPayload(ctx context.Context, emitterPayload *entity.EmitterPayload) (errMsg model.ErrorModel)
}
