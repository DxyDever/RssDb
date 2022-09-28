package iredis

import (
	"testing"
)

func TestRedisSentinel(t *testing.T) {
	sentinelInfo := &RedisSentinelInfo{
		Host: "10.141.3.217",
		Port: "26379",
		Name: "M1",
	}
	masterHost, masterPort, sentinelErr := sentinelInfo.GetMasterInfo()
	if sentinelErr != nil {
		t.Fatal(sentinelErr)
	} else {
		t.Logf("get master. host:%s, port:%s\n", masterHost, masterPort)
	}
}
