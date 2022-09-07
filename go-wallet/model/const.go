package model

import redisservice "github.com/Calmantara/go-common/service/redis"

const (
	WALLET_KEY  redisservice.RedisKey = "WALLET:"
	BALANCE_KEY redisservice.RedisKey = "BALANCE:"
)
