package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sugarscat/seetime/server/account"
	"github.com/sugarscat/seetime/server/module"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

// 任务的本地数据
type TaskInfoData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Info    string `json:"info"`
	Success bool   `json:"success"`
	Diy     bool   `json:"diy"`
	Run     bool   `json:"run"`
	Cycle   string `json:"cycle"`
	Command string `json:"command"`
	Lastime string `json:"lastime"`
}

// 回复信息
type TaskResponse struct {
	Code    int          `json:"code"`    // 返回代码
	Success bool         `json:"success"` // 验证成功
	Message string       `json:"message"` // 消息
	Data    TaskInfoData `json:"data"`
}

// SaveTaskInfo 保存任务信息
func SaveTaskInfo(id int, task TaskData) bool {
	file, _ := os.OpenFile(Tasks[id].Location+"config.json", os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()

	json, _ := json.Marshal(task)

	file.Truncate(0)
	_, err := io.WriteString(file, string(json))
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}
	return true
}

// RunTask 运行任务
func RunTask(id int) {
	go func() {
		TaskInfo := ReadTaskInfo(id)

		defer func() { // 捕获异常，避免任务执行错误，导致整个系统挂掉
			if err := recover(); err != nil {
				fmt.Println(TaskInfo.Name, "任务执行错误")
			}
		}()

		// 任务开始
		run := "cd " + Tasks[id].Location + " && " + TaskInfo.Command
		cmd := exec.Command(runStart, runCode, run)
		err := cmd.Start()
		if err != nil {
			Tasks[id].Success = false
		} else {
			Tasks[id].Success = true
		}
		TaskInfo.Lastime = time.Now().Unix()
		SaveTaskInfo(id, TaskInfo)
	}()
}

// ReadTaskInfo 读取任务信息
func ReadTaskInfo(id int) TaskData {
	var TaskInfo TaskData
	// 读取任务配置
	taskFile, _ := os.ReadFile(Tasks[id].Location + "config.json")
	json.Unmarshal(taskFile, &TaskInfo)
	return TaskInfo
}

// ReadTaskLog 读取任务日志
func ReadTaskLog(id int) (data string, success bool) {
	log, err := os.ReadFile(Tasks[id].Location + "log.log")
	if err != nil {
		return "null", false
	} else {
		return string(log), true
	}
}

// StopTask 停止任务
func StopTask(id int) {
	Crons[id].Stop()
	task := ReadTaskInfo(id)
	task.Run = false
	SaveTaskInfo(id, task)
}

func ActivateTask(id int) {
	task := ReadTaskInfo(id)
	task.Run = true
	SaveTaskInfo(id, task)
	Crons[id].Start()
}

func AddTaskResponse(code int, success bool, message string, id int) TaskResponse {
	if id == -1 {
		return TaskResponse{
			Code:    code,
			Success: success,
			Message: message,
			Data: TaskInfoData{
				Id: id,
			},
		}
	}
	TaskInfo := ReadTaskInfo(id)
	return TaskResponse{
		Code:    code,
		Success: success,
		Message: message,
		Data: TaskInfoData{
			Id:      id,
			Name:    TaskInfo.Name,
			Info:    TaskInfo.Info,
			Success: Tasks[id].Success,
			Diy:     TaskInfo.Diy,
			Run:     TaskInfo.Run,
			Cycle:   TaskInfo.Cycle,
			Command: TaskInfo.Command,
			Lastime: module.GetTime(TaskInfo.Lastime),
		},
	}
}

// UpdateCron 更新定时器信息
func UpdateCron(id int, task TaskData, file *multipart.FileHeader, change bool, ctx *gin.Context) {
	Crons[id].Stop() // 停止定时器

	taskInfo := ReadTaskInfo(id)
	taskInfo.Name = task.Name
	taskInfo.Info = task.Info
	taskInfo.Cycle = task.Cycle
	taskInfo.Command = task.Command

	if change {
		// 删除文件
		os.Remove(Tasks[id].Location + task.File)
		// 将文件保存到服务器
		filepath := filepath.Join(Tasks[id].Location, file.Filename)
		ctx.SaveUploadedFile(file, filepath)
		taskInfo.File = file.Filename
	}

	SaveTaskInfo(id, taskInfo)
	cron := cron.New()
	cron.AddFunc(taskInfo.Cycle, func() {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(taskInfo.Name, "任务执行错误")
				}
			}()

			// 任务开始
			run := "cd " + Tasks[id].Location + " && " + taskInfo.Command
			cmd := exec.Command(runStart, runCode, run)
			err := cmd.Start()
			if err != nil {
				Tasks[id].Success = false
			} else {
				Tasks[id].Success = true
			}
			taskInfo.Lastime = time.Now().Unix()
			SaveTaskInfo(id, taskInfo)
		}()
	})
	if taskInfo.Run {
		cron.Start()
	}
	Crons[id] = *cron // 更新定时器信息
}

