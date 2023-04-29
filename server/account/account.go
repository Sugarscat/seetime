package account

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
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

type Permissions struct {
	Situation    bool `json:"situation"`    // 系统信息
	AddTask      bool `json:"addtask"`      // 添加任务
	ChangeTask   bool `json:"changetask"`   // 修改任务
	DeleteTask   bool `json:"deletetask"`   // 删除任务
	DownloadTask bool `json:"downloadtask"` // 下载任务
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

// LeakyBucket 漏桶算法的桶定义
type LeakyBucket struct {
	capacity     float64   // 桶容量
	rate         float64   // 漏水速率
	water        float64   // 当前水量
	lastLeakTime time.Time // 上次漏水时间
}

// SaveInfo 保存用户信息
func SaveInfo(id int) bool {

	if id == 0 {
		fileAdmin, err := os.OpenFile("./data/users/admin.json", os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fmt.Println(err) // ---日志
			return false
		}
		defer func(fileAdmin *os.File) {
			err := fileAdmin.Close()
			if err != nil {
				fmt.Println(err) // ---日志
			}
		}(fileAdmin)

		jsonDataA, err := json.Marshal(
			adminJson{
				Id:       0,
				Name:     Users[0].Name,
				Password: Users[0].Password,
				Identity: true,
			})
		if err != nil {
			fmt.Println(err) // ---日志
			return false
		}

		fileAdmin.Truncate(0)
		_, err = io.WriteString(fileAdmin, string(jsonDataA))
		if err != nil {
			fmt.Println(err) // ---日志
			return false
		}

		return true
	}

	fileUsers, err := os.OpenFile("./data/users/users.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}
	defer func(fileUsers *os.File) {
		err := fileUsers.Close()
		if err != nil {
			fmt.Println(err) // ---日志
		}
	}(fileUsers)

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

	jsonDataU, err := json.Marshal(usersJsonFile)
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}

	fileUsers.Truncate(0)
	_, err = io.WriteString(fileUsers, string(jsonDataU))
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}

	return true
}

// GetTime 转换时间
func GetTime(timestamp int64) string {
	if timestamp == 0 {
		return "error"
	}
	datetime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	return datetime
}

/* NewLeakyBucket 漏桶算法 限制请求次数 */
func NewLeakyBucket(capacity, rate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity:     capacity,
		rate:         rate,
		water:        0,
		lastLeakTime: time.Now(),
	}
}

func (b *LeakyBucket) AddWater(amount float64) bool {
	// 先漏水
	b.Leak()
	// 再加水
	if b.water+amount <= b.capacity {
		b.water += amount
		return true // 添加成功
	} else {
		return false // 添加失败
	}
}

func (b *LeakyBucket) Leak() {
	now := time.Now()
	elapsed := now.Sub(b.lastLeakTime).Seconds() // 计算距离上次漏水时间
	b.water = Max(b.water-elapsed*b.rate, 0)     // 漏掉一定量的水
	b.lastLeakTime = now                         // 更新上次漏水时间
}

func Max(x, y float64) float64 {
	if x > y {
		return x
	}
	return y
}

/* 漏桶算法 END */

// AddUser 添加信息，解析json
func AddUser(adminInfo []byte, userInfo []byte) {

	AdminInfo = adminInfo
	UsersInfo = userInfo

	addAdmin()
	addUser()
}

func addAdmin() {
	err := json.Unmarshal(AdminInfo, &adminData)
	if err != nil {
		fmt.Println(err) // ---日志
	}

	admin := User{
		Id:          adminData.Id,
		Name:        adminData.Name,
		Password:    adminData.Password,
		Token:       "null",
		Identity:    adminData.Identity,
		Permissions: 11111,
	}

	Users = append(Users, admin)
}

func addUser() {
	err := json.Unmarshal(UsersInfo, &userData)
	if err != nil {
		fmt.Println(err) // ---日志
	}

	for i := 0; i < len(userData.Users); i++ {
		user := User{
			Id:          userData.Users[i].Id,
			Name:        userData.Users[i].Name,
			Password:    userData.Users[i].Password,
			Token:       "null",
			Identity:    userData.Users[i].Identity,
			Permissions: userData.Users[i].Permissions,
		}
		Users = append(Users, user)
	}
}

func init() {

}
