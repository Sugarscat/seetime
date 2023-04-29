package tasks

import "github.com/gin-gonic/gin"

type TasksList struct {
}

type TasksListData struct {
	Total   int         `json:"total"`
	Content []TasksList `json:"content"`
}

type TasksResponse struct {
	Code    int           `json:"code"`    // 返回代码
	Success bool          `json:"success"` // 验证成功
	Message string        `json:"message"` // 消息
	Data    TasksListData `json:"data"`
}

func HandleTasks(ctx *gin.Context) {
	var response TasksList

	ctx.JSON(200, response)
}
