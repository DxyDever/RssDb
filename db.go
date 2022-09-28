package idb

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/DxyDever/RssDb/idbinfo"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
	// time.Millisecond to dial tcp timeout
	ConnTimeout int
	// time.Millisecond to read tcp timeout
	ReadTimeout int
	// time.Millisecond to write tcp timeout
	WriteTimeout int

	MaxOpenConn  int
	MaxIdleConn  int
	RetryCount   int
	TimeInterval int
	// open ssl or not, true is open, false is close
	FlagSsl bool

	// support mysql/redis/redis-cluster/redis-sentinel
	// TODO add doc
	DBType string `json:"db_type"`

	// The sentinel master name.
	// Only failover clients.
	MasterName string
}

func (this *DBConfig) Init() error {
	oriDBInfo := new(idbinfo.DBInfo)
	oriDBInfo.Host = this.Host
	oriDBInfo.Port = strconv.Itoa(this.Port)
	oriDBInfo.UserName = this.User
	oriDBInfo.Password = this.Password
	oriDBInfo.DBName = this.DbName
	oriDBInfo.DBType = this.DBType
	oriDBInfo.FlagSsl = this.FlagSsl
	if err := idbinfo.AddOriginalDBInfo(oriDBInfo); err != nil {
		return fmt.Errorf("err:%v, dbinfo:%s", err, oriDBInfo.Host)
	}
	return nil
}

func (this *DBConfig) JoinHostPort() string {
	return net.JoinHostPort(this.Host, strconv.Itoa(this.Port))
}

func (this DBConfig) String() string {
	bts, _ := json.Marshal(this)
	return string(bts)
}

func (this *DBConfig) FormatErr(err interface{}) string {
	return fmt.Sprintf("db error, dbinfo:%s, error:%v", this.Host, err)
}
