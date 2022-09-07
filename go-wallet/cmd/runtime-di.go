package main

import (
	"encoding/json"
	"fmt"

	"github.com/Calmantara/go-common/infra/gorm/transaction"
	"github.com/Calmantara/go-common/infra/kafka/consumer"
	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-wallet/repository/balance"
	"github.com/Calmantara/go-wallet/repository/wallet"
	"go.uber.org/dig"

	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"

	configgorm "github.com/Calmantara/go-common/infra/gorm"
	grpcserver "github.com/Calmantara/go-common/infra/grpc/server"
	redisconfig "github.com/Calmantara/go-common/infra/redis"
	redisservice "github.com/Calmantara/go-common/service/redis"
	balanceserver "github.com/Calmantara/go-wallet/handler/grpc/balance/server"
	walletserver "github.com/Calmantara/go-wallet/handler/grpc/wallet/server"
	balancesvc "github.com/Calmantara/go-wallet/service/balance"
	walletsvc "github.com/Calmantara/go-wallet/service/wallet"
)

// initiate all grouped DI
func commonDependencies() []any {
	return []any{logger.NewCustomLogger, config.NewConfigSetup,
		serviceutil.NewUtilService, serviceassert.NewAssert,
		redisconfig.NewRedisConfig, redisservice.NewRedisService}
}

func svcDependencies() []any {
	return []any{walletsvc.NewWalletSvc, balancesvc.NewBalanceSvc}
}

func handlerDependencies() []any {
	return []any{walletserver.NewWalletServer, balanceserver.NewBalanceServer}

}

func BuildKafkaWorker(sugar logger.CustomLogger, conf config.ConfigSetup, balanceSvc balancesvc.BalanceSvc) consumer.Consumer {
	// getting all config
	balanceTransaction := consumer.KafkaWorker{}
	conf.GetConfig("balancetransaction", &balanceTransaction)
	balanceTransaction.Method = balanceSvc.ConsumeKafkaPayload()

	return consumer.NewKafkaConsumer(sugar, conf, balanceTransaction)
}

func BuildRepoDependencies(sugar logger.CustomLogger, conf config.ConfigSetup) (transaction.Transaction, wallet.WalletRepo, balance.BalanceRepo) {
	readCln := configgorm.NewPostgresConfig(sugar, conf, configgorm.WithPostgresMode("read"))
	writeCln := configgorm.NewPostgresConfig(sugar, conf, configgorm.WithPostgresMode("write"))
	return transaction.NewTransaction(sugar, readCln), wallet.NewWalletRepo(sugar, readCln, writeCln), balance.NewBalanceRepo(sugar, readCln, writeCln)
}

func BuildServerDependencies(sugar logger.CustomLogger, conf config.ConfigSetup) grpcserver.GRPCServer {
	var config map[string]any
	conf.GetConfig("service", &config)
	return grpcserver.NewGRPCServer(sugar, fmt.Sprintf("%v", config["grpcport"]))
}

func BuildInRuntime() (serviceConf map[string]any, grpcServer grpcserver.GRPCServer, err error) {
	c := dig.New()
	// define all generic
	var constructor []any
	constructor = append(constructor, BuildRepoDependencies, BuildServerDependencies, BuildKafkaWorker)
	constructor = append(constructor, commonDependencies()...)
	constructor = append(constructor, svcDependencies()...)
	constructor = append(constructor, handlerDependencies()...)

	// provide all generic
	for _, service := range constructor {
		if err := c.Provide(service); err != nil {
			return nil, nil, err
		}
	}
	if err = c.Invoke(func(config config.ConfigSetup, grpc grpcserver.GRPCServer,
		walletServer walletserver.WalletServer, balanceServer balanceserver.BalanceServer,
		kafkaConsumer consumer.Consumer, blcSvc balancesvc.BalanceSvc) {
		// start consumer
		go kafkaConsumer.Consume()
		// service information
		app, _ := json.Marshal(config.GetRawConfig()["service"])
		// init grpc server
		json.Unmarshal(app, &serviceConf)
		grpcServer = grpc
	}); err != nil {
		panic(err)
	}
	return serviceConf, grpcServer, err
}
