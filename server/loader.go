package server

import (
	"fmt"
	"io/ioutil"
	"seetime/server/account"
)

var (
	adminInfo []byte // 读取json文件
	usersInfo []byte
	err       error
)

func SendInfo() {
	account.AddUser(adminInfo, usersInfo)
}

func init() {
	adminInfo, err = ioutil.ReadFile("./data/Users/admin.json")
	if err != nil {
		fmt.Println(err) // ---日志
	}
	usersInfo, err = ioutil.ReadFile("./data/Users/Users.json")
	if err != nil {
		fmt.Println(err) // ---日志
	}
}
