package idb

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/DxyDever/RssDb/idbinfo"
	"github.com/go-redis/redis/v8"
)

// RedisConfig /*
/*
redis 配置示例：
{
    "Host":"10.66.121.171",
    "Port":6379,
    "Password":"shumei123",
    "ConnTimeout":1000,
    "ReadTimeout":1000,
    "WriteTimeout":1000,
    "MaxIdle":50,
    "MaxConn":200,
    "DBType":"redis-cluster"
}
*/
type RedisConfig struct {
	DBConfig
}

func (this *RedisConfig) Init() error {
	if len(this.DBType) == 0 {
		this.DBType = DB_TYPE_REDIS
	}
	return this.DBConfig.Init()
}

func (this *RedisConfig) OpenDB() (redis.UniversalClient, error) {
	dbInfo := idbinfo.GetOneRandomDBInfo(this.Host, strconv.Itoa(this.Port), this.DbName)
	if dbInfo == nil {
		return nil, fmt.Errorf("dbinfo:%s, err:not get dbInfo", this.DBConfig.JoinHostPort())
	}
	redisConfig := &redis.UniversalOptions{
		Addrs:            dbInfo.GetAddr(),
		Username:         dbInfo.UserName,
		Password:         dbInfo.Password,
		SentinelPassword: dbInfo.Password,
		MaxRetries:       this.RetryCount,
		DialTimeout:      time.Duration(this.ConnTimeout) * time.Millisecond,
		ReadTimeout:      time.Duration(this.ReadTimeout) * time.Millisecond,
		WriteTimeout:     time.Duration(this.WriteTimeout) * time.Millisecond,
		MasterName:       this.MasterName,
	}
	if dbInfo.FlagSsl {
		redisConfig.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	redisClient := redis.NewUniversalClient(redisConfig)
	return redisClient, nil
}
