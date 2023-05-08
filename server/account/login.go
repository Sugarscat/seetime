package account

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sugarscat/seetime/server/module"
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
	Token:   "",
}

var bucket = module.NewLeakyBucket(6, 0.1) // 桶

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
func checkInfo(name string, password string, ctx *gin.Context) LoginResponse {
	var response LoginResponse
	var allowPass bool

	for _, user := range Users {
		if user.Name == name && user.Password == password {
			allowPass = true
			token := GenerateToken(user.Id, user.Name)
			Users[user.Id].Token = token
			Users[user.Id].LastIp = Users[user.Id].ClientIp
			Users[user.Id].ClientIp = ctx.ClientIP()
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

	response = AddLoginResponse(404, false, "用户名或密码错误", -1, "")
	return response
}

// HandleLogin 登录
func HandleLogin(ctx *gin.Context) {

	name := ctx.Query("name")
	password := ctx.Query("pwd")
	response := checkInfo(name, password, ctx)

	ctx.JSON(200, response)
}

func init() {}
