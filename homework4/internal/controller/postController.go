package controller

import (
	"homework4/internal/middleware/auth"
	"homework4/internal/middleware/response"
	"homework4/internal/service"

	"github.com/gin-gonic/gin"
)

type PostController struct {
	postService *service.PostService
}

func NewPostController() *PostController {
	return &PostController{
		postService: service.NewPostService(),
	}
}

/**
 * @Description: 创建文章
 * @param c
 * @return error
 */
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
func (ctrl *PostController) CreatePost(c *gin.Context) error {
	var req service.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}
	//查询登录信息
	authUser := auth.GetCurrentAuthUser(c)
	//创建文章
	post, err := ctrl.postService.CreatePost(&req, authUser.UserID)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}

	response.SendJSON(c, post)
	return nil
}

/**
 * @Description: 更新文章
 * @param c
 * @return error
 */
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
func (ctrl *PostController) UpdatePost(c *gin.Context) error {
	var req service.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}
	//查询登录信息
	authUser := auth.GetCurrentAuthUser(c)
	//更新文章
	post, err := ctrl.postService.UpdatePost(&req, authUser.UserID)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}

	response.SendJSON(c, post)
	return nil
}

/**
 * @Description: 删除文章
 * @param c
 * @return error
 */
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
func (ctrl *PostController) DeletePost(c *gin.Context) error {
	var req service.DeletePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}
	//查询登录信息
	authUser := auth.GetCurrentAuthUser(c)
	//删除文章会校验是不是当前用户的文章
	err := ctrl.postService.DeletePost(&req, authUser.UserID)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}

	response.SendJSON(c, gin.H{
		"message": "删除成功",
	})
	return nil
}

/**
 * @Description: 获取文章分页
 * @param c
 * @return error
 */
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
func (ctrl *PostController) GetPostList(c *gin.Context) error {
	var req service.GetPostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}

	posts, total, err := ctrl.postService.GetPostList(&req)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}

	response.SendJSON(c, gin.H{
		"list":  posts,
		"total": total,
	})
	return nil
}
