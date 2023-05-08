package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sugarscat/seetime/server/account"
	"github.com/sugarscat/seetime/server/module"

	"github.com/gin-gonic/gin"
)

type TasksListInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Info    string `json:"info"`
	Success bool   `json:"success"`
	Diy     bool   `json:"diy"`
	Lastime string `json:"lastime"`
}

type TasksList struct {
	Total   int             `json:"total"`
	Content []TasksListInfo `json:"content"`
}

type TasksResponse struct {
	Code    int       `json:"code"`    // 返回代码
	Success bool      `json:"success"` // 验证成功
	Message string    `json:"message"` // 消息
	Data    TasksList `json:"data"`
}

// SaveTask 保存任务
func SaveTasks() bool {
	file, _ := os.OpenFile("./data/tasks/tasks.json", os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()

	var taskJsonFile tasksJson
	taskJsonFile.Tasks = make([]Task, 0, 1)
	for _, task := range Tasks {
		taskJsonFile.Tasks = append(taskJsonFile.Tasks, Task{
			Id:       task.Id,
			Success:  task.Success,
			Location: task.Location,
		})
	}

	jsonDataU, _ := json.Marshal(taskJsonFile)

	file.Truncate(0)
	_, err := io.WriteString(file, string(jsonDataU))
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}

	return true
}

func addTasksList() []TasksListInfo {
	var tasksList = make([]TasksListInfo, 0, 1)
	for _, task := range Tasks {
		TaskInfo := ReadTaskInfo(task.Id)
		taskOne := TasksListInfo{
			Id:      task.Id,
			Name:    TaskInfo.Name,
			Info:    TaskInfo.Info,
			Success: task.Success,
			Diy:     TaskInfo.Diy,
			Lastime: module.GetTime(TaskInfo.Lastime),
		}
		tasksList = append(tasksList, taskOne)
	}
	return tasksList
}

func AddTasksResponse(code int, success bool, message string, tasksList []TasksListInfo) TasksResponse {
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

func ReloadTasksInfo() {
	for id := range Tasks {
		Tasks[id].Id = id
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

func HandleTasksAdd(ctx *gin.Context) {
	var response TasksResponse
	token := ctx.Request.Header.Get("Authorization")
	name := ctx.PostForm("name")
	info := ctx.PostForm("info")
	cycle := ctx.PostForm("cycle")
	command := ctx.PostForm("command")
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.String(400, "Bad request")
		return
	}

	success, requestId := account.ChecKToken(token)

	if success {
		if account.ParsingPermissions(requestId, "addTask") {
			task := Task{
				Id:       len(Tasks),
				Success:  true,
				Location: "./resources/tasks/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "/",
			}
			taskData := TaskData{
				Name:    name,
				Info:    info,
				Diy:     true,
				Cycle:   cycle,
				Command: command,
				File:    file.Filename,
				Lastime: 0,
			}

			err := os.MkdirAll(task.Location, 0755)
			if err != nil {
				response = AddTasksResponse(500, false, "添加失败，请重试", addTasksList())
				ctx.JSON(200, response)
				return
			}

			// 将文件保存到服务器
			filepath := filepath.Join(task.Location, file.Filename)
			ctx.SaveUploadedFile(file, filepath)

			Tasks = append(Tasks, task) // 添加任务

			if SaveTasks() && SaveTaskInfo(task.Id, taskData) {
				AddCron(task.Id)
				response = AddTasksResponse(200, true, "添加成功", addTasksList())
			} else {
				// 添加失败后删除上传的信息
				Tasks = append(Tasks[:task.Id], Tasks[task.Id+1:]...)
				os.RemoveAll(task.Location)
				response = AddTasksResponse(500, false, "添加失败，请重试", addTasksList())
			}

		} else {
			response = AddTasksResponse(400, false, "无权限", nil)
		}

	} else {
		response = AddTasksResponse(403, false, "身份令牌过期，请重新登录", nil)
	}

	ctx.JSON(200, response)
}

func HandleTasksDelete(ctx *gin.Context) {
	var response TasksResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")
	success, requestId := account.ChecKToken(token)
	if success {
		if account.ParsingPermissions(requestId, "deleteTask") {
			lastTask := Tasks[id]
			Tasks = append(Tasks[:id], Tasks[id+1:]...)
			ReloadTasksInfo()
			if SaveTasks() {
				os.RemoveAll(lastTask.Location)
				response = AddTasksResponse(200, true, "删除成功", addTasksList())
			} else {
				// 若保存失败则回档
				newSlice := make([]Task, len(Tasks)+1)
				copy(newSlice[:id], Tasks[:id])
				newSlice[id] = lastTask
				copy(newSlice[id+1:], Tasks[id:])
				Tasks = newSlice
				ReloadTasksInfo()
				response = AddTasksResponse(500, false, "删除失败，请重试", addTasksList())
			}
			for _, cron := range Crons {
				cron.Stop()
			}
			Crons = Crons[:0]
			PlanningTasks()
		} else {
			response = AddTasksResponse(400, false, "无权限", nil)
		}
	} else {
		response = AddTasksResponse(403, false, "身份令牌过期，请重新登录", nil)
	}
	ctx.JSON(200, response)
}

func HandleTasksCount(ctx *gin.Context) {
	var count tasksNumJson
	token := ctx.Request.Header.Get("Authorization")
	success, id := account.ChecKToken(token)

	if success {
		if account.ParsingPermissions(id, "situation") {
			countInfo, _ := os.ReadFile("./data/tasks/count.json")
			json.Unmarshal(countInfo, &count)
			ctx.JSON(200, gin.H{
				"code":    200,
				"success": true,
				"message": "加载成功",
				"data":    count.Count,
			})
		} else {
			ctx.JSON(200, gin.H{
				"code":    400,
				"success": false,
				"message": "无权限",
				"data":    nil,
			})
		}
		return
	}

	ctx.JSON(200, gin.H{
		"code":    403,
		"success": false,
		"message": "身份令牌过期，请重新登录",
		"data":    nil,
	})
}
