package balance

import (
	gingroup "github.com/Calmantara/go-common/infra/gin/group"
	ginrouter "github.com/Calmantara/go-common/infra/gin/router"
	"github.com/Calmantara/go-gateway/handler/http/balance"
)

type BalanceRouter interface {
	Routers()
}

type BalanceRouterImpl struct {
	group      gingroup.GinGroup
	balanceHdl balance.BalanceHdl
}

func NewBalanceRouter(ginRouter ginrouter.GinRouter, balanceHdl balance.BalanceHdl) BalanceRouter {
	group := ginRouter.GROUP("api/v2/balance")
	return &BalanceRouterImpl{group: group, balanceHdl: balanceHdl}
}

func (w *BalanceRouterImpl) get() {
	w.group.GET("", w.balanceHdl.GetBalanceDetail)
}

func (w *BalanceRouterImpl) Routers() {
	w.get()
}
