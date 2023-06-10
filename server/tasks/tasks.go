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

// 任务列表中单个任务的信息
type TaskOneInfo struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Info    string `json:"info"`
	Diy     bool   `json:"diy"`
	Run     bool   `json:"run"`
	Success bool   `json:"success"`
	Lastime string `json:"lastime"`
}

type TasksList struct {
	Total   int           `json:"total"`
	Content []TaskOneInfo `json:"content"`
}

type TasksResponse struct {
	Code    int       `json:"code"`    // 返回代码
	Success bool      `json:"success"` // 验证成功
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
			Id:   task.Id,
			Path: task.Path,
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

func addTasksList() []TaskOneInfo {
	var tasksList = make([]TaskOneInfo, 0, 1)
	var tasksNotRun = make([]TaskOneInfo, 0, 1)
	var tasksSuccess = make([]TaskOneInfo, 0, 1)
	var tasksFail = make([]TaskOneInfo, 0, 1)
	for _, task := range Tasks {
		TaskInfo := ReadTaskInfo(task.Id)
		taskOne := TaskOneInfo{
			Id:      task.Id,
			Name:    TaskInfo.Name,
			Info:    TaskInfo.Info,
			Diy:     TaskInfo.Diy,
			Run:     TaskInfo.Run,
			Success: TaskInfo.Success,
			Lastime: module.GetTime(TaskInfo.Lastime),
		}
		tasksList = append(tasksList, taskOne)
	}
	for _, task := range tasksList {
		if !task.Run {
			tasksNotRun = append(tasksNotRun, task)
		} else {
			if task.Success {
				tasksSuccess = append(tasksSuccess, task)
			} else {
				tasksFail = append(tasksFail, task)
			}
		}
	}
	tasksList = make([]TaskOneInfo, 0, 1)
	tasksList = append(tasksList, tasksFail...)
	tasksList = append(tasksList, tasksSuccess...)
	tasksList = append(tasksList, tasksNotRun...)
	return tasksList
}

func AddTasksResponse(code int, success bool, tasksList []TaskOneInfo) TasksResponse {
	return TasksResponse{
		Code:    code,
		Success: success,
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

func anyPermissions(requestId int) bool {
	if account.ParsingPermissions(requestId, "situation") {
		return true
	}
	if account.ParsingPermissions(requestId, "addTask") {
		return true
	}
	if account.ParsingPermissions(requestId, "updateTask") {
		return true
	}
	if account.ParsingPermissions(requestId, "deleteTask") {
		return true
	}
	if account.ParsingPermissions(requestId, "exportTask") {
		return true
	}
	return false
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
				"data":    count.Count,
			})
		} else {
			ctx.JSON(200, gin.H{
				"code":    400,
				"success": false,
				"data":    nil,
			})
		}
		return
	}

	ctx.JSON(200, gin.H{
		"code":    403,
		"success": false,
		"data":    nil,
	})
}

// HandleTasks 回复任务列表
func HandleTasks(ctx *gin.Context) {
	var response TasksResponse
	token := ctx.Request.Header.Get("Authorization")

	success, requestId := account.ChecKToken(token)

	if success {
		if anyPermissions(requestId) {
			response = AddTasksResponse(200, true, addTasksList())
		} else {
			response = AddTasksResponse(400, false, nil)
		}

	} else {
		response = AddTasksResponse(403, false, nil)
	}

	ctx.JSON(200, response)
}

func HandleTasksAdd(ctx *gin.Context) {

	if !bucket.AddWater(1) {
		ctx.JSON(200, gin.H{
			"code":    429,
			"success": false,
			"data":    "null",
		})
		return
	}

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
				Id:   len(Tasks),
				Path: "./resources/tasks/" + strconv.FormatInt(time.Now().Unix(), 10) + "/",
			}
			taskData := TaskData{
				Name:    name,
				Info:    info,
				Diy:     true,
				Run:     true,
				Success: true,
				Cycle:   cycle,
				Command: command,
				File:    file.Filename,
				Lastime: 0,
			}

			err := os.MkdirAll(task.Path, 0755)
			if err != nil {
				response = AddTasksResponse(500, false, addTasksList())
				ctx.JSON(200, response)
				return
			}

			// 将文件保存到服务器
			filepath := filepath.Join(task.Path, file.Filename)
			ctx.SaveUploadedFile(file, filepath)

			Tasks = append(Tasks, task) // 添加任务

			if SaveTasks() && SaveTaskInfo(task.Id, taskData) {
				AddCron(task.Id)
				response = AddTasksResponse(200, true, addTasksList())
			} else {
				// 添加失败后删除上传的信息
				Tasks = append(Tasks[:task.Id], Tasks[task.Id+1:]...)
				os.RemoveAll(task.Path)
				response = AddTasksResponse(500, false, addTasksList())
			}

		} else {
			response = AddTasksResponse(400, false, nil)
		}

	} else {
		response = AddTasksResponse(403, false, nil)
	}

	ctx.JSON(200, response)
}

