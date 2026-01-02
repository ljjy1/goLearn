package models

import (
	"encoding/json"
	"time"

	"gorm.io/plugin/soft_delete"
)

/**
 * @description: 用户模型
 */
type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Nickname  string                `json:"nickname" gorm:"not null;size:20;comment:昵称"`                                  //设置不能为null 长度20 唯一索引
	Mobile    string                `json:"mobile" gorm:"not null;size:11;comment:手机号;uniqueIndex:idx_mobile_deleted_at"` //设置不能为null 长度11 唯一索引要和DeletedAt关联
	Age       *uint8                `json:"age" gorm:"comment:年龄"`                                                        //年龄可以为空 使用指针避免查询返回0
	PostCount uint32                `json:"postCount" gorm:"default:0;comment:文章数量"`                                      //默认值为0
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_mobile_deleted_at"`

	//关联文章模型 一对多关系 外键为UserID 引用为ID
	Posts []Post `json:"posts" gorm:"foreignKey:UserID;references:ID;comment:文章"`
}

// 配置表中文注释
func (u *User) TableComment() string {
	return "用户表"
}

// 配置String方法 方便打印
func (u User) ToString() string {
	jsondata, _ := json.MarshalIndent(u, "", "  ")
	return string(jsondata)
}
