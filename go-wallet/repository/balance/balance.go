//go:generate mockgen -source balance.go -destination mock/balance_mock.go -package mock

package balance

import (
	"context"

	"github.com/Calmantara/go-wallet/entity"
	"github.com/Calmantara/go-wallet/model"
)

type BalanceRepo interface {
	// read
	ReadBalance(ctx context.Context, walletBalance *model.WalletBalance) (err error)
	// write
	InsertBalance(ctx context.Context, balance *entity.Balance) (err error)
}
