package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	err            error
	adminInfo      []byte
	adminFilePlace = "./data/Users/admin.json"
)

type adminFile struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Identity bool   `json:"identity"`
}

func GetPwd() {
	var admin adminFile
	err = json.Unmarshal(adminInfo, &admin)
	if err != nil {
		fmt.Println(err) // ---日志
	}
	fmt.Println(admin.Name)
	fmt.Println(admin.Password)
}

func createAdminFile() {
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

func init() {
	adminInfo, err = os.ReadFile(adminFilePlace)
	if err != nil {
		createAdminFile()
	}

	defer func() {
		adminInfo, _ = os.ReadFile(adminFilePlace)
	}()
}
