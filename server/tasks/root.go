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

type tasksNumJson struct {
	Count []tasksCount `json:"count"`
}

type TaskData struct {
	Name    string `json:"name"`
	Info    string `json:"info"`
	Diy     bool   `json:"diy"`
	Cycle   string `json:"cycle"`
	Command string `json:"command"`
	File    string `json:"file"`
	Lastime int64  `json:"lastime"`
}

type tasksCount struct {
	Date    string `json:"date"`
	Total   int    `json:"total"`
	Success int    `json:"success"`
	Fail    int    `json:"fail"`
}

func LoadTasks(tasksInfo []byte) {
	TasksInfo = tasksInfo
	defer addTasks()
}

func addTasks() {
	json.Unmarshal(TasksInfo, &tasksData)
	Tasks = append(Tasks, tasksData.Tasks...)
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
		AddCron(task.Id)
	}
}

func AddCron(id int) {
	// 如任务上次未执行成功则跳过加载定时器
	if !Tasks[id].Success {
		return
	}

	TaskInfo := ReadTaskInfo(id)
	cron := cron.New()

	cron.AddFunc(TaskInfo.Cycle, func() {
		RunTask(id)
	})
	cron.Start()
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
			if task.Success {
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
