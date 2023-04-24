package account

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// MeResponse json 回复
type MeResponse struct {
	Code        int    `json:"code"`    // 返回代码
	Success     bool   `json:"success"` // 验证成功
	Message     string `json:"message"` // 消息
	Id          int    `json:"id"`
	Token       string `json:"token"`
	Time        string `json:"time"` // 上次登录时间
	IP          string `json:"ip"`   // 上次登录IP
	Permissions int    `json:"permissions"`
}

// AddMeResponse 添加回复
func AddMeResponse(code int, success bool, message string, id int, token string, time string, ip string, permissions int) MeResponse {
	return MeResponse{
		Code:        code,
		Success:     success,
		Message:     message,
		Id:          id,
		Token:       token,
		Time:        time,
		IP:          ip,
		Permissions: permissions,
	}
}

func changeMeInfo(id int, name string, password string) MeResponse {
	if id == 0 && Users[id].Name != name {
		return AddMeResponse(403, false, "不可修改根管理员用户名", id, "", "", "", 00000)
	}
	Users[id].Name = name
	Users[id].Password = password
	if !SaveInfo(id) {
		return AddMeResponse(404, false, "修改失败", id, "", "", "", 00000)
	}
	return AddMeResponse(200, true, "修改成功", id, Users[id].Token, GetTime(Users[id].LastTime), Users[id].LastIp, Users[id].Permissions)
}

func HandleMe(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	id := request.FormValue("id")
	token := request.FormValue("token")

	var response MeResponse

	for _, user := range Users {
		if strconv.Itoa(user.Id) == id && user.Token == token {
			response = AddMeResponse(200, true, "认证成功", user.Id, user.Token, GetTime(user.LastTime), user.LastIp, user.Permissions)
			break
		}
		response = AddMeResponse(403, false, "身份令牌过期，请重新登录", -1, "null", "null", "null", 00000)
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

func HandleMeChange(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	id := request.FormValue("id")
	token := request.FormValue("token")
	name := request.FormValue("name")
	password := request.FormValue("password")
	var response MeResponse

	for _, user := range Users {
		if strconv.Itoa(user.Id) == id && user.Token == token {
			response = changeMeInfo(user.Id, name, password)
			break
		}
		response = AddMeResponse(403, false, "身份令牌过期，请重新登录", -1, "null", "null", "null", 00000)
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
