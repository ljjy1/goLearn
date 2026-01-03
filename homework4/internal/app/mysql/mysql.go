package mysql

/**
 * @Description: 初始化数据库连接
 */
import (
	"fmt"
	"homework4/config"
	"homework4/internal/models"
	"homework4/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB() {
	//设置表名规则
	namingStrategy := schema.NamingStrategy{
		TablePrefix:   "table_", // 表名加table_前缀
		SingularTable: true,     // 禁用复数表名（去掉s后缀）
	}
	// 构建DSN
	dsn := config.Cfg.Database.User + ":" + config.Cfg.Database.Password + "@tcp(" + config.Cfg.Database.Host + ":" + fmt.Sprintf("%d", config.Cfg.Database.Port) + ")/" + config.Cfg.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		//禁用配置外键
		DisableForeignKeyConstraintWhenMigrating: true,
		//设置表名规则
		NamingStrategy: namingStrategy,
	})
	if err != nil {
		logger.AppLog.Fatal("初始化数据库连接失败", logger.WrapMeta(err)...)
	}
	logger.AppLog.Info("数据库连接成功")

	logger.AppLog.Info("开始迁移模型--------------------")
	DB.AutoMigrate(&models.Post{}, &models.Comment{}, &models.User{})

}
