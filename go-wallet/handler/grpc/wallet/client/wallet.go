//go:generate mockgen -source wallet.go -destination mock/wallet_mock.go -package mock

package walletclient

import (
	"github.com/Calmantara/go-common/logger"

	grpcclient "github.com/Calmantara/go-common/infra/grpc/client"
	pb "github.com/Calmantara/go-common/pb"
)

type WalletClient interface {
	// get client registration for grpc
	GetClient() pb.WalletServiceClient
}

type WalletClientImpl struct {
	sugar  logger.CustomLogger
	client pb.WalletServiceClient
}

// constructor function to build new object
func NewWalletClient(sugar logger.CustomLogger, clientConnection grpcclient.GRPCClient) WalletClient {
	// create new connection for client
	client := pb.NewWalletServiceClient(clientConnection.GetClient())
	return &WalletClientImpl{sugar: sugar, client: client}
}

// get client registration for grpc
func (c *WalletClientImpl) GetClient() pb.WalletServiceClient {
	c.sugar.Logger().Info("creating grpc connection for wallet client")
	return c.client
}
