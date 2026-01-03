package route

/**
 * @Description: 路由管理
 */
import (
	"homework4/internal/controller"
	"homework4/internal/middleware/auth"
	"homework4/internal/middleware/logger"
	"homework4/internal/middleware/response"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var userController *controller.UserController
var postController *controller.PostController
var commentController *controller.CommentController

// Register godoc
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=map[string]interface{}} "注册成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /user/register [post]
func RegisterHandler(c *gin.Context) {
	response.WrapHandler(userController.Register)(c)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录获取JWT token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body service.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=service.LoginResponse} "登录成功"
// @Failure 400 {object} response.Response "参数错误或登录失败"
// @Router /user/login [post]
func LoginHandler(c *gin.Context) {
	response.WrapHandler(userController.Login)(c)
}

// CreatePost godoc
// @Summary 创建文章
// @Description 创建新文章,需要登录
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param request body service.CreatePostRequest true "文章信息"
// @Security Bearer
// @Success 200 {object} response.Response{data=service.PostResponse} "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /post/create [post]
func CreatePostHandler(c *gin.Context) {
	response.WrapHandler(postController.CreatePost)(c)
}

// UpdatePost godoc
// @Summary 更新文章
// @Description 更新文章信息,需要登录
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param request body service.UpdatePostRequest true "文章信息"
// @Security Bearer
// @Success 200 {object} response.Response{data=service.PostResponse} "更新成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /post/update [put]
func UpdatePostHandler(c *gin.Context) {
	response.WrapHandler(postController.UpdatePost)(c)
}

// DeletePost godoc
// @Summary 删除文章
// @Description 删除文章,需要登录
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param request body service.DeletePostRequest true "文章ID"
// @Security Bearer
// @Success 200 {object} response.Response{data=map[string]interface{}} "删除成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /post/delete [delete]
func DeletePostHandler(c *gin.Context) {
	response.WrapHandler(postController.DeletePost)(c)
}

// GetPostList godoc
// @Summary 获取文章列表
// @Description 分页获取文章列表,不需要登录
// @Tags 文章管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /post/list [get]
func GetPostListHandler(c *gin.Context) {
	response.WrapHandler(postController.GetPostList)(c)
}

// CreateComment godoc
// @Summary 创建评论
// @Description 为文章创建评论,需要登录
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param request body service.CreateCommentRequest true "评论信息"
// @Security Bearer
// @Success 200 {object} response.Response{data=service.CommentResponse} "创建成功"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /comment/create [post]
func CreateCommentHandler(c *gin.Context) {
	response.WrapHandler(commentController.CreateComment)(c)
}

// GetCommentList godoc
// @Summary 获取评论列表
// @Description 分页获取文章评论列表,不需要登录
// @Tags 评论管理
// @Accept json
// @Produce json
// @Param postId query int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Failure 400 {object} response.Response "参数错误"
// @Router /comment/list [get]
func GetCommentListHandler(c *gin.Context) {
	response.WrapHandler(commentController.GetCommentList)(c)
}

func InitRoutes() *gin.Engine {
	// 注册路由
	r := gin.Default()

	// 使用中间件
	r.Use(logger.LoggerMiddleware())
	r.Use(response.ErrorHandlerMiddleware())
	r.Use(gin.Recovery())

	// 初始化服务
	userController = controller.NewUserController()
	postController = controller.NewPostController()
	commentController = controller.NewCommentController()

	// API路由组
	api := r.Group("/api/v1")
	{
		// 用户路由
		userGroup := api.Group("/user")
		{
			userGroup.POST("/register", RegisterHandler)
			userGroup.POST("/login", LoginHandler)
		}

		// 文章路由
		articleGroup := api.Group("/post")
		{
			articleGroup.POST("/create", CreatePostHandler).Use(auth.AuthMiddleware())
			articleGroup.PUT("/update", UpdatePostHandler).Use(auth.AuthMiddleware())
			articleGroup.DELETE("/delete", DeletePostHandler).Use(auth.AuthMiddleware())
			articleGroup.GET("/list", GetPostListHandler)
		}

		// 评论路由
		commentGroup := api.Group("/comment")
		{
			commentGroup.POST("/create", CreateCommentHandler).Use(auth.AuthMiddleware())
			commentGroup.GET("/list", GetCommentListHandler)
		}

		// 健康检查
		// @Summary 健康检查
		// @Description 检查服务是否正常运行
		// @Tags 系统管理
		// @Produce json
		// @Success 200 {object} map[string]string "服务正常"
		// @Router /health [get]
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Blog API is running",
			})
		})

		// Swagger文档路由
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return r
}
