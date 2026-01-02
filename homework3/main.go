package main

import (
	"fmt"
	"homework3/config"
	"homework3/models"

	"gorm.io/gorm"
)

func main() {
	//初始化数据库连接
	config.InitDB()

	//获取数据库连接实例
	db := config.GetDB()

	//查询
	// TestSelect(db)

	//新增文章数据
	// TestCreateHook(db)

	//删除评论
	TestDeleteHook(db)

}

func TestSelect(db *gorm.DB) {
	var userById models.User
	err := db.Where("id = ?", 1).
		Preload("Posts.Comments").
		Find(&userById).Error
	if err != nil {
		fmt.Println("查询用户失败:", err)
	} else {
		fmt.Println("查询到ID等于1的用户:", userById.ToString())
	}

	//查询评论最多的的文章
	var maxCommentPost models.Post
	err = db.Model(&models.Post{}).
		Select("table_post.*, COUNT(table_comment.id) as comment_count").
		Joins("LEFT JOIN table_comment AS table_comment ON table_post.id = table_comment.post_id").
		Group("table_post.id").
		Order("comment_count DESC").
		First(&maxCommentPost).Error
	if err != nil {
		fmt.Println("查询评论最多的文章失败:", err)
	} else {
		if maxCommentPost.ID > 0 {
			//根据文章ID查询所有评论
			var comments []models.Comment
			db.Where(&models.Comment{PostID: maxCommentPost.ID}).Find(&comments)
			if len(comments) > 0 {
				maxCommentPost.Comments = comments
			}
		}
		fmt.Println("查询到评论最多的文章:", maxCommentPost.ToString())
	}

}

func TestCreateHook(db *gorm.DB) {
	//创建文章 设置默认用户ID为1
	post := models.Post{
		Title:         "测试用户1文章3",
		Content:       "测试用户1文章3内容",
		CommentStatus: 1,
		UserID:        1,

		Comments: []models.Comment{
			{
				Content: "测试用户1文章3评论1",
			},
			{
				Content: "测试用户1文章3评论2",
			},
		},
	}
	err := db.Create(&post).Error
	if err != nil {
		fmt.Println("创建文章失败:", err)
		return
	}
	fmt.Println("创建文章:", post.ToString())

	type PostCount struct {
		PostCount uint32 `json:"postCount"`
	}
	//查询用户文章数量
	var postCount PostCount
	err = db.Model(&models.User{}).
		Where("id = ?", 1).
		Scan(&postCount).Error
	if err != nil {
		fmt.Println("查询用户文章数量失败:", err)
		return
	}
	fmt.Printf("用户1文章数量:%d\n", postCount.PostCount)
}

func TestDeleteHook(db *gorm.DB) {
	//查询文章ID为4的评论列表
	var comments []models.Comment
	err := db.Where(&models.Comment{PostID: 4}).Find(&comments).Error
	if err != nil {
		fmt.Println("查询评论失败:", err)
		return
	}
	fmt.Println("查询到文章ID等于4的评论列表:", comments)

	//遍历单个删除 触发删除钩子
	// Where条件批量删除 钩子不会生效
	// 当有大批量删除的或者经常删除的情况下 也不会在业务表存储状态 而且通过查询另外一个表数据来判断状态
	//如果特频繁还可以通过redis 计数 判断
	//所以这里不处理通过Where批量删除的情况
	for _, comment := range comments {
		err = db.Delete(&comment).Error
		if err != nil {
			fmt.Println("删除评论失败:", err)
			return
		}
		fmt.Println("删除评论:", comment.ToString())
		//查询文章评论状态
		type PostCommentStatus struct {
			CommentStatus uint8 `json:"commentStatus"`
		}
		var postCommentStatus PostCommentStatus
		err = db.Model(&models.Post{}).
			Where("id = ?", comment.PostID).
			Scan(&postCommentStatus).Error
		if err != nil {
			fmt.Println("查询文章评论状态失败:", err)
			return
		}
		if postCommentStatus.CommentStatus == 0 {
			fmt.Printf("文章ID等于%d评论状态:无评论\n", comment.PostID)
		} else {
			fmt.Printf("文章ID等于%d评论状态:有评论\n", comment.PostID)
		}
	}
}
