//go:generate mockgen -source redis.go -destination mock/redis_mock.go -package mock

package redisservice

import (
	"context"
	"time"
)

type RedisService interface {
	Get(ctx context.Context, key RedisKey, model interface{}) error
	Set(ctx context.Context, key RedisKey, model interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key RedisKey) error
}
