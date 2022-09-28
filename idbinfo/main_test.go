package idbinfo

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	managerInitFlow()
	exitCode := m.Run()
	DestroyDBInfoManager()
	os.Exit(exitCode)
}

func managerInitFlow() {
	// step 1：初始化 dbinfo
	option := DBInfoManagerOption{
		EncryptedPwdFlag:      false,
		YinlianDynamicPwdFlag: false,
	}
	if err := InitDBInfoManager(&option); err != nil {
		panic(err)
	}

	// step 2：添加元数据库信息
	oriDBInfoStr := `
	{
		"host":"0.0.0.0",
		"port":"3306",
		"user_name":"test",
		"password":"test",
		"db_type":"mysql"
	}
	`
	if err := AddOriginalDBInfoByJson(oriDBInfoStr); err != nil {
		panic(err)
	}

	// step 3：添加备用数据库信息，注意 备用数据库信息以元数据库的 host:port 为 key
	backupDBInfoStr := `
	{
		"0.0.0.0:3306/sentry":[
			{
				"host":"10.141.0.234",
				"port":"3306",
				"user_name":"root",
				"password":"shumeiShumei2016",
				"db_type":"mysql"
			}
		],
		"0.0.0.0:3306":[
			{
				"host":"10.141.0.234",
				"port":"3306",
				"user_name":"root",
				"password":"shumeiShumei2016",
				"db_type":"mysql"
			}
		]
	}
	`
	if err := AddBackupDBInfoByJson(backupDBInfoStr); err != nil {
		panic(err)
	}

	// step 4：根据元数据库的 host、port 获取可用的数据库信息
	oriHost := "0.0.0.0"
	oriPort := "3306"
	dbInfo := GetOneRandomDBInfo(oriHost, oriPort, "sentry")
	log.Printf("get db info:%v\n", dbInfo)

	// step 5：根据数据库信息做业务操作，eg.连接数据库
}

func initDBInfo(oriDBInfo, backupDBInfo string) error {
	if err := AddOriginalDBInfoByJson(oriDBInfo); err != nil {
		return err
	}

	if err := AddBackupDBInfoByJson(backupDBInfo); err != nil {
		return err
	}
	return nil
}

func doTestManagerKeepalive(t *testing.T, oriHost, oriPort, backupDBInfoStr string) {
	{
		t.Log("====== test dbinfo ======")
		dbInfo := GetOneRandomDBInfo(oriHost, oriPort, "")
		if dbInfo == nil {
			t.Fatal("empty db info!", oriHost, oriPort)
		}
		t.Log("get db info:", dbInfo)
	}

	{
		t.Log("====== test invalid dbinfo ======")
		GetDBInfoManager().KeepaliveDBList()
		t.Log("KeepaliveDBList")
		dbInfo := GetOneRandomDBInfo(oriHost, oriPort, "")
		if dbInfo != nil {
			t.Fatal("not empty db info!")
		}
		t.Log("get db info:", dbInfo)
	}

	{
		t.Log("====== test valid dbinfo ======")
		t.Log("AddBackupDBInfo")
		if err := AddBackupDBInfoByJson(backupDBInfoStr); err != nil {
			t.Fatal(err)
		}
		GetDBInfoManager().KeepaliveDBList()
		t.Log("KeepaliveDBList")

		dbInfo := GetOneRandomDBInfo(oriHost, oriPort, "")
		if dbInfo == nil {
			t.Fatal("empty db info!")
		}
		t.Log("get db info:", dbInfo)
	}
}
