package idbinfo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/DxyDever/RssDb/ipassword"
	"github.com/DxyDever/RssDb/iredis"
)

const (
	DB_UPDATE_METHOD_PING = "ping"
	DB_UPDATE_METHOD_ISM  = "ism"
)

type DBInfoManagerOption struct {
	EncryptedPwdFlag        bool   `json:"encrypted_pwd_flag"`
	YinlianDynamicPwdFlag   bool   `json:"yinlian_dynamic_pwd_flag"`
	SwitchDBInfoUpdate      bool   `json:"switch_db_info_update"`
	DBInfoUpdateIntervalSec int    `json:"db_info_update_interval_sec"`
	DBInfoUpdateMethod      string `json:"db_info_update_method"`
}

func (this *DBInfoManagerOption) init() {
	if this.DBInfoUpdateIntervalSec == 0 {
		this.DBInfoUpdateIntervalSec = 2
	}
	if len(this.DBInfoUpdateMethod) == 0 {
		this.DBInfoUpdateMethod = DB_UPDATE_METHOD_PING
	}
}

type RedisSentinelInfo struct {
	iredis.RedisSentinelInfo
}

type dbInfoManager struct {
	DBInfoManagerOption
	goodHostSet      HostInfoSet
	badHostSet       HostInfoSet
	pwdServer        *ipassword.PasswordServer
	originalDBMap    map[string]*DBInfo
	redisSentinelMap map[string][]*RedisSentinelInfo
	mutex            sync.RWMutex
	ctx              context.Context
	ctxCancelFunc    context.CancelFunc
}

func (this *dbInfoManager) Init(option *DBInfoManagerOption) {
	this.DBInfoManagerOption = *option
	this.DBInfoManagerOption.init()
	this.originalDBMap = make(map[string]*DBInfo, 10)
	this.goodHostSet.init()
	this.badHostSet.init()
	this.ctx, this.ctxCancelFunc = context.WithCancel(context.Background())
	go this.CycleKeepaliveDBList()
}

func (this *dbInfoManager) GetOneRandomDBInfo(host, port, dbName string) *DBInfo {
	info := this.goodHostSet.getOneRandomDBInfo(host, port, dbName)
	if info == nil {
		key := generateDBKey(host, port, dbName)
		return this.originalDBMap[key]
	}
	return info
}

