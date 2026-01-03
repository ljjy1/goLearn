package models

import (
	"gorm.io/gorm"
)

/**
 * @description: 评论模型
 */
type Comment struct {
	gorm.Model
	UserID  uint   `json:"userId" gorm:"not null;index:idx_user_id;comment:用户ID"` //设置不能为null 引用用户模型ID 添加一个索引
	PostID uint `json:"postId" gorm:"not null;index:idx_post_id;comment:文章ID"` //设置不能为null 引用文章模型ID 添加一个索引
	Content string `json:"content" gorm:"not null;size:200;comment:内容"`  //设置不能为null 长度200

	//用户信息
	User User `json:"user" gorm:"foreignKey:UserID;references:ID;comment:用户"`
}

//配置表中文注释
func (c *Comment) TableComment() string {
	return "评论表"
}