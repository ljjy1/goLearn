package controller

import (
	"homework4/internal/middleware/auth"
	"homework4/internal/middleware/response"
	"homework4/internal/service"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentService *service.CommentService
}

func NewCommentController() *CommentController {
	return &CommentController{
		commentService: service.NewCommentService(),
	}
}

/**
 * @Description: 创建评论
 * @param c
 * @return error
 */
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
func (ctrl *CommentController) CreateComment(c *gin.Context) error {
	var req service.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}
	//查询登录信息
	authUser := auth.GetCurrentAuthUser(c)
	//创建评论
	comment, err := ctrl.commentService.CreateComment(&req, authUser.UserID)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}

	response.SendJSON(c, comment)
	return nil
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
func (ctrl *CommentController) GetCommentList(c *gin.Context) error {
	var req service.GetCommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return response.NewBadRequestError("参数错误: " + err.Error())
	}

	comments, total, err := ctrl.commentService.GetCommentList(&req)
	if err != nil {
		return response.NewBizError(response.CodeBadRequest, err.Error(), nil)
	}

	response.SendJSON(c, gin.H{
		"list":  comments,
		"total": total,
	})
	return nil
}
