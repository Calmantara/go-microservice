package main

import (
	"encoding/json"

	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/middleware/cors"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-gateway/router/v1/balance"
	"github.com/Calmantara/go-gateway/router/v1/wallet"
	"go.uber.org/dig"

	ginrouter "github.com/Calmantara/go-common/infra/gin/router"
	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"
	balancehdl "github.com/Calmantara/go-gateway/handler/http/balance"
	wallethdl "github.com/Calmantara/go-gateway/handler/http/wallet"
	balancev2 "github.com/Calmantara/go-gateway/router/v2/balance"
)

// initiate all grouped DI
func commonDependencies() []any {
	return []any{logger.NewCustomLogger, config.NewConfigSetup,
		serviceutil.NewUtilService, serviceassert.NewAssert,
		ginrouter.NewGinRouter}
}

func hdlDependencies() []any {
	return []any{wallethdl.NewWalletHdl, balancehdl.NewBalanceHdl}
}

func routerDependencies() []any {
	return []any{wallet.NewWalletRouter, balance.NewBalanceRouter, balancev2.NewBalanceRouter}
}

func BuildInRuntime() (serviceConf map[string]any, ginRouter ginrouter.GinRouter, err error) {
	c := dig.New()
	// define all generic
	var constructor []any
	constructor = append(constructor, commonDependencies()...)
	constructor = append(constructor, routerDependencies()...)
	constructor = append(constructor, hdlDependencies()...)

	// provide all generic
	for _, service := range constructor {
		if err := c.Provide(service); err != nil {
			return nil, nil, err
		}
	}
	if err = c.Invoke(func(config config.ConfigSetup, http ginrouter.GinRouter,
		walletR wallet.WalletRouter, balanceR balance.BalanceRouter, balanceRv2 balancev2.BalanceRouter) {
		// service information
		app, _ := json.Marshal(config.GetRawConfig()["service"])
		// init http server
		json.Unmarshal(app, &serviceConf)
		ginRouter = http
		// init cors
		http.USE(cors.NewCorsMiddleware().Cors)
		walletR.Routers()
		balanceR.Routers()
		balanceRv2.Routers()
	}); err != nil {
		panic(err)
	}
	return serviceConf, ginRouter, err
}
