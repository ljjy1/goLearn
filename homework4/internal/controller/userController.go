package controller

import (
	"homework4/internal/middleware/response"
	"homework4/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

/**
 * @Description: 注册用户
 * @param c
 * @return error
 */
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
func (ctrl *UserController) Register(c *gin.Context) error {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}

	user, err := ctrl.userService.Register(&req)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}

	response.SendJSON(c, gin.H{
		"user_id":  user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"email":    user.Email,
	})
	return nil
}

/**
 * @Description: 用户登录
 * @param c
 * @return error
 */
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
func (ctrl *UserController) Login(c *gin.Context) error {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}

	resp, err := ctrl.userService.Login(&req)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}
	response.SendJSON(c, resp)
	return nil
}
