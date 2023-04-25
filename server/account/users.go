package account

import (
	"github.com/gin-gonic/gin"
)

type UsersList struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Identity bool   `json:"identity"`
}

type UsersListResponse struct {
	Code      int         `json:"code"`    // 返回代码
	Success   bool        `json:"success"` // 验证成功
	Message   string      `json:"message"` // 消息
	UsersList []UsersList `json:"userslist"`
}

type UserResponse struct {
	Code      int    `json:"code"`    // 返回代码
	Success   bool   `json:"success"` // 验证成功
	Message   string `json:"message"` // 消息
	RequestId int    `json:"request"` // 请求的id
	UserId    int    `json:"id"`      // 添加/修改/删除的id
}

var usersList = make([]UsersList, 0, 1)

func addUsersList() []UsersList {
	for _, user := range Users {
		aUser := UsersList{
			Id:       user.Id,
			Name:     user.Name,
			Identity: user.Identity,
		}
		usersList = append(usersList, aUser)
	}
	return usersList
}

func AddUsersListResponse(code int, success bool, message string, userslist []UsersList) UsersListResponse {
	return UsersListResponse{
		Code:      code,
		Success:   success,
		Message:   message,
		UsersList: userslist,
	}
}

func AddUserResponse(code int, success bool, message string, request int, id int) UserResponse {
	return UserResponse{
		Code:      code,
		Success:   success,
		Message:   message,
		RequestId: request,
		UserId:    id,
	}
}

func HandleUsers(ctx *gin.Context) {
	var response UsersListResponse
	token := ctx.Request.Header.Get("Authorization")

	success, id := ChecKToken(token)

	if success {
		if Users[id].Identity {
			response = AddUsersListResponse(200, true, "加载成功", addUsersList())
		} else {
			response = AddUsersListResponse(403, false, "无权限", nil)
		}

	} else {
		response = AddUsersListResponse(403, false, "身份令牌过期，请重新登录", nil)
	}

	ctx.JSON(200, response)
}
