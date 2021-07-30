package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Command struct {
	SecondMap map[int]int //秒(0 - 59)
	MinuteMap map[int]int //分钟 (0 - 59)
	HourMap   map[int]int //小时 (0 - 23)
	DayMap    map[int]int //一个月中的第几天 (1 - 31)
	MonthMap  map[int]int // 月份 (1 - 12)
	WeekMap   map[int]int //星期中星期几 (0 - 6) (星期天 为0)
	Plan      string      //执行计划
	Cmd       string
	mux       sync.Mutex
	isRuning  bool
	log       *Logger
	funcJob   func() //命令执行方法
}

//执行命令
func (c *Command) Run() {
	c.mux.Lock()
	c.isRuning = true
	defer func() {
		c.mux.Unlock()
		c.isRuning = false
	}()
	if len(c.Cmd) == 0 {
		c.runFunc()
	} else {
		c.runPlan()
	}

}

//方法执行
func (c *Command) runFunc() {
	c.funcJob()
}

//执行计划
func (c *Command) runPlan() {
	cmds := strings.Fields(c.Cmd)
	exc := exec.Command(cmds[0], cmds[1:]...)
	if err := exc.Start(); err != nil {
		c.log.Error("Execute failed when Start :" + err.Error())
		return
	}
	if err := exc.Wait(); err != nil {
		c.log.Error("Execute failed when Wait :" + err.Error())
		return
	}
	c.log.Info(c.Cmd + " : Execute finished ")
}

//序列化对象
func (c *Command) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		return fmt.Sprintf("%+v", &c)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", &c)
	}
	return out.String()
}

type Task struct {
	Id      int
	preTime time.Time //上次执行时间
	cmd     Command   //执行命令
	mux     sync.Mutex
}

//任务开始执行
func (t *Task) Start(dateTime time.Time) {
	t.mux.Lock()
	defer t.mux.Unlock()
	ok := t.timeVerify(dateTime)
	if ok && !t.cmd.isRuning {
		t.cmd.log.Info("Task id " + strconv.Itoa(t.Id) + " 开始执行")
		t.preTime = dateTime
		t.cmd.Run()
	}
}

//时间验证
func (t *Task) timeVerify(dateTime time.Time) bool {
	month := dateTime.Local().Month()
	_, ok := t.cmd.MonthMap[int(month)]
	if !ok {
		return false
	}
	week := dateTime.Local().Weekday()
	_, ok = t.cmd.WeekMap[int(week)]
	if !ok {
		return false
	}
	day := dateTime.Local().Day()
	_, ok = t.cmd.DayMap[int(day)]
	if !ok {
		return false
	}
	hour := dateTime.Local().Hour()
	_, ok = t.cmd.HourMap[int(hour)]
	if !ok {
		return false
	}
	minute := dateTime.Local().Minute()
	_, ok = t.cmd.MinuteMap[int(minute)]
	if !ok {
		return false
	}
	second := dateTime.Local().Second()
	_, ok = t.cmd.SecondMap[int(second)]
	if !ok {
		return false
	}
	return ok
}

//序列化对象
func (t *Task) String() string {
	b, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%+v", &t)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", &t)
	}
	return out.String()
}

//获取上一次执行时间
func (t *Task) GetLastRunTime() time.Time {
	return t.preTime
}
