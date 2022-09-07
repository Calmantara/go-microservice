//go:generate mockgen -source emitter.go -destination mock/emitter_mock.go -package mock

package emitterserver

import (
	"context"

	"github.com/Calmantara/go-common/pb"
)

type EmitterServer interface {
	SendEmitterPayload(ctx context.Context, emitter *pb.Emitter) (emitterResp *pb.EmitterResponse, err error)
}
