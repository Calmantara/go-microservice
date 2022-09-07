package wallet

import (
	gingroup "github.com/Calmantara/go-common/infra/gin/group"
	ginrouter "github.com/Calmantara/go-common/infra/gin/router"
	"github.com/Calmantara/go-gateway/handler/http/wallet"
)

type WalletRouter interface {
	Routers()
}

type WalletRouterImpl struct {
	group     gingroup.GinGroup
	walletHdl wallet.WalletHdl
}

func NewWalletRouter(ginRouter ginrouter.GinRouter, walletHdl wallet.WalletHdl) WalletRouter {
	group := ginRouter.GROUP("api/v1/wallet")
	return &WalletRouterImpl{group: group, walletHdl: walletHdl}
}

func (w *WalletRouterImpl) get() {
	w.group.GET("", w.walletHdl.GetWalletDetail)
}

func (w *WalletRouterImpl) post() {
	w.group.POST("", w.walletHdl.PostWallet)

}

func (w *WalletRouterImpl) Routers() {
	w.get()
	w.post()
}
