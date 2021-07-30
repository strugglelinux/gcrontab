package tool

import (
	"os"
	"strings"
	"sync"
	"testing"
)

func TestCommand_String(t *testing.T) {
	file := "/var/gospace/src/gcron/gcrontab.conf"
	pwd, _ := os.Getwd()
	pos := strings.LastIndex(pwd, "/")
	if pos != -1 {
		lastDir := pwd[pos+1:]
		if lastDir == "tool" {
			file = pwd[:pos] + "/gcrontab.conf"
		}
	}
	parse, err := NewParse(file)
	if err != nil {
		t.Errorf("%T\n", parse)
	}
	task, _ := parse.parsePlan("0,1,3 0,2,3 0,2,3 0,2,3 0,2,3 0,2,3")
	want := task.cmd.String()
	type fields struct {
		SecondMap map[int]int
		MinuteMap map[int]int
		HourMap   map[int]int
		DayMap    map[int]int
		MonthMap  map[int]int
		WeekMap   map[int]int
		Plan      string
		Cmd       string
		mux       sync.Mutex
		isRuning  bool
		log       *Logger
		funcJob   func()
	}

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Start",
			fields: fields{
				SecondMap: map[int]int{0: 1, 1: 1, 3: 1},
				MinuteMap: map[int]int{0: 1, 2: 1, 3: 1},
				HourMap:   map[int]int{0: 1, 2: 1, 3: 1},
				DayMap:    map[int]int{0: 1, 2: 1, 3: 1},
				MonthMap:  map[int]int{0: 1, 2: 1, 3: 1},
				WeekMap:   map[int]int{0: 1, 2: 1, 3: 1},
				Plan:      "0,1,3 0,2,3 0,2,3 0,2,3 0,2,3 0,2,3",
			},
			want: want,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				SecondMap: tt.fields.SecondMap,
				MinuteMap: tt.fields.MinuteMap,
				HourMap:   tt.fields.HourMap,
				DayMap:    tt.fields.DayMap,
				MonthMap:  tt.fields.MonthMap,
				WeekMap:   tt.fields.WeekMap,
				Plan:      tt.fields.Plan,
				Cmd:       tt.fields.Cmd,
				mux:       tt.fields.mux,
				isRuning:  tt.fields.isRuning,
				log:       tt.fields.log,
				funcJob:   tt.fields.funcJob,
			}
			if got := c.String(); got != tt.want {
				t.Errorf("Command.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
