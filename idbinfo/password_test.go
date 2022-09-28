package idbinfo

import (
	"testing"
)

func TestEncryptedPwd(t *testing.T) {
	manager := GetDBInfoManager()
	oldFlag := manager.EncryptedPwdFlag
	manager.EncryptedPwdFlag = true
	defer func() {
		manager.EncryptedPwdFlag = oldFlag
	}()
	oriHost := "127.0.0.1"
	oriPort := "6379"
	backupDBInfoStr := `
	{
		"127.0.0.1:6379":[
			{
				"host":"10.66.121.171",
				"port":"6379",
				"password":"gugOYNXFwl2N5eZSotHjNZUg87bF6nTj9w9aE5XZBjk=",
				"db_type":"redis"
			}
		]
	}
	`
	t.Log("AddBackupDBInfo")
	if err := AddBackupDBInfoByJson(backupDBInfoStr); err != nil {
		t.Fatal(err)
	}
	GetDBInfoManager().KeepaliveDBList()
	t.Log("KeepaliveDBList")

	dbInfo := GetOneRandomDBInfo(oriHost, oriPort, "")
	if dbInfo == nil {
		t.Fatal("not empty db info!")
	}
	t.Log("get db info:", dbInfo)
}
