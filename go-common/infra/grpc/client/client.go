//go:generate mockgen -source client.go -destination mock/client_mock.go -package mock

package grpcclient

import (
	"context"

	"github.com/Calmantara/go-common/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient interface {
	// return current grpc client connection
	GetClient() (con *grpc.ClientConn)
}

type GRPCClientImpl struct {
	sugar      logger.CustomLogger
	host       string
	connection *grpc.ClientConn
}

func NewGRPCClientConnection(sugar logger.CustomLogger, host string) GRPCClient {
	// trying to connect to server first
	clientConnection, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		sugar.WithContext(context.Background()).Errorf("cannot connect: %v", err.Error())
	}

	return &GRPCClientImpl{
		sugar:      sugar,
		host:       host,
		connection: clientConnection,
	}
}

// return current grpc client connection
func (g *GRPCClientImpl) GetClient() (con *grpc.ClientConn) {
	return g.connection
}
