package idbinfo

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/DxyDever/RssDb/ipassword"
)

var (
	gManager *dbInfoManager
)

type DBInfoConfig struct {
	ManagerOption    *DBInfoManagerOption            `json:"manager_option"`
	BackupDbinfoMap  map[string][]*DBInfo            `json:"backup_dbinfo_map"`
	RedisSentinelMap map[string][]*RedisSentinelInfo `json:"redis_sentinel_map"`
	DynamicPwdServer ipassword.PasswordServer        `json:"dynamic_pwd_server"`
}

func InitDBInfoManagerByConfigFile(configFile string) error {
	configBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	conf := new(DBInfoConfig)
	if err := json.Unmarshal(configBytes, conf); err != nil {
		return err
	}

	if err := InitDBInfoManager(conf.ManagerOption); err != nil {
		return err
	}
	if err := AddBackupDBInfo(conf.BackupDbinfoMap); err != nil {
		return err
	}
	if err := AddRedisSentinel(conf.RedisSentinelMap); err != nil {
		return err
	}
	if err := AddPasswordServer(conf.DynamicPwdServer); err != nil {
		return err
	}
	return nil
}

func InitDBInfoManager(option *DBInfoManagerOption) error {
	if gManager != nil {
		return errors.New("initialize repeatedly!")
	}
	gManager = new(dbInfoManager)
	gManager.Init(option)
	return nil
}

func GetDBInfoManager() *dbInfoManager {
	if gManager == nil {
		panic("idbinfo must be inited firstly")
	}
	return gManager
}

func DestroyDBInfoManager() {
	if gManager == nil {
		return
	}
	gManager.Destroy()
	gManager = nil
}

// 随机获取一个可用的 db info
func GetOneRandomDBInfo(host, port, dbName string) *DBInfo {
	if gManager == nil {
		panic("idbinfo must be inited firstly")
	}
	return gManager.GetOneRandomDBInfo(host, port, dbName)
}

func ClearBackupDB() {
	if gManager == nil {
		return
	}
	gManager.ClearBackupDB()
}

// 添加 backup db 信息，map key 为原始 db 的 host:port
func AddBackupDBInfo(dbs map[string][]*DBInfo) error {
	if gManager == nil {
		panic("idbinfo must be inited firstly")
	}
	return gManager.AddBackupDBInfo(dbs)
}

// 添加 backup db 信息（通过 json string 格式）
// json string 格式如下：
// {
// 	"127.0.0.1:6379":[
// 		{
// 			"host":"127.0.0.2",
// 			"port":"6379",
// 			"password":"test",
// 			"db_type":"redis"
// 		},
// 		{
// 			"host":"127.0.0.3",
// 			"port":"6379",
// 			"password":"test",
// 			"db_type":"redis"
// 		}
// 	]
// }
func AddBackupDBInfoByJson(backupDBInfoStr string) error {
	dbListMap := make(map[string][]*DBInfo)
	if err := json.Unmarshal([]byte(backupDBInfoStr), &dbListMap); err != nil {
		return err
	}
	if err := AddBackupDBInfo(dbListMap); err != nil {
		return err
	}
	return nil
}

func AddPasswordServer(pwdServer ipassword.PasswordServer) error {
	if gManager == nil {
		panic("idbinfo must be inited firstly")
	}
	return gManager.AddPasswordServer(pwdServer)
}

// 添加原始 db 信息（通过 json string 格式）
// json string 格式如下：
// {
// 	"host":"127.0.0.1",
// 	"port":"6379",
// 	"password":"test",
// 	"db_type":"redis"
// }
func AddOriginalDBInfoByJson(dbInfoStr string) error {
	originalDB := new(DBInfo)
	if err := json.Unmarshal([]byte(dbInfoStr), originalDB); err != nil {
		return err
	}
	return AddOriginalDBInfo(originalDB)
}

// 添加原始 db 信息
func AddOriginalDBInfo(db *DBInfo) error {
	if gManager == nil {
		panic("idbinfo must be inited firstly")
	}
	return gManager.AddOriginalDBInfo(db)
}

func AddRedisSentinelByJson(sentinelInfoStr string) error {
	sentinelMap := make(map[string][]*RedisSentinelInfo)
	if err := json.Unmarshal([]byte(sentinelInfoStr), &sentinelMap); err != nil {
		return err
	}
	return AddRedisSentinel(sentinelMap)
}

func AddRedisSentinel(sentinelMap map[string][]*RedisSentinelInfo) error {
	if gManager == nil {
		panic("idbinfo must be inited firstly")
	}
	return gManager.AddRedisSentinel(sentinelMap)
}
