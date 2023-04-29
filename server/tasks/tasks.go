package tasks

import (
	"seetime/server/account"
	"seetime/server/module"

	"github.com/gin-gonic/gin"
)

type TasksList struct {
	Total   int            `json:"total"`
	Content []TaskInfoData `json:"content"`
}

type TasksResponse struct {
	Code    int       `json:"code"`    // 返回代码
	Success bool      `json:"success"` // 验证成功
	Message string    `json:"message"` // 消息
	Data    TasksList `json:"data"`
}

func addTasksList() []TaskInfoData {
	var tasksList = make([]TaskInfoData, 0, 1)
	for _, task := range Tasks {
		TaskInfo := ReadTaskInfo(task.Id)
		taskOne := TaskInfoData{
			Id:      task.Id,
			Name:    task.Name,
			Info:    task.Info,
			Success: task.Success,
			Cycle:   TaskInfo.Cycle,
			Lastime: module.GetTime(TaskInfo.Lastime),
			Command: TaskInfo.Command,
		}
		tasksList = append(tasksList, taskOne)
	}
	return tasksList
}

func AddTasksResponse(code int, success bool, message string, tasksList []TaskInfoData) TasksResponse {
	return TasksResponse{
		Code:    code,
		Success: success,
		Message: message,
		Data: TasksList{
			Total:   len(tasksList),
			Content: tasksList,
		},
	}
}

func HandleTasks(ctx *gin.Context) {
	var response TasksResponse
	token := ctx.Request.Header.Get("Authorization")

	success, requestId := account.ChecKToken(token)

	if success {
		if account.Users[requestId].Identity {
			response = AddTasksResponse(200, true, "加载成功", addTasksList())
		} else {
			response = AddTasksResponse(400, false, "无权限", nil)
		}

	} else {
		response = AddTasksResponse(403, false, "身份令牌过期，请重新登录", nil)
	}

	ctx.JSON(200, response)
}
