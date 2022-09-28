package iredis

import (
	"encoding/json"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisSentinelInfo struct {
	Host string
	Port string
	Name string
}

func (this *RedisSentinelInfo) String() string {
	bts, _ := json.Marshal(this)
	return string(bts)
}

func (this *RedisSentinelInfo) GetMasterInfo() (host, port string, err error) {
	c, err := redis.DialTimeout("tcp", net.JoinHostPort(this.Host, this.Port),
		time.Duration(20)*time.Millisecond,
		time.Duration(20)*time.Millisecond,
		time.Duration(20)*time.Millisecond)
	if err != nil {
		return
	}
	defer c.Close()
	reply, err := c.Do("SENTINEL", "MASTER", this.Name)
	if err != nil {
		return
	}
	host = string((reply.([]interface{}))[3].([]byte))
	port = string((reply.([]interface{}))[5].([]byte))
	return
}
