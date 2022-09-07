package redisservice

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redisconfig "github.com/Calmantara/go-common/infra/redis"
)

type RedisKey string

func (r RedisKey) String() string {
	return string(r)
}

func (r RedisKey) Append(str any) RedisKey {
	return RedisKey(fmt.Sprintf("%v%v", string(r), str))
}

type RedisServiceImpl struct {
	redisConfig redisconfig.RedisConfig
}

func NewRedisService(redisConfig redisconfig.RedisConfig) RedisService {
	return &RedisServiceImpl{
		redisConfig: redisConfig,
	}
}

func (c *RedisServiceImpl) Get(ctx context.Context, key RedisKey, model interface{}) error {
	client := c.redisConfig.GetClient()

	//get from redis
	val, err := client.Get(ctx, fmt.Sprintf("%v", key)).Result()
	if err != nil {
		return err
	}

	// get string value
	if val == "" {
		return nil
	}

	// unmarshal string to model
	if err = json.Unmarshal([]byte(val), &model); err != nil {
		return err
	}
	return nil
}

func (c *RedisServiceImpl) Set(ctx context.Context, key RedisKey, model interface{}, ttl time.Duration) error {
	client := c.redisConfig.GetClient()

	// marshal model
	val, err := json.Marshal(&model)
	if err != nil {
		return err
	}

	// set to redis
	status := client.Set(ctx, key.String(), string(val), ttl)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (c *RedisServiceImpl) Delete(ctx context.Context, key RedisKey) error {
	client := c.redisConfig.GetClient()
	iter := client.Scan(ctx, 0, key.String(), 0).Iterator()

	for iter.Next(ctx) {
		err := client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
