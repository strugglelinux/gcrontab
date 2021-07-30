package tool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Cron struct {
	initFilePath string  //配置文件路径
	taskList     []*Task //任务列表
	//waiter       sync.WaitGroup
	mux      sync.Mutex
	parser   *Parse
	location time.Time
	//taskChan     chan tool.Task
	isRuning bool
	log      *Logger
}

var filepath string = "gcrontab.conf"

func NewCron() (*Cron, error) {
	cron := &Cron{}
	cron.initFilePath = "./" + filepath
	pwd, _ := os.Getwd()
	dir := pwd + "/log"
	pos := strings.LastIndex(pwd, "/")
	if pos != -1 {
		lastDir := pwd[pos+1:]
		if lastDir == "tool" {
			dir = pwd[:pos] + "/log"
			cron.initFilePath = pwd[:pos] + "/" + filepath
		}
	}
	cron.log = NewLogger(dir)
	cron.Init()
	return cron, nil
}

//添加行计划和函数
func (c *Cron) AddFunc(plan string, funcJob func()) {
	count := len(c.taskList)
	t, err := c.parser.parsePlan(plan)
	if err != nil {
		panic(err.Error())
	}
	t.Id = count + 1
	t.cmd.funcJob = funcJob
	c.taskList = append(c.taskList, t)
}

//初始化信息
func (c *Cron) Init() {
	tcount := len(c.taskList)
	if ok, _ := c.parser.FileExists(filepath); !ok {
		f, err := os.Create(filepath)
		if err != nil {
			c.log.Error(err.Error())
			panic(err.Error())
		}
		c.initFilePath = filepath
		defer f.Close()
	}
	_parser, err := NewParse(c.initFilePath)
	if err != nil {
		c.log.Error(err.Error())
		panic(err.Error())
	}
	c.parser = _parser
	c.parser.SetLogger(c.log)
	taskList, err := c.parser.Load()
	if err != nil {
		c.log.Error(err.Error())
		if tcount == 0 {
			panic(err.Error())
		}
	}
	c.location = time.Now()
	c.taskList = taskList
	c.mux = sync.Mutex{}
	c.isRuning = false
}

//开始
func (c *Cron) Start() {
	c.log.Info("Cron Start ")
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.isRuning {
		return
	}
	c.isRuning = true
	c.log.Init()
	c.run()
}

//执行
func (c *Cron) run() {
	c.log.Info("Cron run ")
	ticker := time.NewTicker(time.Millisecond * 500)
	for now := range ticker.C {
		for _, task := range c.taskList {
			t := task
			format := "2006-01-02 15:04:05"
			if now.Local().Format(format) != t.preTime.Local().Format(format) {
				c.startTask(now, t)
			}
		}
	}
}

//开始任务
func (c *Cron) startTask(d time.Time, t *Task) {
	go func() {
		t.Start(d)
	}()
}

//序列化
func (c *Cron) String() string {
	b, err := json.Marshal(&c)
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
