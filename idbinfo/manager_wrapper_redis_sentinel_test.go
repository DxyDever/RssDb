package idbinfo

import (
	"testing"
)

// go test -v -test.run TestManagerRedisInfo code.aliyun.com/module-go/idb/idbinfo
func TestManagerRedisInfo(t *testing.T) {
	oriDBInfoStr := `
	{
		"host":"10.141.3.217",
		"port":"26379",
		"db_type":"redis"
	}
	`
	if err := AddOriginalDBInfoByJson(oriDBInfoStr); err != nil {
		t.Fatal(err)
	}
	sentinelInfoStr := `
	{
		"10.141.3.217:26379":[
			{
				"host":"10.141.3.217",
				"port":"26379",
				"name": "M1"
			},
			{
				"host":"10.141.3.217",
				"port":"26380",
				"name": "M1"
			},
			{
				"host":"10.141.3.217",
				"port":"26381",
				"name": "M1"
			}
		]
	}
	`
	oriHost := "10.141.3.217"
	oriPort := "26379"
	t.Log("AddRedisSentinel")
	if err := AddRedisSentinelByJson(sentinelInfoStr); err != nil {
		t.Fatal(err)
	}
	t.Log("KeepaliveDBList")
	GetDBInfoManager().KeepaliveDBList()

	dbInfo := GetOneRandomDBInfo(oriHost, oriPort, "")
	if dbInfo == nil {
		t.Fatal("empty db info!")
	}
	t.Log("get db info:", dbInfo)
}