func HandleTask(ctx *gin.Context) {
	var response TaskResponse
	id, _ := strconv.Atoi(ctx.Query("id"))
	token := ctx.Request.Header.Get("Authorization")

	success, requestId := account.ChecKToken(token)

	if success {
		if account.Users[requestId].Identity {
			if id < len(Tasks) && id > -1 {
				response = AddTaskResponse(200, true, "查询成功", id)
			} else {
				response = AddTaskResponse(404, false, "无此任务", -1)
			}
		} else {
			response = AddTaskResponse(400, false, "无权限", -1)
		}
	} else {
		response = AddTaskResponse(403, false, "身份令牌过期，请重新登录", -1)
	}

	ctx.JSON(200, response)
}

func HandleTaskUpdate(ctx *gin.Context) {

	if !bucket.AddWater(1) {
		ctx.JSON(200, gin.H{
			"code":    429,
			"success": false,
			"message": "请求次数过多",
			"data":    "null",
		})
		return
	}

	var response TaskResponse
	var change bool
	token := ctx.Request.Header.Get("Authorization")
	id, _ := strconv.Atoi(ctx.Query("id"))
	name := ctx.PostForm("name")
	info := ctx.PostForm("info")
	cycle := ctx.PostForm("cycle")
	command := ctx.PostForm("command")
	file, err := ctx.FormFile("file")
	if err != nil {
		change = false
	} else {
		change = true
	}
	success, requestId := account.ChecKToken(token)

	if success {
		if account.ParsingPermissions(requestId, "changeTask") {
			if id < len(Tasks) && id > -1 {
				task := TaskData{
					Name:    name,
					Info:    info,
					Cycle:   cycle,
					Command: command,
					File:    "",
				}
				UpdateCron(id, task, file, change, ctx)
				response = AddTaskResponse(200, true, "修改成功", id)
			} else {
				response = AddTaskResponse(404, false, "无此任务", -1)
			}
		} else {
			response = AddTaskResponse(400, false, "无权限", -1)
		}
	} else {
		response = AddTaskResponse(403, false, "身份令牌过期，请重新登录", -1)
	}

	ctx.JSON(200, response)
}

func HandleTaskRunOne(ctx *gin.Context) {
	var response TaskResponse
	token := ctx.Request.Header.Get("Authorization")
	id, _ := strconv.Atoi(ctx.Query("id"))

	success, requestId := account.ChecKToken(token)

	if success {
		// 拥有添加任务或修改任务权限的任意一个
		if account.ParsingPermissions(requestId, "addTask") || account.ParsingPermissions(requestId, "changeTask") {
			if id < len(Tasks) && id > -1 {
				RunTask(id)
				if Tasks[id].Success {
					response = AddTaskResponse(200, true, "执行成功", id)
				} else {
					response = AddTaskResponse(500, false, "执行失败", id)
				}
			} else {
				response = AddTaskResponse(404, false, "无此任务", -1)
			}
		} else {
			response = AddTaskResponse(400, false, "无权限", -1)
		}
	} else {
		response = AddTaskResponse(403, false, "身份令牌过期，请重新登录", -1)
	}

	ctx.JSON(200, response)
}

func HandleTaskLog(ctx *gin.Context) {
	var response gin.H
	token := ctx.Request.Header.Get("Authorization")
	id, _ := strconv.Atoi(ctx.Query("id"))

	success, requestId := account.ChecKToken(token)

	if success {
		if account.ParsingPermissions(requestId, "situation") {
			if id < len(Tasks) && id > -1 {
				log, can := ReadTaskLog(id)
				if can {
					response = gin.H{
						"code":    200,
						"success": true,
						"message": "读取成功",
						"data":    log,
					}
				} else {
					response = gin.H{
						"code":    404,
						"success": false,
						"message": "没有日志",
						"data":    "没有日志",
					}
				}
			} else {
				response = gin.H{
					"code":    404,
					"success": false,
					"message": "无此任务",
					"data":    "null",
				}
			}
		} else {
			response = gin.H{
				"code":    400,
				"success": false,
				"message": "无权限",
				"data":    "null",
			}
		}
	} else {
		response = gin.H{
			"code":    403,
			"success": false,
			"message": "身份令牌过期，请重新登录",
			"data":    "null",
		}
	}

	ctx.JSON(200, response)
}
