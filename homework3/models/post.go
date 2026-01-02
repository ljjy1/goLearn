package models

import (
	"encoding/json"
	"fmt"

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
	CommentStatus uint8  `json:"commentStatus" gorm:"default:0;comment:评论状态"`           //评论状态 0-无评论(默认) 1-有评论

	//关联评论模型 一对多关系 外键为PostID 引用为ID
	Comments []Comment `json:"comments" gorm:"foreignKey:PostID;references:ID;comment:评论"`
}

// 配置表中文注释
func (p *Post) TableComment() string {
	return "文章表"
}

// 配置String方法 方便打印
func (p Post) ToString() string {
	jsondata, _ := json.MarshalIndent(p, "", "  ")
	return string(jsondata)
}

// 给文章添加一个后置创建钩子
func (p *Post) AfterCreate(db *gorm.DB) (err error) {
	//更新用户文章数量
	updateErr := db.Model(&User{}).Where("id = ?", p.UserID).
		Update("post_count", gorm.Expr("post_count + ?", 1)).Error
	if updateErr != nil {
		fmt.Printf("更新用户%d文章数量失败:%v\n", p.UserID, updateErr)
		return
	}
	return
}
