package main

import (
	"fmt"
	"homework4/config"
	"homework4/internal/api/route"
	"homework4/internal/app/mysql"
	"homework4/internal/app/redis"
	"homework4/internal/common"
	"homework4/pkg/logger"
	"time"

	_ "homework4/docs"
	_ "homework4/internal/controller"
)

// @title Blog API
// @version 1.0
// @description 这是一个博客API系统,包含用户注册登录、文章管理、评论管理等功能
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host ${HOST}:${PORT}
// @BasePath /api/v1
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 请输入JWT token,格式为Bearer {token}
func main() {
	// 初始化日志
	logger.InitAppLog(
		// logger.WithDisableConsole(),  //禁用控制台日志
		logger.WithTimeLayout(time.DateTime),
		logger.WithFileP(common.LogFile),
	)
	//初始化配置
	config.LoadConfig()

	//初始化数据库
	mysql.InitDB()
	//初始化redis
	redis.InitRedis()

	//初始化路由
	r := route.InitRoutes()

	// 启动服务
	port := fmt.Sprintf("%d", config.Cfg.App.Port)

	r.Run(":" + port)

	logger.AppLog.Info("服务启动成功,端口：" + port)
}
