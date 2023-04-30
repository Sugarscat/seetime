package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/robfig/cron"
)

var (
	TasksInfo []byte
	tasksData tasksJson
	Tasks     []Task
	runStart  string
	runCode   string
	Crons     = make([]cron.Cron, 0, 1)
)

type Task struct {
	Id       int    `json:"id"`
	Success  bool   `json:"success"`
	Location string `json:"location"`
}

// 用于解析 json 文件
type tasksJson struct {
	Tasks []Task `json:"tasks"`
}

type TaskData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Info    string `json:"info"`
	Cycle   string `json:"cycle"`
	Command string `json:"command"`
	Lastime int64  `json:"lastime"`
	Logtime int64  `json:"logtime"`
}

// 回复信息
type TaskInfoData struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Info    string `json:"info"`
	Success bool   `json:"success"`
	Cycle   string `json:"cycle"`
	Lastime string `json:"lastime"`
	Command string `json:"command"`
}

func LoadTasks(tasksInfo []byte) {

	TasksInfo = tasksInfo

	defer addTasks()
}

func addTasks() {
	err := json.Unmarshal(TasksInfo, &tasksData)
	if err != nil {
		fmt.Println(err) // ---日志
	}

	Tasks = tasksData.Tasks

	defer PlanningTasks()
}

func PlanningTasks() {
	if runtime.GOOS == "windows" { // Windows
		runStart = "cmd"
		runCode = "/c"
	} else if runtime.GOOS == "darwin" { // MacOS
		runStart = "/bin/bash"
		runCode = "-c"
	} else { // Linux
		runStart = "sh"
		runCode = "-c"
	}

	for _, task := range Tasks {
		AddCron(task)
	}

}

func AddCron(task Task) {
	// 如任务上次未执行成功则跳过加载定时器
	if !task.Success {
		return
	}

	cron := cron.New()
	TaskInfo := ReadTaskInfo(task.Id)

	cron.AddFunc(TaskInfo.Cycle, func() {
		if !Tasks[task.Id].Success {
			return // 发现上次未执行成功，跳过执行任务
		}
		var fileLog *os.File
		// 只保存两天的日志
		timeLog := time.Now().Unix() - TaskInfo.Logtime
		if timeLog >= 172800 || timeLog == 0 {
			TaskInfo.Logtime = time.Now().Unix()
			fileLog, _ = os.OpenFile(task.Location+"log.log", os.O_WRONLY|os.O_CREATE, 0644)
			fileLog.Truncate(0)
		} else {
			fileLog, _ = os.OpenFile(task.Location+"log.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		}
		// 记录日志
		defer func(fileLog *os.File) {
			fileLog.Close()
		}(fileLog)

		// 设置格式
		log.SetOutput(fileLog)
		log.SetFlags(log.Ldate | log.Ltime)
		// 任务开始
		run := "cd " + task.Location + " && " + TaskInfo.Command
		cmd := exec.Command(runStart, runCode, run)
		output, err := cmd.Output()
		log.Println(string(output))
		if err != nil {
			log.Println(err)
			Tasks[task.Id].Success = false
		}
		TaskInfo.Lastime = time.Now().Unix()
		SaveTaskInfo(TaskInfo)
	})
	cron.Start()
	Crons = append(Crons, *cron)
}

func DeleteCron(id int) {
	Crons[id].Stop()                            // 停止定时器
	Crons = append(Crons[:id], Crons[id+1:]...) // 删除定时器
}
