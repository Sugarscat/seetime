package account

import (
	"encoding/json"
	"net/http"
	"time"
)

type LoginResponse struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Message string `json:"message"` // 消息
	Id      int    `json:"id"`
	Token   string `json:"token"`
}

var refused = LoginResponse{
	Code:    429,
	Success: false,
	Message: "请求次数过多",
	Id:      -1,
	Token:   "null",
}

var bucket = NewLeakyBucket(6, 0.1) // 桶

func AddLoginResponse(code int, success bool, message string, id int, token string) LoginResponse {
	return LoginResponse{
		Code:    code,
		Success: success,
		Message: message,
		Id:      id,
		Token:   token,
	}
}

// checkInfo 检测信息
func checkInfo(name string, password string, r *http.Request) LoginResponse {
	var response LoginResponse
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
			response = AddLoginResponse(200, true, "验证成功", user.Id, token)
			return response
		}
		allowPass = false
	}

	if !allowPass {
		if !bucket.AddWater(1) {
			return refused
		}
	}

	response = AddLoginResponse(404, false, "用户名或密码错误", -1, "null")
	return response
}

// HandleLogin 登录
func HandleLogin(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	password := r.FormValue("password")
	response := checkInfo(name, password, r)

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		// ---日志
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		// ---日志
		return
	}
}

func init() {

}
