package emitter

import (
	"context"

	"github.com/Calmantara/go-emitter/entity"
)

type EmitterRepo interface {
	InsertEmitter(ctx context.Context, emitterPayload *entity.EmitterPayload) (err error)
}
