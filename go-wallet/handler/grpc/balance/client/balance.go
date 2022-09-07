//go:generate mockgen -source balance.go -destination mock/balance_mock.go -package mock

package balanceclient

import (
	"github.com/Calmantara/go-common/logger"

	grpcclient "github.com/Calmantara/go-common/infra/grpc/client"
	pb "github.com/Calmantara/go-common/pb"
)

type BalanceClient interface {
	// get client registration for grpc
	GetClient() pb.BalanceServiceClient
}

type BalanceClientImpl struct {
	sugar  logger.CustomLogger
	client pb.BalanceServiceClient
}

// constructor function to build new object
func NewBalanceClient(sugar logger.CustomLogger, clientConnection grpcclient.GRPCClient) BalanceClient {
	// create new connection for client
	client := pb.NewBalanceServiceClient(clientConnection.GetClient())
	return &BalanceClientImpl{sugar: sugar, client: client}
}

// get client registration for grpc
func (c *BalanceClientImpl) GetClient() pb.BalanceServiceClient {
	c.sugar.Logger().Info("creating grpc connection for balance client")
	return c.client
}
