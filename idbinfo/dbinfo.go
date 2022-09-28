package idbinfo

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_TYPE_MYSQL        = "mysql"
	DB_TYPE_REDIS        = "redis"
	DB_TYPE_REDISCLUSTER = "redis-cluster"
)

func generateDBKey(host, port, dbName string) string {
	key := net.JoinHostPort(host, port)
	if dbName != "" {
		key += "/" + dbName
	}
	return key
}

type DBInfo struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
	DBType   string `json:"db_type"`
	FlagSsl  bool   `json:"flag_ssl"`
}

func (this *DBInfo) GetAddr() []string {
	if this.DBType == DB_TYPE_REDISCLUSTER {
		return []string{this.JoinHostPort(), this.JoinHostPort()}
	}
	return []string{this.JoinHostPort()}
}

func (this *DBInfo) JoinHostPort() string {
	return net.JoinHostPort(this.Host, this.Port)
}

func (this *DBInfo) String() string {
	bts, _ := json.Marshal(this)
	return string(bts)
}

func (this *DBInfo) Ping() bool {
	if this.DBType == DB_TYPE_MYSQL {
		return this.PingMysql()
	}
	if this.DBType == DB_TYPE_REDIS || this.DBType == DB_TYPE_REDISCLUSTER {
		return this.PingRedis()
	}
	return false
}

func (this *DBInfo) ISmPing() bool {
	if this.DBType == DB_TYPE_MYSQL {
		return this.ISmPingMysql()
	}
	if this.DBType == DB_TYPE_REDIS || this.DBType == DB_TYPE_REDISCLUSTER {
		return this.ISmPingRedis()
	}
	return false
}

func (this *DBInfo) ISmPingMysql() bool {
	var (
		db  *sql.DB
		res *sql.Rows
		err error
	)
	defer func() {
		if err != nil {
			passwordSuffix := ""
			if len(this.Password)-3 > 0 {
				passwordSuffix = this.Password[len(this.Password)-3 : len(this.Password)]
			}
			fmt.Fprintf(os.Stderr, "ism db ping error.db host:%s,port:%s,user:%s,passwordSuffix:%s,err:%v\n", this.Host, this.Port, this.UserName, passwordSuffix, err)
		}
	}()
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", this.UserName, this.Password, this.Host, this.Port)
	if this.FlagSsl {
		connectStr = connectStr + "&tls=skip-verify"
	}
	db, err = sql.Open("mysql", connectStr)
	if err != nil {
		return false
	}
	defer db.Close()

	sqlStr := "SELECT VERSION();"
	if res, err = db.Query(sqlStr); err != nil {
		return false
	}
	defer res.Close()
	return true
}

func (this *DBInfo) PingMysql() bool {
	var (
		db  *sql.DB
		err error
	)
	defer func() {
		if err != nil {
			passwordSuffix := ""
			if len(this.Password)-3 > 0 {
				passwordSuffix = this.Password[len(this.Password)-3 : len(this.Password)]
			}
			fmt.Fprintf(os.Stderr, "db ping error.db host:%s,port:%s,user:%s,passwordSuffix:%s,err:%v\n", this.Host, this.Port, this.UserName, passwordSuffix, err)
		}
	}()
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", this.UserName, this.Password, this.Host, this.Port)
	if this.FlagSsl {
		connectStr = connectStr + "&tls=skip-verify"
	}
	db, err = sql.Open("mysql", connectStr)
	if err != nil {
		return false
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		return false
	}
	return true
}

