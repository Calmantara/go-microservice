//go:generate mockgen -source balance.go -destination mock/balance_mock.go -package mock

package balanceserver

import (
	"context"

	"github.com/Calmantara/go-common/pb"
)

type BalanceServer interface {
	GetBalance(ctx context.Context, balance *pb.Wallet) (balanceResp *pb.BalanceResponse, err error)
}
