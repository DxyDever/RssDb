package idb

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	redisConnMap    sync.Map
	redisConnMapMtx sync.Mutex
)

func GetRedisConnection(config RedisConfig) (redis.UniversalClient, error) {
	if conn, isExist := redisConnMap.Load(config); isExist {
		return conn.(redis.UniversalClient), nil
	}

	redisConnMapMtx.Lock()
	defer redisConnMapMtx.Unlock()
	if conn, isExist := redisConnMap.Load(config); isExist {
		return conn.(redis.UniversalClient), nil
	}
	conn, err := config.OpenDB()
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("redis connection is nil")
	}
	redisConnMap.Store(config, conn)
	return conn, nil
}

func DeleteRedisConnection(config RedisConfig) {
	if _, isExist := redisConnMap.Load(config); !isExist {
		return
	}
	redisConnMapMtx.Lock()
	defer redisConnMapMtx.Unlock()
	if conn, isExist := redisConnMap.Load(config); isExist {
		redisConnMap.Delete(config)
		conn.(redis.UniversalClient).Close()
	}
}
