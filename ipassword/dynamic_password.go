package ipassword

import (
	"bytes"
	"errors"
	"net"
	"strings"
)

type PasswordServer struct {
	IP      string `json:"ip"`
	Port    string `json:"port"`
	BakIP   string `json:"bak_ip"`
	BakPort string `json:"bak_port"`
}

func (this *PasswordServer) GetDynamicPasswd(dbName string, dbUserName string) (string, error) {
	passwd, err := getPasswd(this.IP, this.Port, dbName, dbUserName)
	if err != nil {
		passwd, err = getPasswd(this.BakIP, this.BakPort, dbName, dbUserName)
		return passwd, err
	}
	if len(passwd) == 0 {
		passwd, err = getPasswd(this.BakIP, this.BakPort, dbName, dbUserName)
	}
	return passwd, err
}

func getPasswd(pwdServerIp string, pwdServerPort string, dbName string, dbUserName string) (string, error) {
	host := pwdServerIp + ":" + pwdServerPort
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return "", errors.New("connect to " + host + "failed, error: " + err.Error())
	}
	defer conn.Close()
	//send data
	var data bytes.Buffer
	data.WriteByte(byte(14))
	data.WriteByte(byte(2))
	data.WriteByte(byte(' '))
	data.Write([]byte(dbName))
	data.WriteByte(byte(' '))
	data.Write([]byte(dbUserName))
	data.WriteByte(byte(13))
	data.WriteByte(byte('\n'))
	conn.Write(data.Bytes())
	//receive response
	buffer := make([]byte, 10240)
	n, rerr := conn.Read(buffer)
	if rerr != nil {
		return "", errors.New("receive passwd failed, error: " + rerr.Error())
	}
	passwd := ""
	response := string(buffer[:n])
	if strings.Index(response, "OK") == 0 {
		passwd = strings.TrimSpace(response[4:])
	}
	return passwd, nil
}
