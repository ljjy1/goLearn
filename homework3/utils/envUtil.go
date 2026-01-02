package utils

import (
	"os"

	"github.com/joho/godotenv"
)

/**
 * @description: 获取环境变量 如果为空 则返回默认值
 * @param {string} key 环境变量键名
 * @param {string} defaultValue 默认值
 * @return {string} 环境变量值或默认值
 */
func LoadEnvFile() {
	err := godotenv.Load()
	if err != nil {
		panic("错误: 无法加载 .env 文件, 请检查文件是否存在: " + err.Error())
	}
}

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

/**
 * @description: 设置环境变量 如果值为空 则不设置
 * @param {string} key 环境变量键名
 * @param {string} value 环境变量值
 */
func SetEnv(key, value string) {
	if value == "" {
		return
	}
	os.Setenv(key, value)
}
