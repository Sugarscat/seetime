package account

import (
	"encoding/json"
	"net/http"
	"time"
)

var refused = Response{
	Code:    429,
	Success: false,
	Message: "请求次数过多",
	Id:      -1,
	Token:   "error",
	Time:    "",
	IP:      "",
}

var bucket = NewLeakyBucket(6, 0.1) // 桶

func openAPiLogin() {
	http.HandleFunc("/api/login", HandleLogin)
	http.HandleFunc("/api/me", HandleMe)
	err := http.ListenAndServe(":6060", nil)
	if err != nil {
		return
	}
}

// checkInfo 检测信息
func checkInfo(name string, password string, r *http.Request) Response {
	var response Response
	var allowPass bool

	for _, user := range Users {
		if user.Name == name && user.Password == password {
			allowPass = true
			token := GenerateToken(user.Id, user.Name)
			Users[user.Id].Token = token
			Users[user.Id].LastIp = Users[user.Id].ClientIp
			Users[user.Id].ClientIp = GetIP(r)
			Users[user.Id].LastTime = Users[user.Id].LoginTime
			Users[user.Id].LoginTime = time.Now().Unix()
			response = AddResponse(200, true, "验证成功", user.Id, token, GetTime(user.LastTime), user.LastIp)
			return response
		}
		allowPass = false
	}

	if !allowPass {
		if !bucket.AddWater(1) {
			return refused
		}
	}

	response = AddResponse(404, false, "用户名或密码错误", -1, "error", "", "")
	return response
}

// HandleLogin 登录
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	password := r.FormValue("password")
	response := checkInfo(name, password, r)
	jsonBytes, err := json.Marshal(response)

	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		return
	}
}

func init() {

}
