package account

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var (
	AdminInfo []byte
	UsersInfo []byte
	adminData adminJson // 解析json文件
	userData  usersJson
	Users     = make([]User, 0, 1)
)

type adminJson struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Identity bool   `json:"identity"`
}

type usersCom struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	Identity    bool   `json:"identity"`
	Permissions int    `json:"permissions"`
}

type usersJson struct {
	Users []usersCom `json:"users"`
}

type User struct {
	Id          int
	Name        string
	Password    string
	Token       string
	Identity    bool   // 是否是管理员
	LastIp      string // 上次登录的IP
	ClientIp    string //本次登录的IP
	LastTime    int64  //上次登录的时间
	LoginTime   int64  //本次登录的时间
	Permissions int    //权限
}

// SaveInfo 保存用户信息
func SaveInfo(id int) bool {

	if id == 0 {
		fileAdmin, _ := os.OpenFile("./data/users/admin.json", os.O_WRONLY|os.O_CREATE, 0644)
		defer fileAdmin.Close()

		jsonDataA, _ := json.Marshal(
			adminJson{
				Id:       0,
				Name:     Users[0].Name,
				Password: Users[0].Password,
				Identity: true,
			})

		fileAdmin.Truncate(0)
		_, err := io.WriteString(fileAdmin, string(jsonDataA))
		if err != nil {
			fmt.Println(err) // ---日志
			return false
		}

		return true
	}

	fileUsers, _ := os.OpenFile("./data/users/users.json", os.O_WRONLY|os.O_CREATE, 0644)
	defer fileUsers.Close()

	var usersJsonFile usersJson
	usersJsonFile.Users = make([]usersCom, 0, 1)
	for i, user := range Users {
		if i == 0 {
			continue
		}
		usersJsonFile.Users = append(usersJsonFile.Users, usersCom{
			Id:          user.Id,
			Name:        user.Name,
			Password:    user.Password,
			Identity:    user.Identity,
			Permissions: user.Permissions})
	}

	jsonDataU, _ := json.Marshal(usersJsonFile)

	fileUsers.Truncate(0)
	_, err := io.WriteString(fileUsers, string(jsonDataU))
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}

	return true
}

// ParsingPermissions 分析用户权限
func ParsingPermissions(id int, work string) bool {
	switch work {
	case "situation":
		return Users[id].Permissions/10000 >= 1
	case "addTask":
		return Users[id].Permissions%10000/1000 >= 1
	case "updateTask":
		return Users[id].Permissions%10000%1000/100 >= 1
	case "deleteTask":
		return Users[id].Permissions%10000%1000%100/10 >= 1
	case "exportTask":
		return Users[id].Permissions%10000%1000%100%10 >= 1
	default:
		return false
	}
}

// LoadUsers 添加信息，解析json
func LoadUsers(adminInfo []byte, userInfo []byte) {

	AdminInfo = adminInfo
	UsersInfo = userInfo

	addAdmin()
	addUser()
}

func addAdmin() {
	json.Unmarshal(AdminInfo, &adminData)
	admin := User{
		Id:          adminData.Id,
		Name:        adminData.Name,
		Password:    adminData.Password,
		Token:       "",
		Identity:    true,
		Permissions: 11111,
	}
	Users = append(Users, admin)
}

func addUser() {
	json.Unmarshal(UsersInfo, &userData)
	for _, userOne := range userData.Users {
		user := User{
			Id:          userOne.Id,
			Name:        userOne.Name,
			Password:    userOne.Password,
			Token:       "",
			Identity:    userOne.Identity,
			Permissions: userOne.Permissions,
		}
		Users = append(Users, user)
	}
}

func init() {}
