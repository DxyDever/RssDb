package idb

import (
	"flag"

	"github.com/DxyDever/RssDb/idbinfo"
)

var (
	iDBInfoConfigFile string
)

const (
	DB_TYPE_MYSQL = "mysql"
	DB_TYPE_REDIS = "redis"
)

func init() {
	flag.StringVar(&iDBInfoConfigFile, "idbinfo_config_file", "", "idbinfo config file")
}

func InitDBInfoManager() error {
	return idbinfo.InitDBInfoManagerByConfigFile(iDBInfoConfigFile)
}

func InitDBInfoManager_default() error {
	option := idbinfo.DBInfoManagerOption{
		EncryptedPwdFlag:      false,
		YinlianDynamicPwdFlag: false,
	}
	if err := idbinfo.InitDBInfoManager(&option); err != nil {
		return err
	}
	return nil
}

func GetOneRandomDBInfo(host, port, dbName string) *idbinfo.DBInfo {
	return idbinfo.GetOneRandomDBInfo(host, port, dbName)
}
