package server

import (
	"fmt"
	"os"
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
	adminInfo, err = os.ReadFile("./data/Users/admin.json")
	if err != nil {
		fmt.Println(err) // ---日志
	}
	usersInfo, err = os.ReadFile("./data/Users/Users.json")
	if err != nil {
		fmt.Println(err) // ---日志
	}
}
