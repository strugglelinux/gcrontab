package tool

import (
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
	"time"
)

func TestCron_run(t *testing.T) {
	dir := "/private/var/gospace/src/gcron/log"
	Log := NewLogger(dir)
	Log.SetVerbose(true)
	file := "/private/var/gospace/src/gcron/gcrontab.conf"
	parse, err := NewParse(file)
	if err != nil {
		t.Errorf("%T\n", parse)
	}
	parse.SetLogger(Log)
	taskList, _ := parse.Load()

	if Log == nil {
		t.Errorf("创建日志错误")
	}
	type fields struct {
		initFilePath string
		taskList     []*Task
		mux          sync.Mutex
		parser       *Parse
		location     time.Time
		isRuning     bool
		log          *Logger
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Run",
			fields: fields{
				initFilePath: file,
				taskList:     taskList,
				parser:       parse,
				location:     time.Time{},
				isRuning:     false,
				log:          Log,
				mux:          sync.Mutex{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if p := recover(); p != nil {
					fmt.Printf("panic recover! p: %v", p)
					debug.PrintStack()
				}
			}()
			c := &Cron{
				initFilePath: tt.fields.initFilePath,
				taskList:     tt.fields.taskList,
				parser:       tt.fields.parser,
				location:     tt.fields.location,
				isRuning:     tt.fields.isRuning,
				log:          tt.fields.log,
				mux:          tt.fields.mux,
			}
			c.Start()
		})
	}
}
