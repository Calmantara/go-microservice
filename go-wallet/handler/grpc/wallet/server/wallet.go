//go:generate mockgen -source wallet.go -destination mock/wallet_mock.go -package mock

package walletserver

import (
	"context"

	"github.com/Calmantara/go-common/pb"
)

type WalletServer interface {
	GetWallet(ctx context.Context, wallet *pb.Wallet) (walletResp *pb.WalletResponse, err error)
}
