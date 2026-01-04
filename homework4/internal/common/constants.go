package common

import (
	"path/filepath"
	"runtime"
)

/**
 * @Description: 定义常量
 */

// GetLogFile 获取日志文件的绝对路径
func GetLogFile() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../../logs/app-run.log")
}
