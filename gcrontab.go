package main

import (
	"fmt"
	"tool"
)

func main() {
	cron, err := tool.NewCron()
	if err != nil {
		panic(err.Error())
	}
	cron.AddFunc("0 1-5,10-20,30-40 * * * *", func() {
		fmt.Println("hello 0 1-5,10-20,30-40 15 * * *")
	})
	cron.AddFunc("0 20,30,32,36 14 * * *", func() {
		fmt.Println("hello 0 20,30,32,36 * * * *")
	})
	cron.Start()
}
