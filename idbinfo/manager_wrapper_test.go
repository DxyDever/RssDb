package idbinfo

import (
	"testing"
)

func TestManagerKeepaliveForRedis(t *testing.T) {
	originalDBInfoStr := `
	{
		"host":"127.0.0.1",
		"port":"6379",
		"password":"test",
		"db_type":"redis"
	}
	`
	invalidBackupDBInfoStr := `
	{
		"127.0.0.1:6379":[
			{
				"host":"127.0.0.2",
				"port":"6379",
				"password":"test",
				"db_type":"redis"
			},
			{
				"host":"127.0.0.3",
				"port":"6379",
				"password":"test",
				"db_type":"redis"
			}
		]
	}
	`
	if err := initDBInfo(originalDBInfoStr, invalidBackupDBInfoStr); err != nil {
		t.Fatal(err)
	}
	oriHost := "127.0.0.1"
	oriPort := "6379"
	validBackupDBInfoStr := `
	{
		"127.0.0.1:6379":[
			{
				"host":"10.66.121.171",
				"port":"6379",
				"password":"crs-huhkvgpi:shumei123",
				"db_type":"redis"
			}
		]
	}
	`
	doTestManagerKeepalive(t, oriHost, oriPort, validBackupDBInfoStr)
}

func TestManagerKeepaliveForMysql(t *testing.T) {
	originalDBInfoStr := `
	{
		"host":"127.0.0.1",
		"port":"3306",
		"user_name":"test",
		"password":"test",
		"db_type":"mysql"
	}
	`
	invalidBackupDBInfoStr := `
	{
		"127.0.0.1:3306":[
			{
				"host":"127.0.0.2",
				"port":"3306",
				"user_name":"test",
				"password":"test",
				"db_type":"mysql"
			},
			{
				"host":"127.0.0.3",
				"port":"3306",
				"user_name":"test",
				"password":"test",
				"db_type":"mysql"
			}
		]
	}
	`
	if err := initDBInfo(originalDBInfoStr, invalidBackupDBInfoStr); err != nil {
		t.Fatal(err)
	}
	oriHost := "127.0.0.1"
	oriPort := "3306"
	validBackupDBInfoStr := `
	{
		"127.0.0.1:3306":[
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
	doTestManagerKeepalive(t, oriHost, oriPort, validBackupDBInfoStr)
}
