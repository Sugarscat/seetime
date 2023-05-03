package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"seetime/server/account"
	"seetime/server/module"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

// 回复信息
type TaskInfoData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Info    string `json:"info"`
	Success bool   `json:"success"`
	Diy     bool   `json:"diy"`
	Cycle   string `json:"cycle"`
	Lastime string `json:"lastime"`
	Command string `json:"command"`
}

type TaskResponse struct {
	Code    int          `json:"code"`    // 返回代码
	Success bool         `json:"success"` // 验证成功
	Message string       `json:"message"` // 消息
	Data    TaskInfoData `json:"data"`
}

// SaveTaskInfo 保存任务信息
func SaveTaskInfo(task TaskData) bool {
	file, _ := os.OpenFile(Tasks[task.Id].Location+"config.json", os.O_WRONLY|os.O_CREATE, 0644)
	defer func(file *os.File) {
		file.Close()
	}(file)

	json, _ := json.Marshal(task)

	file.Truncate(0)
	_, err := io.WriteString(file, string(json))
	if err != nil {
		fmt.Println(err) // ---日志
		return false
	}
	return true
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
				Diy:     false,
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
			Name:    TaskInfo.Name,
			Info:    TaskInfo.Info,
			Success: Tasks[id].Success,
			Diy:     TaskInfo.Diy,
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
	var response TaskResponse
	token := ctx.Request.Header.Get("Authorization")
	id, _ := strconv.Atoi(ctx.Query("id"))
	name := ctx.PostForm("name")
	info := ctx.PostForm("info")
	cycle := ctx.PostForm("cycle")
	command := ctx.PostForm("command")
	file, err := ctx.FormFile("file")

	success, requestId := account.ChecKToken(token)

	if success {
		if account.ParsingPermissions(requestId, "changeTask") {
			if id < len(Tasks) && id > -1 {
				task := TaskData{
					Id:      id,
					Name:    name,
					Info:    info,
					Cycle:   cycle,
					Command: command,
					File:    file.Filename,
				}
				UpdateCron(task, file, err, ctx)
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

func UpdateCron(task TaskData, file *multipart.FileHeader, err error, ctx *gin.Context) {
	taskInfo := ReadTaskInfo(task.Id)

	if err == nil {
		// 删除文件
		os.Remove(Tasks[taskInfo.Id].Location + task.File)
		// 将文件保存到服务器
		filepath := filepath.Join(Tasks[taskInfo.Id].Location, file.Filename)
		ctx.SaveUploadedFile(file, filepath)
	}

	taskInfo.Name = task.Name
	taskInfo.Info = task.Info
	taskInfo.Cycle = task.Cycle
	taskInfo.Command = task.Command
	taskInfo.File = task.File

	cron := cron.New()
	cron.AddFunc(taskInfo.Cycle, func() {
		if !Tasks[task.Id].Success {
			return // 发现上次未执行成功，跳过执行任务
		}
		var fileLog *os.File
		// 只保存两天的日志
		timeLog := time.Now().Unix() - taskInfo.Logtime
		if timeLog >= 172800 || timeLog == 0 {
			taskInfo.Logtime = time.Now().Unix()
			fileLog, _ = os.OpenFile(Tasks[task.Id].Location+"log.log", os.O_WRONLY|os.O_CREATE, 0644)
			fileLog.Truncate(0)
		} else {
			fileLog, _ = os.OpenFile(Tasks[task.Id].Location+"log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		}
		// 记录日志
		defer func(fileLog *os.File) {
			fileLog.Close()
		}(fileLog)

		// 设置格式
		log.SetOutput(fileLog)
		log.SetFlags(log.Ldate | log.Ltime)
		// 任务开始
		run := "cd " + Tasks[task.Id].Location + " && " + taskInfo.Command
		cmd := exec.Command(runStart, runCode, run)
		output, err := cmd.Output()
		log.Println(string(output))
		if err != nil {
			log.Println(err)
			Tasks[task.Id].Success = false
		}
		taskInfo.Lastime = time.Now().Unix()
		SaveTaskInfo(taskInfo)
	})
	cron.Start()
	Crons[task.Id] = *cron // 更新定时器信息
}
