package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	data []byte
	err  error
)

type Admin struct {
	Id       int
	Name     string
	Password string
}

func getPwd() {
	var admin Admin
	err = json.Unmarshal(data, &admin)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(admin.Password)
}

func init() {
	data, err = ioutil.ReadFile("./data/users/admin.json")
	if err != nil {
		fmt.Println(err)
	}
}
