package account

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sugarscat/seetime/server/module"
)

type LoginResponse struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Id      int    `json:"id"`
	Token   string `json:"token"`
}

var refused = LoginResponse{
	Code:    429,
	Success: false,
	Id:      -1,
	Token:   "",
}

var bucket = module.NewLeakyBucket(6, 0.1) // 桶

func AddLoginResponse(code int, success bool, id int, token string) LoginResponse {
	return LoginResponse{
		Code:    code,
		Success: success,
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
			response = AddLoginResponse(200, true, user.Id, token)
			return response
		}
		allowPass = false
	}

	if !allowPass {
		if !bucket.AddWater(1) {
			return refused
		}
	}

	response = AddLoginResponse(404, false, -1, "")
	return response
}

// HandleLogin 登录
func HandleLogin(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("pwd")
	response := checkInfo(name, password, ctx)
	ctx.JSON(200, response)
}

func init() {}
