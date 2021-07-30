package tool

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	file := "/var/gospace/src/gcron/gcrontab.conf"
	pwd, _ := os.Getwd()
	dir := pwd + "/log"
	pos := strings.LastIndex(pwd, "/")
	if pos != -1 {
		lastDir := pwd[pos+1:]
		if lastDir == "tool" {
			dir = pwd[:pos] + "/log"
			file = pwd[:pos] + "/gcrontab.conf"
		}
	}
	log := NewLogger(dir)
	parse, err := NewParse(file)
	if err != nil {
		t.Errorf("%T\n", parse)
	}
	parse.SetLogger(log)
	taskList, err := parse.Load()
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Printf("%v\n", taskList)

}
