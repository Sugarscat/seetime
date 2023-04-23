package server

import (
	"fmt"
	"io/ioutil"
	"seetime/server/account"
)

var (
	adminInfo []byte // 读取json文件
	usersInfo []byte
	Err       error
)

func SendInfo() {
	account.AddInfo(adminInfo, usersInfo)
}

func init() {
	adminInfo, Err = ioutil.ReadFile("./data/Users/admin.json")
	if Err != nil {
		// ---日志
		fmt.Println(1, Err)
	}
	usersInfo, Err = ioutil.ReadFile("./data/Users/Users.json")
	if Err != nil {
		// ---日志
		fmt.Println(2, Err)
	}
}
