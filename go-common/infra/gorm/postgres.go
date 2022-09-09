//go:generate mockgen -source postgres.go -destination mock/postgres_mock.go -package mock

package configgorm

import (
	"context"
	"fmt"
	"time"

	"github.com/Calmantara/go-common/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "github.com/Calmantara/go-common/setup/config"
	pglog "gorm.io/gorm/logger"
)

type GormKey string

const (
	TRANSACTION_KEY GormKey = "GORM_TRANSACTION"
)

type PostgresParam struct {
	Host                 string `json:"host"`
	Database             string `json:"database"`
	Username             string `json:"username"`
	Password             string `json:"password"`
	Timezone             string `json:"timezone"`
	Mode                 string `json:"mode"`
	Automigrate          bool   `json:"automigrate"`
	EnableLog            bool   `json:"enablelog"`
	Port                 int    `json:"port"`
	MaxConnection        int    `json:"maxconnection"`
	MaxIdleConnection    int    `json:"maxidleconnection"`
	MaxIdleConnectionTtl int    `json:"maxidleconnectionttl"`
}

type PostgresConfig interface {
	GetClient() *gorm.DB
	GetParam() PostgresParam
	GenerateTransaction(ctx context.Context) (tx *gorm.DB)
}

type PostgresConfigImpl struct {
	sugar         logger.CustomLogger
	postgresParam PostgresParam
	gormDB        *gorm.DB
}

type Option func(*PostgresParam)

func NewPostgresConfig(sugar logger.CustomLogger, config config.ConfigSetup, ops ...Option) PostgresConfig {
	sugar.Logger().Info("Initialization Postgres Configuration. . .")
	postgresParam := PostgresParam{Mode: "read"}

	//iterate all option function
	for _, v := range ops {
		v(&postgresParam)
	}
	//get config
	mode := "postgres" + postgresParam.Mode
	config.GetConfig(mode, &postgresParam)

	gormDB := openPostgresConnection(postgresParam)
	pg := PostgresConfigImpl{
		sugar:         sugar,
		gormDB:        gormDB,
		postgresParam: postgresParam,
	}

	sugar.Logger().Infof("database connection error %v:%v", mode, gormDB.Error)
	return &pg
}

func openPostgresConnection(param PostgresParam) (gormDB *gorm.DB) {
	// Logger enabled
	pgLog := pglog.Discard
	if param.EnableLog {
		pgLog = pglog.Default
	}

	// open postgres connection
	gormDB, err := gorm.Open(postgres.Open(param.postgresConfigString()), &gorm.Config{Logger: pgLog})
	if err != nil {
		panic(err)
	}

	// setting connection limit
	db, err := gormDB.DB()
	if err != nil {
		panic(err)
	}

	// max connection
	if param.MaxConnection == 0 {
		param.MaxConnection = 50
	}

	if param.MaxIdleConnection == 0 {
		param.MaxIdleConnection = 1
	}

	db.SetMaxIdleConns(param.MaxIdleConnection)
	db.SetMaxOpenConns(param.MaxConnection)
	db.SetConnMaxIdleTime(time.Duration(param.MaxIdleConnectionTtl * int(time.Minute)))
	// finish setup
	return gormDB
}

// postgres config methods
func (config *PostgresConfigImpl) GetClient() *gorm.DB {
	return config.gormDB
}

func (config *PostgresConfigImpl) GetParam() PostgresParam {
	return config.postgresParam
}

func (config *PostgresConfigImpl) GenerateTransaction(ctx context.Context) (tx *gorm.DB) {
	val := ctx.Value((TRANSACTION_KEY))
	if val == nil {
		val = ctx.Value(string(TRANSACTION_KEY))
	}

	if val != nil {
		tx = val.(*gorm.DB)
	} else {
		tx = config.GetClient().WithContext(ctx)
	}
	return tx
}

// Options Function
func WithPostgresMode(mode string) Option {
	return func(rc *PostgresParam) { rc.Mode = mode }
}

func ActiveRecordQuery() string {
	return `record_flag = 'ACTIVE'`
}

func (p *PostgresParam) postgresConfigString() string {
	return fmt.Sprintf("host=%s port=%v user=%s password=%s dbname=%s sslmode=disable TimeZone=%v",
		p.Host,
		p.Port,
		p.Username,
		p.Password,
		p.Database,
		p.Timezone,
	)
}
