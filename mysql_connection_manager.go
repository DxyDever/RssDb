package idb

import (
	"errors"
	stdlog "log"
	"sync"
)

var (
	mysqlConnMap    sync.Map
	mysqlConnMapMtx sync.Mutex
)

// 可重复获取 mysql 连接，支持连接复用功能
// 获取之后，不要 Close。下次获取会直接返回之前创建的连接
// 如果连接使用过程中出错，可以调用：DeleteMysqlConnection 进行销毁
func GetMysqlConnection(config MysqlConfig) (*MysqlDB, error) {
	if conn, isExist := mysqlConnMap.Load(config); isExist {
		return conn.(*MysqlDB), nil
	}

	mysqlConnMapMtx.Lock()
	defer mysqlConnMapMtx.Unlock()
	if conn, isExist := mysqlConnMap.Load(config); isExist {
		return conn.(*MysqlDB), nil
	}
	conn, err := config.OpenDB()
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("mysql connection is nil")
	}
	mysqlConnMap.Store(config, conn)
	return conn, nil
}

func ResetErrorMysqlConnection(config MysqlConfig) {
	if _, isExist := mysqlConnMap.Load(config); !isExist {
		return
	}
	mysqlConnMapMtx.Lock()
	defer mysqlConnMapMtx.Unlock()
	if conn, isExist := mysqlConnMap.Load(config); isExist {
		if err := conn.(*MysqlDB).Ping(); err == nil {
			return
		}

		mysqlConnMap.Delete(config)
		conn.(*MysqlDB).CloseDB()

		conn, err := config.OpenDB()
		if err != nil {
			stdlog.Println("[idb] mysql_connection_manager ResetErrorMysqlConnection OpenDB error:", err)
			return
		}
		if conn == nil {
			return
		}
		mysqlConnMap.Store(config, conn)
	}
}

func DeleteMysqlConnection(config MysqlConfig) {
	if _, isExist := mysqlConnMap.Load(config); !isExist {
		return
	}
	mysqlConnMapMtx.Lock()
	defer mysqlConnMapMtx.Unlock()
	if conn, isExist := mysqlConnMap.Load(config); isExist {
		mysqlConnMap.Delete(config)
		conn.(*MysqlDB).CloseDB()
	}
}
