package account

import (
	"github.com/gin-gonic/gin"
)

// MeInfoResponse json 回复
type MeInfoResponse struct {
	Code        int    `json:"code"`        // 返回代码
	Success     bool   `json:"success"`     // 验证成功
	Message     string `json:"message"`     // 消息
	Id          int    `json:"id"`          //id
	Time        string `json:"time"`        // 上次登录时间
	IP          string `json:"ip"`          // 上次登录IP
	Permissions int    `json:"permissions"` // 权限
}

type MeResponse struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Message string `json:"message"` // 消息
}

type MeUpdateResponse struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Message string `json:"message"` // 消息
}

// AddMeInfoResponse 添加回复
func AddMeInfoResponse(code int, success bool, message string, id int, time string, ip string, permissions int) MeInfoResponse {
	return MeInfoResponse{
		Code:        code,
		Success:     success,
		Message:     message,
		Id:          id,
		Time:        time,
		IP:          ip,
		Permissions: permissions,
	}
}

func AddMeUpdateResponse(code int, success bool, message string) MeUpdateResponse {
	return MeUpdateResponse{
		Code:    code,
		Success: success,
		Message: message,
	}
}

func UpdateMeInfo(id int, name string, password string) MeUpdateResponse {
	if id == 0 && Users[id].Name != name {
		return AddMeUpdateResponse(423, false, "不可修改根管理员用户名，如需修改请在服务器上修改文件")
	}
	for _, user := range Users {
		if name == user.Name && user.Id != id {
			return AddMeUpdateResponse(409, false, "重复用户名")
		}
	}
	Users[id].Name = name
	Users[id].Password = password
	if !SaveInfo(id) {
		return AddMeUpdateResponse(500, false, "修改失败，请重试")
	}
	return AddMeUpdateResponse(200, true, "修改成功")
}

func HandleMe(ctx *gin.Context) {
	var response MeResponse
	token := ctx.Request.Header.Get("Authorization")

	success, _ := ChecKToken(token)

	if success {
		response = MeResponse{
			Code:    200,
			Success: true,
			Message: "认证成功",
		}
	} else {
		response = MeResponse{
			Code:    403,
			Success: false,
			Message: "身份令牌过期，请重新登录",
		}

	}

	ctx.JSON(200, response)
}

func HandleMeUpdate(ctx *gin.Context) {
	var response MeUpdateResponse
	name := ctx.PostForm("name")
	password := ctx.PostForm("pwd")
	token := ctx.Request.Header.Get("Authorization")

	success, id := ChecKToken(token)

	if success {
		response = UpdateMeInfo(id, name, password)
	} else {
		response = AddMeUpdateResponse(403, false, "身份令牌过期，请重新登录")
	}

	ctx.JSON(200, response)
}

func HandleMeInfo(ctx *gin.Context) {
	var response MeInfoResponse
	token := ctx.Request.Header.Get("Authorization")

	success, id := ChecKToken(token)

	if success {
		response = AddMeInfoResponse(200, true, "认证成功", id, GetTime(Users[id].LastTime), Users[id].LastIp, Users[id].Permissions)
	} else {
		response = AddMeInfoResponse(403, false, "身份令牌过期，请重新登录", -1, "null", "null", 0)
	}

	ctx.JSON(200, response)
}
