//go:generate mockgen -source redis.go -destination mock/redis_mock.go -package mock

package redisconfig

import (
	"context"
	"fmt"
	"log"

	"github.com/Calmantara/go-common/setup/config"
	"github.com/go-redis/redis/v8"
)

type RedisParam struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     int    `json:"port"`
	Database int    `json:"db"`
}
type RedisConfig interface {
	GetClient() *redis.Client
	GetParam() RedisParam
}
type RedisConfigImpl struct {
	redis      *redis.Client
	redisParam RedisParam
}
type Option func(*RedisParam)

func NewRedisConfig(config config.ConfigSetup, ops ...Option) RedisConfig {
	// get config
	var redisConfig RedisParam
	config.GetConfig("redis", &redisConfig)
	//iterate all option function
	for _, v := range ops {
		v(&redisConfig)
	}
	rc := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Database,
	})

	log.Println("Initialization Redis Configuration. . .",
		"ping: "+rc.Ping(context.Background()).String())
	return &RedisConfigImpl{
		redis:      rc,
		redisParam: redisConfig,
	}
}

// Redis Config function
func (rc *RedisConfigImpl) GetClient() *redis.Client {
	return rc.redis
}
func (rc *RedisConfigImpl) GetParam() RedisParam {
	return rc.redisParam
}

// Options Function
func WithRedisHost(host string) Option {
	return func(rc *RedisParam) { rc.Host = host }
}
func WithRedisPort(port int) Option {
	return func(rc *RedisParam) { rc.Port = port }
}
func WithRedisDB(db int) Option {
	return func(rc *RedisParam) { rc.Database = db }
}
func WithRedisPassword(password string) Option {
	return func(rc *RedisParam) { rc.Password = password }
}
