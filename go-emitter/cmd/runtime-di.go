package main

import (
	"encoding/json"
	"fmt"

	"github.com/Calmantara/go-common/infra/gorm/transaction"
	"github.com/Calmantara/go-common/infra/kafka/producer"
	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-common/setup/config"
	"github.com/Calmantara/go-emitter/repository/emitter"
	"go.uber.org/dig"

	serviceassert "github.com/Calmantara/go-common/service/assert"
	serviceutil "github.com/Calmantara/go-common/service/util"

	configgorm "github.com/Calmantara/go-common/infra/gorm"
	grpcserver "github.com/Calmantara/go-common/infra/grpc/server"
	emitterserver "github.com/Calmantara/go-emitter/handler/grpc/emitter/server"
	emittersvc "github.com/Calmantara/go-emitter/service/emitter"
)

// initiate all grouped DI
func commonDependencies() []any {
	return []any{logger.NewCustomLogger, config.NewConfigSetup,
		serviceutil.NewUtilService, serviceassert.NewAssert, producer.NewKafkaProducer}
}

func svcDependencies() []any {
	return []any{emittersvc.NewEmitterSvc}
}

func handlerDependencies() []any {
	return []any{emitterserver.NewEmitterServer}

}

func BuildRepoDependencies(sugar logger.CustomLogger, conf config.ConfigSetup) (transaction.Transaction, emitter.EmitterRepo) {
	readCln := configgorm.NewPostgresConfig(sugar, conf, configgorm.WithPostgresMode("read"))
	writeCln := configgorm.NewPostgresConfig(sugar, conf, configgorm.WithPostgresMode("write"))
	return transaction.NewTransaction(sugar, readCln), emitter.NewEmitterRepo(sugar, readCln, writeCln)
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
	constructor = append(constructor, BuildRepoDependencies, BuildServerDependencies)
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
		emitterServer emitterserver.EmitterServer) {
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
