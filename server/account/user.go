package account

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type UsersList struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Identity bool   `json:"identity"`
}

type UsersResponse struct {
	Code      int         `json:"code"`    // 返回代码
	Success   bool        `json:"success"` // 验证成功
	Message   string      `json:"message"` // 消息
	UsersList []UsersList `json:"userslist"`
}

type UserResponse struct {
	Code    int    `json:"code"`    // 返回代码
	Success bool   `json:"success"` // 验证成功
	Message string `json:"message"` // 消息
	Request int    `json:"request"` // 请求的id
	Id      int    `json:"id"`      // 添加/修改/删除的id
}

func AddUserResponse(code int, success bool, message string, request int, id int) UserResponse {
	return UserResponse{
		Code:    code,
		Success: success,
		Message: message,
		Request: request,
		Id:      id,
	}
}

func HandleUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	id := request.FormValue("id")
	token := request.FormValue("token")
	var response UsersResponse

	for _, user := range Users {
		if strconv.Itoa(user.Id) == id && user.Token == token {
			if user.Identity {
				response = UsersResponse{
					Code:      403,
					Success:   false,
					Message:   "无权限",
					UsersList: nil,
				}
				break
			}
			response = UsersResponse{
				Code:      200,
				Success:   true,
				Message:   "加载成功",
				UsersList: nil,
			}
			break
		}
		response = UsersResponse{
			Code:      403,
			Success:   false,
			Message:   "身份令牌过期，请重新登录",
			UsersList: nil,
		}
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		// ---日志
		return
	}
	_, err = writer.Write(jsonBytes)
	if err != nil {
		// ---日志
		return
	}
}

func HandleUserAdd(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	id := request.FormValue("id")
	token := request.FormValue("token")
	name := request.FormValue("name")
	password := request.FormValue("password")
	// identity := request.FormValue("identity")
	var response MeResponse

	for _, user := range Users {
		if strconv.Itoa(user.Id) == id && user.Token == token {
			if user.Identity {
				response = AddMeResponse(403, false, "无权限", -1, "error", "", "", EmptyPermissions)
				break
			}
			response = changeMeInfo(user.Id, name, password)
			break
		}
		response = AddMeResponse(403, false, "身份令牌过期，请重新登录", -1, "error", "", "", EmptyPermissions)
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		// ---日志
		return
	}
	_, err = writer.Write(jsonBytes)
	if err != nil {
		// ---日志
		return
	}
}
