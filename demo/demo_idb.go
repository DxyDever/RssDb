package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/DxyDever/RssDb"
)

// export GOPATH=/your/project
// go run $GOPATH/src/code.aliyun.com/module-go/idb/demo/demo_idb.go -idbinfo_config_file=$GOPATH/src/code.aliyun.com/module-go/idb/demo/config/idbinfo.json -mysql_config_file=$GOPATH/src/code.aliyun.com/module-go/idb/demo/config/mysql_config.json -redis_config_file=$GOPATH/src/code.aliyun.com/module-go/idb/demo/config/redis_config.json

var (
	mysqlConfigFile string
	redisConfigFile string

	// 实际使用中，将 idb.MysqlConfig、idb.RedisConfig 放到自己的配置文件结构体中
	mysqlConfig idb.MysqlConfig
	redisConfig idb.RedisConfig
)

func init() {
	flag.StringVar(&mysqlConfigFile, "mysql_config_file", "", "mysql config file")
	flag.StringVar(&redisConfigFile, "redis_config_file", "", "mysql config file")
	flag.Parse()
}

func main() {
	if mysqlConfigFile == "" || redisConfigFile == "" {
		panic("config file cannot be empty")
	}

	// step1: 初始化 idb 组件，在使用
	if err := idb.InitDBInfoManager(); err != nil {
		panic(err)
	}

	// step2: 初始化项目的原始数据库
	if err := initOriginalDBConfig(); err != nil {
		panic(err)
	}

	// steo3: 开始使用数据库
	useDBByDBConfig()
	useDBByHostAndPort("0.0.0.0", "3306", "sentry", idb.DB_TYPE_MYSQL)
	useDBByHostAndPort("10.66.121.171", "6379", "", idb.DB_TYPE_REDIS)
}

func initOriginalDBConfig() error {
	// mysql
	{
		configBytes, err := ioutil.ReadFile(mysqlConfigFile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(configBytes, &mysqlConfig); err != nil {
			return err
		}
		if err := mysqlConfig.Init(); err != nil {
			return err
		}
	}

	// redis
	{
		configBytes, err := ioutil.ReadFile(redisConfigFile)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(configBytes, &redisConfig); err != nil {
			return err
		}
		if err := redisConfig.Init(); err != nil {
			return err
		}
	}
	return nil
}

// 通过解析到的 MysqlConfig/RedisConfig，获取数据库连接
// 提倡使用此法方法，原因是此种方法可以直接通过 xxConfig 获取到数据库连接实例，然后直接进行数据库相关的操作。不用自己写这些重复代码
func useDBByDBConfig() {
	// to do something for mysql
	mysql, mysqlErr := mysqlConfig.OpenDB()
	if mysqlErr != nil {
		log.Fatal(mysqlErr)
	}
	defer mysql.CloseDB()
	if pErr := mysql.CheckDB(); pErr != nil {
		log.Fatal(pErr)
	} else {
		log.Println("mysql ping sucessfully!")
	}

	// to do something for redis pool
	redis, err := redisConfig.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	defer redis.Close()
	if _, pErr := redis.Do(context.Background(), "PING").Result(); pErr != nil {
		log.Fatal(pErr)
	} else {
		log.Println("redis pool ping sucessfully!")
	}
}

// 指定原始数据库的 host、port，获取数据库信息
// 应该避免使用此方法，原因是此种方法，只是通过 idb 获取到数据库信息，具体的创建数据库连接操作，需要自己在项目内完成
func useDBByHostAndPort(oriHost, oriPort, dbName, dbType string) {
	log.Printf("====== use %s ======\n", dbType)
	for i := 0; i < 3; i++ {
		oneDBInfo := idb.GetOneRandomDBInfo(oriHost, oriPort, dbName)
		log.Printf("get dbinfo:%v\n", oneDBInfo)
		// 使用 oneDBInfo，自己创建数据库连接
	}
}

func printRandomDBInfo(oriHost, oriPort string) {
}
