//go:generate mockgen -source wallet.go -destination mock/wallet_mock.go -package mock

package wallet

import (
	"context"

	"github.com/Calmantara/go-wallet/entity"
)

type WalletRepo interface {
	// read
	ReadWallet(ctx context.Context, wallet *entity.Wallet) (err error)
	// write
	UpsertWallet(ctx context.Context, wallet *entity.Wallet) (err error)
}
