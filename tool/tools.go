package tool

import "os"

type FileToolBox struct {
}

//创建目录
func CreateDir(dir string, mode os.FileMode) bool {
	err := os.Mkdir(dir, mode)
	if err != nil {
		panic(err.Error())
	}
	return true
}

//判断目录
func IsDir(dir string) bool {
	s, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//判断文件是否存在
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
