package common

import (
	"os"
	"path/filepath"
)

//GetLastAbsPathByExecutable 获取执行程序所在的上一层级目录的绝对路径
func GetLastAbsPathByExecutable() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	res, _ := filepath.EvalSymlinks(exePath)
	res = filepath.Dir(res) //获取程序所在目录
	res = filepath.Dir(res) //获取程序所在上一层级目录
	return res, nil
}