func (this *dbInfoManager) AddOriginalDBInfo(db *DBInfo) error {
	if db.Host == "" {
		return errors.New("host cannot be empty")
	}
	key := generateDBKey(db.Host, db.Port, db.DBName)
	if _, isExist := this.originalDBMap[key]; isExist {
		return nil
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()
	if _, isExist := this.originalDBMap[key]; isExist {
		return nil
	}
	this.originalDBMap[key] = db
	if err := this.addDBInfos(key, []*DBInfo{db}); err != nil {
		return err
	}
	log.Printf("idb add dbkey:%s\n", key)
	return nil
}

func (this *dbInfoManager) AddBackupDBInfo(dbs map[string][]*DBInfo) error {
	for key, dbList := range dbs {
		this.addDBInfos(key, dbList)
	}
	return nil
}

func (this *dbInfoManager) addDBInfos(key string, dbList []*DBInfo) error {
	if err := this.dealDBInfosForAdd(key, dbList); err != nil {
		return fmt.Errorf("AddBackupDBInfo dealDBInfosForAdd-err:%v", err)
	}
	this.goodHostSet.AddDBInfoList(key, dbList)
	return nil
}

func (this *dbInfoManager) AddPasswordServer(pwdServer ipassword.PasswordServer) error {
	this.pwdServer = &pwdServer
	return nil
}

func (this *dbInfoManager) AddRedisSentinel(sentinelMap map[string][]*RedisSentinelInfo) error {
	if this.redisSentinelMap == nil {
		this.redisSentinelMap = sentinelMap
	} else {
		for key, sentinelList := range sentinelMap {
			this.redisSentinelMap[key] = sentinelList
		}
	}
	return nil
}

func (this *dbInfoManager) dealDBInfosForAdd(dbKey string, dbInfos []*DBInfo) error {
	for _, dbInfo := range dbInfos {
		if err := this.dealDBInfoForYinlianDynamicPwd(dbKey, dbInfo); err != nil {
			return err
		}
		if err := this.dealDBInfoForEncryptedPwd(dbKey, dbInfo); err != nil {
			return err
		}
		if err := this.dealDBInfoForRedisSentinel(dbKey, dbInfo); err != nil {
			return err
		}
	}
	return nil
}

// redis 哨兵处理
func (this *dbInfoManager) dealDBInfoForRedisSentinel(dbKey string, dbInfo *DBInfo) error {
	sentinelList, isSentinel := this.redisSentinelMap[dbKey]
	if !isSentinel {
		return nil
	}
	var (
		masterHost  string
		masterPort  string
		sentinelErr error
	)

	badSentinelIndexMap := make(map[int]bool)
	for {
		rand.Seed(time.Now().UnixNano())
		index := rand.Intn(len(sentinelList))
		sentinelInfo := sentinelList[index]
		masterHost, masterPort, sentinelErr = sentinelInfo.GetMasterInfo()
		if sentinelErr == nil {
			break
		}
		badSentinelIndexMap[index] = true
		if len(badSentinelIndexMap) == len(sentinelList) {
			break
		}
		log.Printf("redis sentinel get master error!dbkey:%s, sentinel:%s, error:%v\n", dbKey, sentinelInfo, sentinelErr)
	}

	if sentinelErr != nil {
		log.Printf("redis sentinel donot get master info!dbkey:%s, error:%v\n", dbKey, sentinelErr)
		return sentinelErr
	}
	if masterHost == "" || masterPort == "" {
		return errors.New("redis sentinel get empty master info")
	}
	dbInfo.Host = masterHost
	dbInfo.Port = masterPort
	return nil
}

// 密码解密
func (this *dbInfoManager) dealDBInfoForEncryptedPwd(dbKey string, dbInfo *DBInfo) error {
	// 动态密码开启时，针对 mysql 不使用密码解密功能
	if !this.EncryptedPwdFlag || (this.YinlianDynamicPwdFlag && dbInfo.DBType == DB_TYPE_MYSQL) {
		return nil
	}

	// 空密码，不作处理
	if dbInfo.Password == "" {
		return nil
	}

	password, err := ipassword.Decrypt(dbInfo.Password, ipassword.ENCRYPT_ALGO_AES)
	if err != nil {
		return err
	}
	dbInfo.Password = password
	return nil
}

// 获取动态密码
func (this *dbInfoManager) dealDBInfoForYinlianDynamicPwd(dbKey string, dbInfo *DBInfo) error {
	if !this.YinlianDynamicPwdFlag || dbInfo.DBType != DB_TYPE_MYSQL {
		return nil
	}
	if this.pwdServer == nil {
		return errors.New("password server info must be added firstly!")
	}
	password, err := this.pwdServer.GetDynamicPasswd(dbInfo.DBName, dbInfo.UserName)
	if err != nil {
		return err
	}
	dbInfo.Password = password
	return nil
}

func (this *dbInfoManager) IsExistOriginalDB(key string) bool {
	_, isExist := this.originalDBMap[key]
	return isExist
}

func (this *dbInfoManager) CycleKeepaliveDBList() {
	if !this.SwitchDBInfoUpdate {
		return
	}
	t := time.NewTicker(time.Duration(time.Duration(this.DBInfoUpdateIntervalSec) * time.Second))
	defer t.Stop()
	for {
		select {
		case <-t.C:
			this.KeepaliveDBList()
		case <-this.ctx.Done():
			fmt.Fprintln(os.Stderr, "CycleKeepaliveDBList exit")
			return
		}
	}
}

func (this *dbInfoManager) KeepaliveDBList() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintln(os.Stderr, "CycleKeepaliveDBList error", err)
		}
	}()

	for dbKey, _ := range this.originalDBMap {
		newGoodDBList := make([]*DBInfo, 0, 10)
		newBadDBList := make([]*DBInfo, 0, 10)

		// check good list
		for _, db := range this.goodHostSet.GetDBInfoList(dbKey) {
			if err := this.dealDBInfoForKeepalive(dbKey, db); err != nil {
				fmt.Fprintf(os.Stderr, "deal ihostinfo cycle keepalive error.db:%v.err:%v\n", (*db).Host, err)
			}
			if this.checkDBInfoIsAlive(db) {
				newGoodDBList = append(newGoodDBList, db)
			} else {
				newBadDBList = append(newBadDBList, db)
			}
		}

		// check bad list
		for _, db := range this.badHostSet.GetDBInfoList(dbKey) {
			if err := this.dealDBInfoForKeepalive(dbKey, db); err != nil {
				fmt.Fprintf(os.Stderr, "deal ihostinfo cycle keepalive error.db:%v.err:%v\n", (*db).Host, err)
			}
			if this.checkDBInfoIsAlive(db) {
				newGoodDBList = append(newGoodDBList, db)
			} else {
				newBadDBList = append(newBadDBList, db)
			}
		}

		this.goodHostSet.ResetDBInfoList(dbKey, newGoodDBList)
		this.badHostSet.ResetDBInfoList(dbKey, newBadDBList)
	}
}

func (this *dbInfoManager) checkDBInfoIsAlive(db *DBInfo) bool {
	switch this.DBInfoUpdateMethod {
	case DB_UPDATE_METHOD_ISM:
		return db.ISmPing()
	default:
		return db.Ping()
	}
}

func (this *dbInfoManager) dealDBInfoForKeepalive(dbKey string, dbInfo *DBInfo) error {
	if err := this.dealDBInfoForYinlianDynamicPwd(dbKey, dbInfo); err != nil {
		return err
	}
	if err := this.dealDBInfoForRedisSentinel(dbKey, dbInfo); err != nil {
		return err
	}
	return nil
}

func (this *dbInfoManager) ClearBackupDB() {
	this.goodHostSet.Clear()
	this.badHostSet.Clear()
}

func (this *dbInfoManager) Destroy() {
	this.ctxCancelFunc()
}
