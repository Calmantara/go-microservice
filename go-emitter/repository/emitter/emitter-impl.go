package emitter

import (
	"context"

	configgorm "github.com/Calmantara/go-common/infra/gorm"
	"github.com/Calmantara/go-common/logger"
	"github.com/Calmantara/go-emitter/entity"
	"gorm.io/gorm/clause"
)

type EmitterRepoImpl struct {
	sugar    logger.CustomLogger
	readCln  configgorm.PostgresConfig
	writeCln configgorm.PostgresConfig
}

func NewEmitterRepo(sugar logger.CustomLogger, readCln configgorm.PostgresConfig, writeCln configgorm.PostgresConfig) EmitterRepo {
	// read config and decide migrator
	emitter := &EmitterRepoImpl{sugar: sugar, writeCln: writeCln}
	if readCln.GetParam().Automigrate {
		sugar.Logger().Info("automigrate invoked for emitter")
		cln := readCln.GetClient()
		cln.AutoMigrate(entity.EmitterPayload{})
	}
	sugar.Logger().Info("init emitter repo")
	return emitter

}

func (e *EmitterRepoImpl) InsertEmitter(ctx context.Context, emitterPayload *entity.EmitterPayload) (err error) {
	e.sugar.WithContext(ctx).Infof("%T-InsertEmitter is invoked", e)
	// generate transaction
	txn := e.writeCln.GenerateTransaction(ctx)
	txn.Model(entity.EmitterPayload{}).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"status"}),
		}).
		Create(emitterPayload)
	if err = txn.Error; err != nil {
		e.sugar.WithContext(ctx).Errorf("error execute InsertEmitter:%v", err.Error())
	}
	e.sugar.WithContext(ctx).Infof("%T-InsertEmitter executed", e)
	return err
}
