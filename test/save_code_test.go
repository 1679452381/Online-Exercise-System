package test

import (
	"online_exercise_system/utils"
	"os"
	"testing"
)

// 保存代码 返回文件路径
func SaveCode(code []byte) (string, error) {
	fileName := "code/" + utils.GetUUID()
	path := fileName + "/main.go"
	err := os.MkdirAll(fileName, 0777)
	if err != nil {
		return "", err
	}
	//创建并打开文件
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 777)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(code)
	if err != nil {
		return "", err
	}
	return path, nil
}

func TestSaveCode(t *testing.T) {
	SaveCode([]byte("zxczxc"))
}
