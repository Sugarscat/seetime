package account

import (
	"github.com/sugarscat/seetime/server/module"

	"github.com/gin-gonic/gin"
)

type MeInfoData struct {
	Id          int    `json:"id"` // id
	Identity    bool   `json:"identity"`
	Time        string `json:"time"`        // 上次登录时间
	IP          string `json:"ip"`          // 上次登录IP
	Permissions int    `json:"permissions"` // 权限
}

// MeInfoResponse json 回复
type MeInfoResponse struct {
	Code    int        `json:"code"`    // 返回代码
	Success bool       `json:"success"` // 验证成功
	Message string     `json:"message"` // 消息
	Data    MeInfoData `json:"data"`
}

type MeUpdateResponse struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Message string `json:"message"` // 消息
}

// AddMeInfoResponse 添加回复
func AddMeInfoResponse(code int, success bool, message string, data MeInfoData) MeInfoResponse {
	return MeInfoResponse{
		Code:    code,
		Success: success,
		Message: message,
		Data:    data,
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
	lastName := Users[id].Name
	lastPassword := Users[id].Password
	if id == 0 && Users[id].Name != name { // 不可修改根管理用户名，防呆设计
		return AddMeUpdateResponse(423, false, "不可修改根管理员用户名，如需修改请在服务器上修改文件")
	}
	for _, user := range Users {
		if name == user.Name && user.Id != id {
			return AddMeUpdateResponse(409, false, "重复用户名")
		}
	}
	Users[id].Name = name
	// 判断用户是否传入空密码，若不是则改变密码
	if len(password) != 0 {
		Users[id].Password = password
	}
	if !SaveInfo(id) {
		// 若保存失败则回档
		Users[id].Name = lastName
		Users[id].Password = lastPassword
		return AddMeUpdateResponse(500, false, "修改失败，请重试")
	}
	return AddMeUpdateResponse(200, true, "修改成功")
}

func HandleMe(ctx *gin.Context) {
	var code int
	var response MeInfoResponse
	token := ctx.Request.Header.Get("Authorization")

	success, id := ChecKToken(token)

	if success {
		code = 200
		response = AddMeInfoResponse(200, true, "认证成功", MeInfoData{
			id,
			Users[id].Identity,
			module.GetTime(Users[id].LastTime),
			Users[id].LastIp,
			Users[id].Permissions,
		})
	} else {
		code = 403
		response = AddMeInfoResponse(403, false, "身份令牌过期，请重新登录", MeInfoData{
			-1,
			false,
			"---",
			"---",
			0,
		})
	}

	ctx.JSON(code, response)
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
