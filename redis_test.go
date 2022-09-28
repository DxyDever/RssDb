package idb

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRedisCluster(t *testing.T) {

	err := InitDBInfoManager_default()
	assert.NoError(t, err)

	redisConfig := RedisConfig{DBConfig{
		Host:         "10.141.2.209",
		Port:         7004,
		Password:     "shumeiShumei2018",
		ConnTimeout:  20,
		ReadTimeout:  20,
		WriteTimeout: 20,
		DBType:       "redis-cluster",
	}}

	err = redisConfig.Init()
	assert.NoError(t, err)

	client, err := redisConfig.OpenDB()

	res, err := client.Do(context.Background(), "SET", "test_key", "test_value").Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", res)

	res, err = client.Do(context.Background(), "GET", "test_key").Result()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", res)

	res, err = client.Get(context.Background(), "test_key").Result()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", res)

}

func TestRedis(t *testing.T) {

	redisConfig := RedisConfig{DBConfig{
		Host:         "10.66.121.171",
		Port:         6379,
		Password:     "shumei123",
		ConnTimeout:  20,
		ReadTimeout:  20,
		WriteTimeout: 20,
		DBType:       "redis",
	}}

	err := redisConfig.Init()
	assert.NoError(t, err)

	client, err := redisConfig.OpenDB()

	res, err := client.Do(context.Background(), "SET", "test_key", "test_value").Result()
	assert.NoError(t, err)
	assert.Equal(t, "OK", res)

	res, err = client.Do(context.Background(), "GET", "test_key").Result()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", res)

	res, err = client.Get(context.Background(), "test_key").Result()
	assert.NoError(t, err)
	assert.Equal(t, "test_value", res)

}
