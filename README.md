# Gcrontab

## 可执行计划任务

### 使用方式

```go
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
  fmt.Println("0 20,30,32,36 14 * * *")
 })
 cron.Start()
}
```

### 配置文件 gcrontab.conf 设置

```shell
*/5 * * * * cd /project/ && php think a:b
```

### 计划参数

#### 格式  `* * * * * *`

| 设置项   | 值范围 | 字符范围 |
| -------- | :----: | :------: |
| 秒       |  0-59  | * / , -  |
| 分       |  0-59  | * / , -  |
| 时       |  0-59  |  * / ,-  |
| 天(某日) |  1-31  |  * / -   |
| 月       |  1-12  | * / , -  |
| 周       |  0-6   | * / , -  |
