package models

import (
	"gorm.io/gorm"
)

type InitDataConfig struct {
	gorm.Model
	InitDefaultData bool `json:"initDefaultData" gorm:"default:false;comment:是否初始化默认数据 0-未初始化 1-已初始化"`
}
