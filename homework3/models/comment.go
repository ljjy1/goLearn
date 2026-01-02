package models

import (
	"gorm.io/gorm"
	"encoding/json"
	"fmt"
)

/**
 * @description: 评论模型
 */
type Comment struct {
	gorm.Model
	PostID uint `json:"postId" gorm:"not null;index:idx_post_id;comment:文章ID"` //设置不能为null 引用文章模型ID 添加一个索引
	Content string `json:"content" gorm:"not null;size:200;comment:内容"`  //设置不能为null 长度200
}

//配置表中文注释
func (c *Comment) TableComment() string {
	return "评论表"
}

//配置String方法 方便打印
func (c Comment) ToString() string {
	jsondata, _ := json.MarshalIndent(c, "", "  ")
	return string(jsondata)
}

//添加一个后置删除钩子
func (c *Comment) AfterDelete(db *gorm.DB) (err error) {
	//查询当前评论所属的文章的评论数需要排除当前ID
	var count *int64 = new(int64)

	queryErr := db.Model(&Comment{}).Where("post_id = ? and id != ?", c.PostID, c.ID).
		Count(count).Error
	if queryErr != nil {
		fmt.Printf("查询文章%d评论数量失败:%v\n", c.PostID, err)
		return
	}
	if *count == 0 {
		//如果评论数为0 则更新文章评论状态为0
		updateErr := db.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_status", 0).Error
		if updateErr != nil {
			fmt.Printf("更新文章%d评论状态失败:%v\n", c.PostID, updateErr)
			return
		}
	}
	return
}
