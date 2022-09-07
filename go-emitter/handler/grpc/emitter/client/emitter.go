//go:generate mockgen -source emitter.go -destination mock/emitter_mock.go -package mock

package emitterclient

import (
	"github.com/Calmantara/go-common/logger"

	grpcclient "github.com/Calmantara/go-common/infra/grpc/client"
	pb "github.com/Calmantara/go-common/pb"
)

type EmitterClient interface {
	// get client registration for grpc
	GetClient() pb.EmitterServiceClient
}

type EmitterClientImpl struct {
	sugar  logger.CustomLogger
	client pb.EmitterServiceClient
}

// constructor function to build new object
func NewEmitterClient(sugar logger.CustomLogger, clientConnection grpcclient.GRPCClient) EmitterClient {
	// create new connection for client
	client := pb.NewEmitterServiceClient(clientConnection.GetClient())
	return &EmitterClientImpl{sugar: sugar, client: client}
}

// get client registration for grpc
func (c *EmitterClientImpl) GetClient() pb.EmitterServiceClient {
	c.sugar.Logger().Info("creating grpc connection for emitter client")
	return c.client
}
