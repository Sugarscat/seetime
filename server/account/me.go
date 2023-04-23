package account

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func HandleMe(writer http.ResponseWriter, request *http.Request) {

	writer.Header().Set("Content-Type", "application/json")
	id := request.FormValue("id")
	token := request.FormValue("token")
	var response Response

	for _, user := range Users {
		if strconv.Itoa(user.Id) == id && user.Token == token {
			response = AddResponse(200, true, "认证成功", user.Id, user.Token, GetTime(user.LastTime), user.LastIp)
			break
		}
		response = AddResponse(404, false, "认证失败", -1, "error", "", "")
	}

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return
	}
	_, err = writer.Write(jsonBytes)
	if err != nil {
		return
	}
}
