package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/robfig/cron"
	"github.com/sugarscat/seetime/server/module"
)

// 每 2 秒可请求一次 api, (为什么限制？因为本项目的所属者很菜，不限制会导致系统运行异常)
var bucket = module.NewLeakyBucket(1, 0.5) // 桶

var (
	TasksInfo []byte
	tasksData tasksJson
	Tasks     []Task
	runStart  string
	runCode   string
	Crons     = make([]cron.Cron, 0, 1)
)

// 任务位置
type Task struct {
	Id   int    `json:"id"`
	Path string `json:"path"`
}

// 用于解析 json 文件
type tasksJson struct {
	Tasks []Task `json:"tasks"`
}

// 任务统计的 json 文件
type tasksNumJson struct {
	Count []tasksCount `json:"count"`
}

// 任务数据（可用于解析 json 文件）
type TaskData struct {
	Name    string `json:"name"`
	Info    string `json:"info"`
	Diy     bool   `json:"diy"`
	Run     bool   `json:"run"`
	Success bool   `json:"success"`
	Cycle   string `json:"cycle"`
	Command string `json:"command"`
	File    string `json:"file"`
	Lastime int64  `json:"lastime"`
}

// 任务统计
type tasksCount struct {
	Date    string `json:"date"`
	Total   int    `json:"total"`
	Success int    `json:"success"`
	Fail    int    `json:"fail"`
}

// LoadTasks 接收任务信息
func LoadTasks(tasksInfo []byte) {
	TasksInfo = tasksInfo
	defer addTasks()
}

// addTasks 解析任务，添加任务
func addTasks() {
	json.Unmarshal(TasksInfo, &tasksData)
	Tasks = append(Tasks, tasksData.Tasks...)
	defer PlanningTasks()
}

// PlanningTasks 计划任务
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
		AddCron(task.Id)
	}
}

// AddCron 添加定时器
func AddCron(id int) {
	TaskInfo := ReadTaskInfo(id)
	cron := cron.New()
	cron.AddFunc(TaskInfo.Cycle, func() {
		RunTask(id)
	})
	// 只有可执行或上次执行成功的任务才开启定时器
	if TaskInfo.Run && TaskInfo.Success {
		cron.Start()
	}
	Crons = append(Crons, *cron)
}

func DeleteCron(id int) {
	Crons[id].Stop()
	Crons = append(Crons[:id], Crons[id+1:]...) // 删除定时器
}

// init 系统定时器
func init() {
	cron := cron.New()
	cron.AddFunc("0 0 23 * * ?", func() { // 每天 23:00 执行，记录任务执行情况
		var (
			success   int
			failed    int
			JsonFile  tasksNumJson
			countData tasksNumJson
		)

		file, _ := os.OpenFile("./data/tasks/count.json", os.O_WRONLY|os.O_CREATE, 0644)
		defer file.Close()

		countInfo, _ := os.ReadFile("./data/tasks/count.json")
		json.Unmarshal(countInfo, &countData)

		for _, task := range Tasks {
			taskInfo := ReadTaskInfo(task.Id)
			if taskInfo.Success {
				success++
			} else {
				failed++
			}
		}

		JsonFile.Count = make([]tasksCount, 0, 1)
		_, month, day := time.Now().Date()

		if len(countData.Count) > 7 {
			JsonFile.Count = append(JsonFile.Count, countData.Count[len(countData.Count)-7:]...)
		} else {
			JsonFile.Count = append(JsonFile.Count, countData.Count...)
		}

		JsonFile.Count = append(JsonFile.Count, tasksCount{
			Date:    strconv.Itoa(int(month)) + "/" + strconv.Itoa(day),
			Total:   len(Tasks),
			Success: success,
			Fail:    failed,
		})

		jsonData, _ := json.Marshal(JsonFile)

		file.Truncate(0)
		_, err := io.WriteString(file, string(jsonData))
		if err != nil {
			fmt.Println(err) // ---日志
		}
	})
	cron.Start()
}
