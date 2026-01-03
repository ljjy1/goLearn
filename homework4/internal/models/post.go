package models

import (
	"gorm.io/gorm"
)

/**
 * @description: 帖子模型
 */
type Post struct {
	gorm.Model
	UserID        uint   `json:"userId" gorm:"not null;index:idx_user_id;comment:用户ID"` //设置不能为null 引用用户模型ID 添加一个索引
	Title         string `json:"title" gorm:"not null;size:20;comment:标题"`              //设置不能为null 长度20 唯一索引
	Content       string `json:"content" gorm:"not null;size:200;comment:内容"`           //设置不能为null 长度200
	
	//关联评论模型 一对多关系 外键为PostID 引用为ID
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID;references:ID;comment:评论"`
}

// 配置表中文注释
func (p *Post) TableComment() string {
	return "文章表"
}