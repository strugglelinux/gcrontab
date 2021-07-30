package tool

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const notesChar1 = "#"  //注释符1  //。。。。
const notesChar2 = "//" //注释符2  //。。。。

type Parse struct {
	filePath string
	taskList []*Task
	log      *Logger
}

func NewParse(file string) (*Parse, error) {
	if len(file) == 0 {
		return nil, fmt.Errorf("配置文件路径为空")
	}
	_parse := &Parse{file, []*Task{}, &Logger{}}
	return _parse, nil
}

//设置日志
func (p *Parse) SetLogger(l *Logger) {
	p.log = l
}

/**
 *判断文件是否存在
 **/
func (p *Parse) FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/**
 * 读取文件内容
 */
func (p *Parse) readAll(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	var list []string
	buf := bufio.NewReader(f)
	for {
		line := ""
		b, err := buf.ReadBytes('\n')
		if err != nil || err == io.EOF {
			break
		}
		line = string(b)
		post1 := strings.Index(line, notesChar1)
		post2 := strings.Index(line, notesChar2)

		if post1 == 0 || post2 == 0 {
			continue
		}
		if post1 > 0 || post2 > 0 {
			if post1 < post2 && post1 != -1 {
				line = line[:post1]
			}
			if post1 > post2 && post2 != -1 {
				line = line[:post2]
			}
		}
		list = append(list, line)
	}
	return list, err
}

/**
 *载入文件解析文件
 */
func (p *Parse) Load() ([]*Task, error) {
	p.log.Info("解析文件 " + p.filePath)
	if ok, err := p.FileExists(p.filePath); !ok {
		return nil, err
	}
	list, err := p.readAll(p.filePath)
	if err != nil && err != io.EOF {
		return nil, err
	}
	for k, _plan := range list {
		options := strings.Fields(_plan)
		_plan_opts := options[:6]
		_cmd := options[6:]
		if len(_cmd) == 0 {
			err = fmt.Errorf("执行计划格式错误 ：%s", _plan)
			break
		}
		t, err := p.parsePlan(strings.Join(_plan_opts, " "))
		if err != nil {
			break
		}
		t.Id = k + 1
		t.cmd.Cmd = strings.Join(_cmd, " ")
		p.taskList = append(p.taskList, t)
	}
	return p.taskList, err
}

//解析计划
func (p *Parse) parsePlan(plan string) (*Task, error) {
	_plan_opts := strings.Fields(plan)
	task := &Task{cmd: Command{}}
	task.cmd.log = p.log
	task.cmd.Plan = strings.Join(_plan_opts, " ")

	sec_opts, err := p.optionParse(_plan_opts[0], 0, 59)
	if err != nil {
		return task, fmt.Errorf("秒 计划项解析错误")
	}
	task.cmd.SecondMap = sec_opts
	min_opts, err := p.optionParse(_plan_opts[1], 0, 59)
	if err != nil {
		return task, fmt.Errorf("分 计划项解析错误")
	}
	task.cmd.MinuteMap = min_opts
	hour_opts, err := p.optionParse(_plan_opts[2], 0, 59)
	if err != nil {
		return task, fmt.Errorf("时 计划项解析错误")
	}
	task.cmd.HourMap = hour_opts
	day_opts, err := p.optionParse(_plan_opts[3], 1, 31)
	if err != nil {
		return task, fmt.Errorf("日 计划项解析错误")
	}
	task.cmd.DayMap = day_opts
	month_opts, err := p.optionParse(_plan_opts[4], 1, 12)
	if err != nil {
		return task, fmt.Errorf("月 计划项解析错误")
	}
	task.cmd.MonthMap = month_opts
	week_opts, err := p.optionParse(_plan_opts[5], 0, 6)
	if err != nil {
		return task, fmt.Errorf("周 计划项解析错误")
	}
	task.cmd.WeekMap = week_opts
	return task, nil
}

/**
 * 解析计划项
 */
func (p *Parse) optionParse(option string, min, max int) (map[int]int, error) {
	option = strings.Trim(option, " ")
	length := len(option)
	if length == 0 {
		return nil, fmt.Errorf("计划项为空")
	}
	var optionslp = make(map[int]int)
	options := strings.Split(option, ",")
	for _, _option := range options {
		_min := min
		_max := max
		step := 1
		_option = strings.Trim(_option, " ")
		pos1 := strings.Index(_option, "/")
		if pos1 != -1 {
			opts := strings.Split(_option, "/")
			_step, _ := strconv.Atoi(opts[1])
			step = _step
		}
		pos2 := strings.Index(_option, "-")
		if pos2 != -1 {
			opts := strings.Split(_option, "-")
			_a, _ := strconv.Atoi(opts[0])
			_min = _a
			_b, _ := strconv.Atoi(opts[1])
			_max = _b
		}
		if pos1 == -1 && pos2 == -1 && _option != "*" {
			_opt, err := strconv.Atoi(_option)
			if err == nil {
				optionslp[_opt] = 1
				continue
			}
		}

		for i := _min; i <= _max; i += step {
			if _, ok := optionslp[i]; !ok {
				optionslp[i] = 1
			}
		}
	}
	if len(optionslp) == 0 {
		return nil, fmt.Errorf("计划项解析失败")
	}
	return optionslp, nil
}
