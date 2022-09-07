//go:generate mockgen -source wallet.go -destination mock/wallet_mock.go -package mock

package wallet

import (
	"context"

	"github.com/Calmantara/go-common/model"
	"github.com/Calmantara/go-wallet/entity"
)

type WalletSvc interface {
	GetWalletDetail(ctx context.Context, wallet *entity.Wallet) (errModel model.ErrorModel)
	UpsertWallet(ctx context.Context, wallet *entity.Wallet) (errModel model.ErrorModel)
}