func HandleTasksDelete(ctx *gin.Context) {

	if !bucket.AddWater(1) {
		ctx.JSON(200, gin.H{
			"code":    429,
			"success": false,
			"data":    "null",
		})
		return
	}

	var response TasksResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")
	success, requestId := account.ChecKToken(token)
	if success {
		if account.ParsingPermissions(requestId, "deleteTask") {
			if id < len(Tasks) && id > -1 {
				lastTask := Tasks[id]
				Tasks = append(Tasks[:id], Tasks[id+1:]...)
				ReloadTasksInfo()
				if SaveTasks() {
					// 停止所有定时器
					for _, cron := range Crons {
						cron.Stop()
					}
					os.RemoveAll(lastTask.Path)
					response = AddTasksResponse(200, true, addTasksList())
					// 重载所有定时器
					Crons = Crons[:0]
					PlanningTasks()
				} else {
					// 若保存失败则回档
					newSlice := make([]Task, len(Tasks)+1)
					copy(newSlice[:id], Tasks[:id])
					newSlice[id] = lastTask
					copy(newSlice[id+1:], Tasks[id:])
					Tasks = newSlice
					ReloadTasksInfo()
					response = AddTasksResponse(500, false, addTasksList())
				}
			} else {
				response = AddTasksResponse(404, false, addTasksList())
			}
		} else {
			response = AddTasksResponse(400, false, nil)
		}
	} else {
		response = AddTasksResponse(403, false, nil)
	}
	ctx.JSON(200, response)
}

func HandleTaskStop(ctx *gin.Context) {
	var response TasksResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")
	success, requestId := account.ChecKToken(token)
	if success {
		// 拥有添加任务或修改任务权限的任意一个
		if account.ParsingPermissions(requestId, "addTask") || account.ParsingPermissions(requestId, "updateTask") {
			if id < len(Tasks) && id > -1 {
				StopTask(id)
				response = AddTasksResponse(200, true, addTasksList())
			} else {
				response = AddTasksResponse(404, false, addTasksList())
			}
		} else {
			response = AddTasksResponse(400, false, nil)
		}
	} else {
		response = AddTasksResponse(403, false, nil)
	}
	ctx.JSON(200, response)
}

func HandleTaskActivate(ctx *gin.Context) {
	var response TasksResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")
	success, requestId := account.ChecKToken(token)
	if success {
		// 拥有添加任务或修改任务权限的任意一个
		if account.ParsingPermissions(requestId, "addTask") || account.ParsingPermissions(requestId, "updateTask") {
			if id < len(Tasks) && id > -1 {
				ActivateTask(id)
				response = AddTasksResponse(200, true, addTasksList())
			} else {
				response = AddTasksResponse(404, false, addTasksList())
			}
		} else {
			response = AddTasksResponse(400, false, nil)
		}
	} else {
		response = AddTasksResponse(403, false, nil)
	}
	ctx.JSON(200, response)
}

func HandleTasksRunOne(ctx *gin.Context) {
	var response TasksResponse
	token := ctx.Request.Header.Get("Authorization")
	id, _ := strconv.Atoi(ctx.Query("id"))

	success, requestId := account.ChecKToken(token)

	if success {
		// 拥有添加任务或修改任务权限的任意一个
		if account.ParsingPermissions(requestId, "addTask") || account.ParsingPermissions(requestId, "updateTask") {
			if id < len(Tasks) && id > -1 {
				if RunTask(id) {
					response = AddTasksResponse(200, true, addTasksList())
				} else {
					response = AddTasksResponse(500, false, addTasksList())
				}
			} else {
				response = AddTasksResponse(404, false, addTasksList())
			}
		} else {
			response = AddTasksResponse(400, false, nil)
		}
	} else {
		response = AddTasksResponse(403, false, nil)
	}

	ctx.JSON(200, response)
}
