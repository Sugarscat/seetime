package account

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type UsersList struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Identity    bool   `json:"identity"`
	Permissions int    `json:"permissions"`
}

type Data struct {
	Total   int         `json:"total"`
	Content []UsersList `json:"content"`
}

type UsersListResponse struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Message string `json:"message"` // 消息
	Data    Data   `json:"data"`
}

type UserResponse struct {
	Code      int    `json:"code"`    // 返回代码
	Success   bool   `json:"success"` // 验证成功
	Message   string `json:"message"` // 消息
	RequestId int    `json:"request"` // 请求的id
	UserId    int    `json:"id"`      // 添加/修改/删除的id
}

func addUsersList() []UsersList {
	var usersList = make([]UsersList, 0, 1)
	for _, user := range Users {
		aUser := UsersList{
			Id:          user.Id,
			Name:        user.Name,
			Identity:    user.Identity,
			Permissions: user.Permissions,
		}
		usersList = append(usersList, aUser)
	}
	return usersList
}

func AddUsersListResponse(code int, success bool, message string, userslist []UsersList) UsersListResponse {
	return UsersListResponse{
		Code:    code,
		Success: success,
		Message: message,
		Data: Data{
			Total:   len(userslist),
			Content: userslist,
		},
	}
}

func checkUserName(name string, id int) bool {
	for _, user := range Users {
		if name == user.Name && id != user.Id {
			return false
		}
	}
	return true
}

func ReloadUsersInfo() {
	for id := range Users {
		Users[id].Id = id
	}
}

func HandleUsersDelete(ctx *gin.Context) {
	var response UsersListResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")

	success, requestId := ChecKToken(token)

	if success {
		if Users[requestId].Identity {

			if id < len(Users) && id > 0 {
				Users = append(Users[:id], Users[id+1:]...)
				ReloadUsersInfo()
				if SaveInfo(-1) {
					response = AddUsersListResponse(200, true, "删除成功", addUsersList())
				} else {
					response = AddUsersListResponse(500, false, "删除失败，请重试", addUsersList())
				}
			} else if id == 0 {
				response = AddUsersListResponse(423, false, "不可删除根管理员", addUsersList())
			} else {
				response = AddUsersListResponse(404, false, "未找到该用户", addUsersList())
			}
		} else {
			response = AddUsersListResponse(400, false, "无权限", nil)
		}

	} else {
		response = AddUsersListResponse(403, false, "身份令牌过期，请重新登录", nil)
	}

	ctx.JSON(200, response)
}

func HandleUsersUpdate(ctx *gin.Context) {
	var response UsersListResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")
	name := ctx.PostForm("name")
	password := ctx.PostForm("pwd")
	identity, _ := strconv.ParseBool(ctx.PostForm("identity"))
	permissions, _ := strconv.Atoi(ctx.PostForm("permissions"))

	success, requestId := ChecKToken(token)

	if success {
		if Users[requestId].Identity && id != 0 {
			if id < len(Users) && id > -1 {
				if checkUserName(name, id) {
					Users[id].Name = name
					Users[id].Password = password
					Users[id].Identity = identity
					Users[id].Permissions = permissions
					if SaveInfo(-1) {
						response = AddUsersListResponse(200, true, "修改成功", addUsersList())
					} else {
						response = AddUsersListResponse(500, false, "修改失败，请重试", addUsersList())
					}
				} else {
					response = AddUsersListResponse(409, false, "修改失败，重复用户名", addUsersList())
				}
			} else {
				response = AddUsersListResponse(404, false, "未找到该用户", addUsersList())
			}
		} else {
			response = AddUsersListResponse(400, false, "无权限", nil)
		}

	} else {
		response = AddUsersListResponse(403, false, "身份令牌过期，请重新登录", nil)
	}

	ctx.JSON(200, response)
}

func HandleUsersAdd(ctx *gin.Context) {
	var response UsersListResponse
	token := ctx.Request.Header.Get("Authorization")
	name := ctx.PostForm("name")
	password := ctx.PostForm("pwd")
	identity, _ := strconv.ParseBool(ctx.PostForm("identity"))
	permissions, _ := strconv.Atoi(ctx.PostForm("permissions"))

	success, id := ChecKToken(token)

	if success {
		if Users[id].Identity {
			if checkUserName(name, -51) {
				user := User{
					Id:          len(Users),
					Name:        name,
					Password:    password,
					Token:       "null",
					Identity:    identity,
					LastIp:      "",
					ClientIp:    "",
					LastTime:    0,
					LoginTime:   0,
					Permissions: permissions,
				}
				Users = append(Users, user)
				if SaveInfo(-1) {
					response = AddUsersListResponse(200, true, "添加成功", addUsersList())
				} else {
					response = AddUsersListResponse(500, false, "添加失败，请重试", addUsersList())
				}
			} else {
				response = AddUsersListResponse(409, false, "添加失败，重复用户名", addUsersList())
			}

		} else {
			response = AddUsersListResponse(400, false, "无权限", nil)
		}

	} else {
		response = AddUsersListResponse(403, false, "身份令牌过期，请重新登录", nil)
	}

	ctx.JSON(200, response)
}

func HandleUsers(ctx *gin.Context) {
	var response UsersListResponse
	token := ctx.Request.Header.Get("Authorization")

	success, requestId := ChecKToken(token)

	if success {
		if Users[requestId].Identity {
			response = AddUsersListResponse(200, true, "加载成功", addUsersList())
		} else {
			response = AddUsersListResponse(400, false, "无权限", nil)
		}

	} else {
		response = AddUsersListResponse(403, false, "身份令牌过期，请重新登录", nil)
	}

	ctx.JSON(200, response)
}
