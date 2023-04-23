package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	AdminInfo []byte
	UsersInfo []byte
	adminData adminJson // 解析json文件
	userData  usersJson
	Users     = make([]user, 0, 1)
)

type adminJson struct {
	Id       int
	Name     string
	Password string
}

type usersJson struct {
	Users []struct {
		Id       int
		Name     string
		Password string
	}
}

type permissions struct {
	situation    bool // 系统信息
	addTask      bool // 添加任务
	changeTask   bool // 修改任务
	deleteTask   bool // 删除任务
	downloadTask bool // 下载任务
}

type user struct {
	Id         int
	Name       string
	Password   string
	Token      string
	Identity   bool        // 是否是管理员
	LastIp     string      // 上次登录的IP
	ClientIp   string      //本次登录的IP
	LastTime   int64       //上次登录的时间
	LoginTime  int64       //本次登录的时间
	permission permissions //权限
}

// 漏桶算法的桶定义
type LeakyBucket struct {
	capacity     float64   // 桶容量
	rate         float64   // 漏水速率
	water        float64   // 当前水量
	lastLeakTime time.Time // 上次漏水时间
}

// GetTime 转换时间
func GetTime(timestamp int64) string {
	if timestamp == 0 {
		return ""
	}
	datetime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
	return datetime
}

// 获取请求 IP
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

// json 回复
type Response struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Message string `json:"message"` // 消息
	Id      int    `json:"id"`
	Token   string `json:"token"`
	Time    string `json:"time"` // 上次登录时间
	IP      string `json:"ip"`   // 上次登录IP
}

/* 漏桶算法 */
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

// 添加回复
func AddResponse(code int, success bool, message string, id int, token string, time string, ip string) Response {
	return Response{
		Code:    code,
		Success: success,
		Message: message,
		Id:      id,
		Token:   token,
		Time:    time,
		IP:      ip,
	}
}

// AddInfo 添加信息，解析json
func AddInfo(adminInfo []byte, userInfo []byte) {
	defer openAPiLogin()

	AdminInfo = adminInfo
	UsersInfo = userInfo

	addAdmin()
	addUser()
}

func addAdmin() {
	Err := json.Unmarshal(AdminInfo, &adminData)
	if Err != nil {
		fmt.Println(3, Err)
	}

	admin := user{
		Id:       adminData.Id,
		Name:     adminData.Name,
		Password: adminData.Password,
		Token:    "err",
		Identity: true,
		permission: permissions{
			situation:    true,
			addTask:      true,
			changeTask:   true,
			deleteTask:   true,
			downloadTask: true,
		},
	}

	Users = append(Users, admin)
}

func addUser() {
	Err := json.Unmarshal(UsersInfo, &userData)
	if Err != nil {
		fmt.Println(4, Err)
	}

	for i := 0; i < len(userData.Users); i++ {
		user := user{
			Id:       userData.Users[i].Id,
			Name:     userData.Users[i].Name,
			Password: userData.Users[i].Password,
			Token:    "err",
			Identity: false,
			permission: permissions{
				situation:    false,
				addTask:      false,
				changeTask:   false,
				deleteTask:   false,
				downloadTask: false,
			},
		}
		Users = append(Users, user)
	}
}

func init() {

}