func (this *DBInfo) ISmPingRedis() bool {
	var err error
	defer func() {
		if err != nil {
			passwordSuffix := ""
			if len(this.Password)-3 > 0 {
				passwordSuffix = this.Password[len(this.Password)-3 : len(this.Password)]
			}
			fmt.Fprintf(os.Stderr, "ism db ping error.db host:%s,port:%s,user:%s,passwordSuffix:%s,err:%v\n", this.Host, this.Port, this.UserName, passwordSuffix, err)
		}
	}()
	redisConfig := &redis.UniversalOptions{
		Addrs:        this.GetAddr(),
		Username:     this.UserName,
		Password:     this.Password,
		DialTimeout:  time.Duration(20) * time.Millisecond,
		ReadTimeout:  time.Duration(20) * time.Millisecond,
		WriteTimeout: time.Duration(20) * time.Millisecond,
	}
	if this.FlagSsl {
		redisConfig.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	client := redis.NewUniversalClient(redisConfig)
	defer func() {
		client.Close()
		if err != nil {
			passwordSuffix := ""
			if len(this.Password)-3 > 0 {
				passwordSuffix = this.Password[len(this.Password)-3 : len(this.Password)]
			}
			fmt.Fprintf(os.Stderr, "db ping error.db host:%s,port:%s,user:%s,passwordSuffix:%s,err:%v\n", this.Host, this.Port, this.UserName, passwordSuffix, err)
		}
	}()
	if _, err = client.Do(context.Background(), "info").Result(); err != nil {
		return false
	}
	return true
}

func (this *DBInfo) PingRedis() bool {

	var err error
	defer func() {
		if err != nil {
			passwordSuffix := ""
			if len(this.Password)-3 > 0 {
				passwordSuffix = this.Password[len(this.Password)-3 : len(this.Password)]
			}
			fmt.Fprintf(os.Stderr, "db ping error.db host:%s,port:%s,user:%s,passwordSuffix:%s,err:%v\n", this.Host, this.Port, this.UserName, passwordSuffix, err)
		}
	}()
	redisConfig := &redis.UniversalOptions{
		Addrs:        this.GetAddr(),
		Username:     this.UserName,
		Password:     this.Password,
		DialTimeout:  time.Duration(20) * time.Millisecond,
		ReadTimeout:  time.Duration(20) * time.Millisecond,
		WriteTimeout: time.Duration(20) * time.Millisecond,
	}
	if this.FlagSsl {
		redisConfig.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}
	client := redis.NewUniversalClient(redisConfig)

	defer func() {
		client.Close()
		if err != nil {
			passwordSuffix := ""
			if len(this.Password)-3 > 0 {
				passwordSuffix = this.Password[len(this.Password)-3 : len(this.Password)]
			}
			fmt.Fprintf(os.Stderr, "db ping error.db host:%s,port:%s,user:%s,passwordSuffix:%s,err:%v\n", this.Host, this.Port, this.UserName, passwordSuffix, err)
		}
	}()
	if _, err = client.Ping(context.Background()).Result(); err != nil {
		return false
	}
	return true
}

type HostInfoSet struct {
	dbListMap map[string][]*DBInfo
	mutex     sync.RWMutex
}

func (this *HostInfoSet) init() {
	this.dbListMap = make(map[string][]*DBInfo, 10)
}

func (this *HostInfoSet) Clear() {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.dbListMap = make(map[string][]*DBInfo, 10)
}

func (this *HostInfoSet) getOneRandomDBInfo(host, port, dbName string) *DBInfo {
	rand.Seed(time.Now().UnixNano())

	this.mutex.RLock()
	defer this.mutex.RUnlock()
	key := generateDBKey(host, port, dbName)
	if dbs, isExist := this.dbListMap[key]; isExist && len(dbs) > 0 {
		index := rand.Intn(len(dbs))
		return dbs[index]
	}
	return nil
}

func (this *HostInfoSet) AddDBInfoList(key string, dbList []*DBInfo) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, isKeyExist := this.dbListMap[key]; isKeyExist {
		this.dbListMap[key] = append(this.dbListMap[key], dbList...)
	} else {
		this.dbListMap[key] = dbList
	}
}

func (this *HostInfoSet) ResetDBInfoList(key string, dbList []*DBInfo) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.dbListMap[key] = dbList
}

func (this *HostInfoSet) GetDBInfoList(key string) []*DBInfo {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.dbListMap[key]
}
