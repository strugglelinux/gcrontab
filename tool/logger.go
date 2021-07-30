package tool

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	ERROR       = "ERROR"
	WARING      = "WARNING"
	INFO        = "INFO"
	dirOpenMode = 0777

	// 文件写入mode
	fileOpenMode = 0666
	// 文件Flag
	fileFlag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	// 换行符
	newlineFlag = "\n"
	//缓存切片容量
	cacheCap = 64
)

type Cache struct {
	isUse bool //是否使用缓存
	data  []string
	mux   sync.Mutex
	rwMux sync.RWMutex
}

type Logger struct {
	verbose  bool
	mux      sync.Mutex
	dir      string
	file     *os.File
	filePath string
	cache    Cache
}

func NewLogger(dir string) *Logger {
	if is := isDir(dir); !is {
		if ok := createDir(dir); !ok {
			panic(dir + " 目录不存在")
		}
	}
	lg := &Logger{verbose: false, dir: dir}
	lg.cache.isUse = true
	lg.dir = dir
	return lg
}

func (l *Logger) Init() {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for {
			select {
			case <-ticker.C:
				l.cache.rwMux.RLock()
				go l.flush()
				l.cache.rwMux.RUnlock()
			}

		}
	}()
}

//设置是否直接打印输出
func (l *Logger) SetVerbose(isVerbose bool) {
	l.verbose = isVerbose

}

//设置是否使用缓存
func (l *Logger) SetUserCache(isUse bool) {
	l.cache.isUse = isUse

}

//创建目录
func createDir(dir string) bool {
	err := os.Mkdir(dir, dirOpenMode)
	if err != nil {
		panic(err.Error())
	}
	os.Chmod(dir, dirOpenMode) //通过chmod重新赋权限
	return true
}

//判断目录
func isDir(dir string) bool {
	s, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//错误信息
func (l *Logger) Error(msg string) {
	date := l.nowDate()
	l.log(date, ERROR, msg)
}

//警告信息
func (l *Logger) Warning(msg string) {
	date := l.nowDate()
	l.log(date, WARING, msg)
}

//信息
func (l *Logger) Info(msg string) {
	date := l.nowDate()
	l.log(date, INFO, msg)
}

//处理方式
func (l *Logger) log(date, option, msg string) {
	infoString := fmt.Sprintf("[%s] %s:%s", date, option, msg)
	if l.verbose {
		fmt.Println(infoString)
	} else {
		l.save(infoString)
	}
}

func (l *Logger) nowDate() string {
	now := time.Now()
	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	zone, _ := now.Zone()
	return fmt.Sprintf("%d-%d-%d %02d:%02d:%02d %s", year, mon, day, hour, min, sec, zone)
}

//判断文件是否存在
func (l *Logger) fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//获取操作的日志文件
func (l *Logger) getFile() string {
	now := time.Now()
	year, mon, day := now.Date()
	hour, _, _ := now.Clock()
	dateDir := fmt.Sprintf("%d%d%d", year, mon, day)
	dateFile := fmt.Sprintf("%d%d%d%d", year, mon, day, hour)
	fileDir := l.dir + "/" + dateDir
	if !isDir(fileDir) {
		if ok := createDir(fileDir); !ok {
			panic(fileDir + " 目录创建失败")
		}
		l.Info("日志目录：" + fileDir)
	}
	filePath := fileDir + "/" + dateFile + ".log"
	if l.filePath == filePath {
		return l.filePath
	}
	ok, _ := l.fileExists(filePath)
	if !ok {
		f, err := os.Create(filePath)
		if err != nil {
			f, err = os.Create(filePath)
			if err != nil {
				panic(filePath + " Create Fail " + err.Error())
			}
		}
		l.Info("日志文件：" + filePath)
		defer f.Close()
	}
	return filePath
}

//打开日志文件
func (l *Logger) openFile() (*os.File, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	filePath := l.getFile()
	file, err := os.OpenFile(filePath, fileFlag, fileOpenMode)
	if err != nil {
		//重试
		file, err = os.OpenFile(filePath, fileFlag, fileOpenMode)
		if err != nil {
			return file, nil
		}
	}
	if l.file != nil {
		l.file.Close()
	}
	l.file = file
	l.filePath = filePath
	return file, nil
}

func (l *Logger) openFileCache() (*os.File, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	filePath := l.getFile()
	file, err := os.OpenFile(filePath, fileFlag, fileOpenMode)
	if err != nil {
		//重试
		file, err = os.OpenFile(filePath, fileFlag, fileOpenMode)
		if err != nil {
			return file, nil
		}
	}
	return file, nil
}

//写入信息
func (l *Logger) writer(msg []byte) error {
	file, err := l.openFile()
	if err != nil {
		panic(err)
	}
	l.mux.Lock()
	_, err = file.Write(msg)
	l.mux.Unlock()
	return err
}

//写入信息
func (l *Logger) save(msg string) {
	if l.cache.isUse {
		l.cachePush(msg) //写入缓存
	} else {
		msg = msg + newlineFlag
		l.writer([]byte(msg))
	}
}

//追加入缓存
func (l *Logger) cachePush(msg string) {
	l.cache.mux.Lock()
	l.cache.data = append(l.cache.data, msg+newlineFlag)
	l.cache.mux.Unlock()

}

//同步缓存到文件
func (l *Logger) flush() error {
	file, err := l.openFileCache()
	if err != nil {
		panic(err)
	}
	defer file.Close()
	l.cache.mux.Lock()
	cacheData := l.cache.data
	l.cache.data = make([]string, 0, cacheCap)
	l.cache.mux.Unlock()
	if len(cacheData) == 0 {
		return nil
	}
	_, err = file.WriteString(strings.Join(cacheData, ""))
	if err != nil {
		_, err = file.WriteString(strings.Join(cacheData, ""))
		if err != nil {
			panic(err.Error())
		}
	}
	return nil
}

//缓存日志刷新到文件
func (l *Logger) Flush() {
	err := l.flush()
	if err != nil {
		panic(err)
	}
}
