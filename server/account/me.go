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
	Data    MeInfoData `json:"data"`
}

type MeUpdateResponse struct {
	Code    int  `json:"code"`    // 返回代码
	Success bool `json:"success"` // 验证成功
}

func UpdateMeInfo(id int, name string, password string) MeUpdateResponse {
	lastName := Users[id].Name
	lastPassword := Users[id].Password
	for _, user := range Users {
		if name == user.Name && user.Id != id {
			return MeUpdateResponse{409, false}
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
		return MeUpdateResponse{500, false}
	}
	return MeUpdateResponse{200, true}
}

// HandleMe 获取个人资料
func HandleMe(ctx *gin.Context) {
	var code int
	var response MeInfoResponse
	token := ctx.Request.Header.Get("Authorization")

	success, id := ChecKToken(token)

	if success {
		code = 200
		response = MeInfoResponse{200, true, MeInfoData{
			id,
			Users[id].Identity,
			module.GetTime(Users[id].LastTime),
			Users[id].LastIp,
			Users[id].Permissions,
		}}
	} else {
		code = 403
		response = MeInfoResponse{403, false, MeInfoData{
			-1,
			false,
			"",
			"",
			0,
		}}
	}

	ctx.JSON(code, response)
}

// HandleMeUpdate 更新个人资料
func HandleMeUpdate(ctx *gin.Context) {
	var response MeUpdateResponse
	name := ctx.PostForm("name")
	password := ctx.PostForm("pwd")
	token := ctx.Request.Header.Get("Authorization")

	success, id := ChecKToken(token)

	if success {
		response = UpdateMeInfo(id, name, password)
	} else {
		response = MeUpdateResponse{403, false}
	}

	ctx.JSON(200, response)
}
