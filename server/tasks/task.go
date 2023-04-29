package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"seetime/server/account"
	"seetime/server/module"

	"github.com/gin-gonic/gin"
)

type TaskResponse struct {
	Code    int          `json:"code"`    // 返回代码
	Success bool         `json:"success"` // 验证成功
	Message string       `json:"message"` // 消息
	Data    TaskInfoData `json:"data"`
}

// SaveTaskInfo 保存任务信息
func SaveTaskInfo(task TaskData) {
	file, err := os.OpenFile(Tasks[task.Id].Location+"config.json", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err) // ---日志
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err) // ---日志
		}
	}(file)

	json, err := json.Marshal(task)
	if err != nil {
		fmt.Println(err) // ---日志
	}

	file.Truncate(0)
	_, err = io.WriteString(file, string(json))
	if err != nil {
		fmt.Println(err) // ---日志
	}
}

func ReadTaskInfo(id int) TaskData {
	var TaskInfo TaskData
	// 读取任务配置
	taskFile, _ := os.ReadFile(Tasks[id].Location + "config.json")
	err := json.Unmarshal(taskFile, &TaskInfo)
	if err != nil {
		fmt.Println(err) // ---日志
	}
	return TaskInfo
}

func AddTaskResponse(code int, success bool, message string, id int) TaskResponse {
	TaskInfo := ReadTaskInfo(id)
	if id == -1 {
		return TaskResponse{
			Code:    code,
			Success: success,
			Message: message,
			Data: TaskInfoData{
				Id:      id,
				Name:    "",
				Info:    "",
				Success: false,
				Cycle:   "",
				Lastime: "",
				Command: "",
			},
		}
	}
	return TaskResponse{
		Code:    code,
		Success: success,
		Message: message,
		Data: TaskInfoData{
			Id:      id,
			Name:    Tasks[id].Name,
			Info:    Tasks[id].Info,
			Success: Tasks[id].Success,
			Cycle:   TaskInfo.Cycle,
			Lastime: module.GetTime(TaskInfo.Lastime),
			Command: TaskInfo.Command,
		},
	}
}

func HandleTask(ctx *gin.Context) {
	var response TaskResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")

	success, requestId := account.ChecKToken(token)

	if success {
		if account.Users[requestId].Identity {
			response = AddTaskResponse(200, true, "加载成功", id)
		} else {
			response = AddTaskResponse(400, false, "无权限", -1)
		}
	} else {
		response = AddTaskResponse(403, false, "身份令牌过期，请重新登录", -1)
	}

	ctx.JSON(200, response)
}
