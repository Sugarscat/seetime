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
		fileAdmin, _ := os.OpenFile("./data/users/admin.json", os.O_WRONLY|os.O_CREATE, 0644)
		defer func(fileAdmin *os.File) {
			fileAdmin.Close()
		}(fileAdmin)

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
	defer func(fileUsers *os.File) {
		fileUsers.Close()
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

	jsonDataU, _ := json.Marshal(usersJsonFile)

	fileUsers.Truncate(0)
	_, err := io.WriteString(fileUsers, string(jsonDataU))
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}

	return true
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

// LoadUsers 添加信息，解析json
func LoadUsers(adminInfo []byte, userInfo []byte) {

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

	for _, userOne := range userData.Users {
		user := User{
			Id:          userOne.Id,
			Name:        userOne.Name,
			Password:    userOne.Password,
			Token:       "null",
			Identity:    userOne.Identity,
			Permissions: userOne.Permissions,
		}
		Users = append(Users, user)
	}
}

func init() {

}
