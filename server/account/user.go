package account

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserData struct {
	Id          int    `json:"id"` // id
	Name        string `json:"name"`
	Identity    bool   `json:"identity"`
	Permissions int    `json:"permissions"` // 权限
}

type UserResponse struct {
	Code    int      `json:"code"`    // 返回代码
	Success bool     `json:"success"` // 验证成功
	Message string   `json:"message"` // 消息
	Data    UserData `json:"data"`
}

func AddUserResponse(code int, success bool, message string, id int) UserResponse {
	var response UserResponse
	if id == -1 {
		response = UserResponse{
			Code:    code,
			Success: success,
			Message: message,
			Data: UserData{
				Id:          id,
				Name:        "",
				Identity:    false,
				Permissions: 0,
			},
		}
		return response
	}
	response = UserResponse{
		Code:    code,
		Success: success,
		Message: message,
		Data: UserData{
			Id:          id,
			Name:        Users[id].Name,
			Identity:    Users[id].Identity,
			Permissions: Users[id].Permissions,
		},
	}
	return response
}

func HandleUser(ctx *gin.Context) {
	var response UserResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")

	success, requestId := ChecKToken(token)

	if success {
		if Users[requestId].Identity {
			if id < len(Users) && id > -1 {
				response = AddUserResponse(200, true, "查询成功", id)
			} else {
				response = AddUserResponse(404, false, "无此用户", -1)
			}
		} else {
			response = AddUserResponse(400, false, "无权限", -1)
		}
	} else {
		response = AddUserResponse(403, false, "身份令牌过期，请重新登录", -1)
	}

	ctx.JSON(200, response)
}
