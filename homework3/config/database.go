package config

import (
	"fmt"
	"homework3/models"
	"homework3/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

/**
 * @description: 初始化数据库连接
 * @return {*gorm.DB} 数据库连接实例
 */
func InitDB() {
	//加载环境变量
	utils.LoadEnvFile()
	host := utils.GetEnv("DB_HOST", "localhost")
	port := utils.GetEnv("DB_PORT", "3306")
	user := utils.GetEnv("DB_USER", "root")
	password := utils.GetEnv("DB_PASSWORD", "password")
	database := utils.GetEnv("DB_DATABASE", "test")

	//设置表名规则
	namingStrategy := schema.NamingStrategy{
		TablePrefix:   "table_", // 表名加table_前缀
		SingularTable: true,     // 禁用复数表名（去掉s后缀）
	}

	//拼接链接字符串
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//禁用配置外键
		DisableForeignKeyConstraintWhenMigrating: true,
		//设置表名规则
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		panic("连接数据库失败:" + err.Error())
	}
	fmt.Println("数据库连接成功")

	//创建模型表
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.InitDataConfig{})

	//查询InitDataConfig表数据
	var initDataConfig models.InitDataConfig
	initDataErr := DB.First(&initDataConfig).Error
	initData := true
	if initDataErr != nil {
		initData = false
	}
	if !initData || initDataConfig.ID == 0 || !initDataConfig.InitDefaultData {
		InitDefaultData()
		//更新InitDataConfig表数据
		initDataConfig.InitDefaultData = true
		DB.Save(&initDataConfig)
	}
}

func GetDB() *gorm.DB {
	if DB == nil {
		panic("数据库连接未初始化")
	}
	return DB
}

func InitDefaultData() {
	//初始化默认数据
	age := uint8(20)
	DB.Create(&models.User{
		Nickname: "测试用户1",
		Mobile:   "13712341234",
		Age:      &age,
		Posts: []models.Post{
			{
				Title:         "测试用户1文章1",
				Content:       "文章1内容",
				CommentStatus: 0,
			},
			{
				Title:         "测试用户1文章2",
				Content:       "文章2内容",
				CommentStatus: 1,
				Comments: []models.Comment{
					{
						Content: "文章2评论内容1",
					},
					{
						Content: "文章2评论内容2",
					},
				},
			},
		},
	})

	//创建第二个用户数据
	DB.Create(&models.User{
		Nickname: "测试用户2",
		Mobile:   "13712351235",
		Posts: []models.Post{
			{
				Title:         "测试用户2文章1",
				Content:       "测试用户2文章1内容",
				CommentStatus: 1,
				Comments: []models.Comment{
					{
						Content: "文章1评论内容1",
					},
					{
						Content: "文章1评论内容2",
					},
					{
						Content: "文章1评论内容3",
					},
				},
			},
		},
	})

	fmt.Println("初始化默认数据完成")
}
