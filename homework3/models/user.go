package models

import (
	"gorm.io/gorm"
	"encoding/json"
)

/**
 * @description: 用户模型
 */
type User struct {
	gorm.Model
	Nickname string `json:"nickname" gorm:"not null;size:20;comment:昵称"`  //设置不能为null 长度20 唯一索引
	Mobile string `json:"mobile" gorm:"not null;size:11;uniqueIndex:uq_mobile;comment:手机号"`  //设置不能为null 长度11 唯一索引
	Age *uint8 `json:"age" gorm:"comment:年龄"`  //年龄可以为空 使用指针避免查询返回0
	PostCount uint32 `json:"postCount" gorm:"default:0;comment:文章数量"` //默认值为0

	//关联文章模型 一对多关系 外键为UserID 引用为ID
	Posts []Post `json:"posts" gorm:"foreignKey:UserID;references:ID;comment:文章"`
}


//配置表中文注释
func (u *User) TableComment() string {
	return "用户表"
}

//配置String方法 方便打印
func (u User) ToString() string {
	jsondata, _ := json.MarshalIndent(u, "", "  ")
	return string(jsondata)
}






