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
	Code      int    `json:"code"`    // 返回代码
	Success   bool   `json:"success"` // 验证成功
	Message   string `json:"message"` // 消息
	RequestId int    `json:"request"` // 请求的id
	userId    int    `json:"id"`      // 添加/修改/删除的id
}

var userList = make([]UsersList, 0, 1)

func addUserList() []UsersList {
	for _, user := range Users {
		aUser := UsersList{
			Id:       user.Id,
			Name:     user.Name,
			Identity: user.Identity,
		}
		userList = append(userList, aUser)
	}
	return userList
}

func AddUserResponse(code int, success bool, message string, request int, id int) UserResponse {
	return UserResponse{
		Code:      code,
		Success:   success,
		Message:   message,
		RequestId: request,
		userId:    id,
	}
}

func ChangeUsersInfo(code string, id int, userid int, name string, pwd string, identity string, permissions string) UserResponse {
	var response UserResponse

	switch code {
	case "1": // 增

	case "2": // 删

	case "3": // 改

	default:
		response = UserResponse{
			Code:      404,
			Success:   false,
			Message:   "未知代码",
			RequestId: id,
			userId:    userid,
		}
	}

	return response
}

func HandleUser(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	id := request.FormValue("id")
	token := request.FormValue("token")
	var response UsersResponse

	for _, user := range Users {
		if strconv.Itoa(user.Id) == id && user.Token == token {
			if !user.Identity {
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
				UsersList: addUserList(),
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

func HandleUserManage(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	code := request.FormValue("code")
	id := request.FormValue("id")
	token := request.FormValue("token")
	userId := request.FormValue("userid")
	name := request.FormValue("name")
	password := request.FormValue("password")
	identity := request.FormValue("identity")
	permissions := request.FormValue("permissions")
	var response UserResponse

	for _, user := range Users {
		if strconv.Itoa(user.Id) == id && user.Token == token {
			userid, _ := strconv.Atoi(userId)
			if !user.Identity {
				response = AddUserResponse(403, false, "无权限", user.Id, userid)
				break
			}
			response = ChangeUsersInfo(code, user.Id, userid, name, password, identity, permissions)
			break
		}
		response = AddUserResponse(403, false, "身份令牌过期，请重新登录", user.Id, -1)
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
