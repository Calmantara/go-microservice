//go:generate mockgen -source balance.go -destination mock/balance_mock.go -package mock
package balance

import (
	"context"

	"github.com/Calmantara/go-wallet/entity"
	"github.com/Calmantara/go-wallet/model"
	"github.com/lovoo/goka"

	cmodel "github.com/Calmantara/go-common/model"
)

type BalanceSvc interface {
	GetBalanceDetail(ctx context.Context, balance *model.WalletBalance) (errModel cmodel.ErrorModel)
	InsertBalance(ctx context.Context, balance *entity.Balance) (errModel cmodel.ErrorModel)
	// consumer
	ConsumeKafkaPayload() goka.ProcessCallback
}
