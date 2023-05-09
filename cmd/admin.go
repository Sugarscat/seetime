package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	err            error
	adminInfo      []byte
	adminFilePlace = "./data/Users/"
)

type adminFile struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Identity bool   `json:"identity"`
}

func getPwd() {
	var admin adminFile
	err = json.Unmarshal(adminInfo, &admin)
	if err != nil {
		fmt.Println(err) // ---日志
	}
	fmt.Println("Name:", admin.Name)
	fmt.Println("Password:", admin.Password)
}

func createAdminFile() {
	fileData := adminFile{
		Id:       0,
		Name:     "admin",
		Password: "QWQTime",
		Identity: true,
	}

	os.MkdirAll(adminFilePlace, 0755)
	file, err := os.Create(adminFilePlace + "admin.json")
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

var PasswordCmd = &cobra.Command{
	Use:     "admin",
	Aliases: []string{"password"},
	Short:   "Show admin user's info",
	Run: func(cmd *cobra.Command, args []string) {
		getPwd()
	},
}

func init() {
	RootCmd.AddCommand(PasswordCmd)
	adminInfo, err = os.ReadFile(adminFilePlace)
	if err != nil {
		createAdminFile()
	}

	defer func() {
		adminInfo, _ = os.ReadFile(adminFilePlace)
	}()
}
