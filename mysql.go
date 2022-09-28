package idb

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/DxyDever/RssDb/idbinfo"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlDB struct {
	*sql.DB
}

func (this *MysqlDB) CheckDB() error {
	if this == nil {
		return errors.New("db instance is nil")
	}
	return this.Ping()
}

func (this *MysqlDB) CloseDB() error {
	if this == nil {
		return errors.New("db instance is nil")
	}
	return this.Close()
}

// MysqlConfig /*
/*
配置示例：
{
    "Host": "10.141.0.234",
    "Port": 3306,
    "User": "root",
    "Password": "shumeiShumei2016",
    "DbName": "sentry",
    "ConnTimeout": 1000,
    "MaxIdleConn": 50,
    "MaxOpenConn": 200,
    "ReadTimeout": 1000,
    "WriteTimeout": 1000
}
*/
type MysqlConfig struct {
	DBConfig
}

func (this *MysqlConfig) Init() error {
	if len(this.DBType) == 0 {
		this.DBType = idbinfo.DB_TYPE_MYSQL
	}
	return this.DBConfig.Init()
}

func (this *MysqlConfig) OpenDB() (*MysqlDB, error) {
	dbInfo := idbinfo.GetOneRandomDBInfo(this.Host, strconv.Itoa(this.Port), this.DbName)
	if dbInfo == nil {
		return nil, errors.New("not get dbInfo")
	}
	dbName := this.DbName
	if dbInfo.DBName != "" {
		dbName = dbInfo.DBName
	}

	var err error
	connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=%dms&readTimeout=%dms&writeTimeout=%dms&charset=utf8mb4,utf8", dbInfo.UserName, dbInfo.Password, dbInfo.Host, dbInfo.Port, dbName, this.ConnTimeout, this.ReadTimeout, this.WriteTimeout)
	if dbInfo.FlagSsl {
		connectStr = connectStr + "&tls=skip-verify"
	}
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		return nil, fmt.Errorf("dbinfo:%s, err:%v", dbInfo.JoinHostPort(), err)
	}
	if db == nil {
		return nil, fmt.Errorf("dbinfo:%s, err:db is nil", dbInfo.JoinHostPort())
	}
	db.SetMaxOpenConns(this.MaxOpenConn)
	db.SetMaxIdleConns(this.MaxIdleConn)
	if idbinfo.GetDBInfoManager().SwitchDBInfoUpdate {
		if err = db.Ping(); err != nil {
			db.Close()
			return nil, fmt.Errorf("dbinfo:%s, err:%v", dbInfo.JoinHostPort(), err)
		}
	}
	mysql := new(MysqlDB)
	mysql.DB = db
	return mysql, nil
}
