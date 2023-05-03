package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sugarscat/seetime/server/account"
	"github.com/sugarscat/seetime/server/tasks"
)

type adminFile struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Identity bool   `json:"identity"`
}

var (
	adminFilePlace = "./data/Users/admin.json"
	usersFilePlace = "./data/Users/Users.json"
	tasksFilePlace = "./data/tasks/tasks.json"
	adminInfo      []byte // 读取json文件
	usersInfo      []byte
	tasksInfo      []byte
	err            error
)

func CreateAdminFile() {
	fileData := adminFile{
		Id:       0,
		Name:     "admin",
		Password: "QWQTime",
		Identity: true,
	}
	file, err := os.Create(adminFilePlace)
	if err != nil {
		fmt.Println(err) // ---日志
	}

	fileJson, _ := json.Marshal(fileData)
	file.Truncate(0)
	_, err = file.WriteString(string(fileJson))
	if err != nil {
		fmt.Println(err) // ---日志
	}

	defer file.Close()
}

func SendInfo() {
	account.LoadUsers(adminInfo, usersInfo)
	tasks.LoadTasks(tasksInfo)
}

func Loading() {
	SendInfo()
	defer OpenRouter()

}

func init() {
	adminInfo, err = os.ReadFile(adminFilePlace)
	if err != nil {
		CreateAdminFile()
	}

	defer func() {
		adminInfo, _ = os.ReadFile(adminFilePlace)
		usersInfo, _ = os.ReadFile(usersFilePlace)
		tasksInfo, _ = os.ReadFile(tasksFilePlace)
	}()
}
