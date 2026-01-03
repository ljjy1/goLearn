package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

/**
 * @description: 用户模型
 */
type User struct {
	ID        uint   `gorm:"primarykey"`
	Username  string `json:"username" gorm:"not null;size:20;comment:登录账号;uniqueIndex:idx_username_deleted_at"` //设置不能为null 长度20 唯一索引要和DeletedAt关联
	Nickname  string `json:"nickname" gorm:"not null;size:64;comment:昵称"`                                       //设置不能为null 长度64
	Password  string `json:"password" gorm:"not null;size:64;comment:登录密码(加密后的)"`                               //设置不能为null 长度64
	Email     string `json:"email" gorm:"size:128;comment:邮箱"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt soft_delete.DeletedAt `gorm:"uniqueIndex:idx_username_deleted_at"`

	//关联文章模型 一对多关系 外键为UserID 引用为ID
	Posts []Post `json:"posts" gorm:"foreignKey:UserID;references:ID;comment:文章"`
}

// 配置表中文注释
func (u *User) TableComment() string {
	return "用户表"
}
