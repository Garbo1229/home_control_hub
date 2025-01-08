package utils

import (
	"fmt"
	"os"
)

func GetPwd() string {
	// 获取当前工作目录
	root, err := os.Getwd()
	if err != nil {
		// 将错误信息转换为字符串并传递给 handleError 函数
		fmt.Println("[helper]处理root径失败：", err)
	}
	return root
}
