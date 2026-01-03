package service

import (
	"errors"
	"homework4/internal/app/mysql"
	"homework4/internal/models"

	"gorm.io/gorm"
)

type CommentService struct {
	db *gorm.DB
}

func NewCommentService() *CommentService {
	return &CommentService{db: mysql.DB}
}

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	PostID  uint   `json:"postId" binding:"required" example:"1"`                     // 文章ID
	Content string `json:"content" binding:"required,min=1,max=200" example:"很棒的文章!"` // 评论内容
}

// GetCommentListRequest 获取评论列表请求
type GetCommentListRequest struct {
	PostID   uint `form:"postId" binding:"required" example:"1"`                   // 文章ID
	Page     int  `form:"page" binding:"omitempty,min=1" example:"1"`              // 页码
	PageSize int  `form:"pageSize" binding:"omitempty,min=1,max=100" example:"10"` // 每页数量
}

// CommentResponse 评论响应
type CommentResponse struct {
	ID        uint   `json:"id" example:"1"`                          // 评论ID
	PostID    uint   `json:"postId" example:"1"`                      // 文章ID
	UserID    uint   `json:"userId" example:"1"`                      // 用户ID
	Content   string `json:"content" example:"很棒的文章!"`                // 评论内容
	CreatedAt string `json:"createdAt" example:"2024-01-01 12:00:00"` // 创建时间
}

// CommentWithUserResponse 带用户信息的评论响应
type CommentWithUserResponse struct {
	ID        uint   `json:"id" example:"1"`                          // 评论ID
	PostID    uint   `json:"postId" example:"1"`                      // 文章ID
	UserID    uint   `json:"userId" example:"1"`                      // 用户ID
	Username  string `json:"username" example:"testuser"`             // 用户名
	Nickname  string `json:"nickname" example:"测试用户"`                 // 昵称
	Content   string `json:"content" example:"很棒的文章!"`                // 评论内容
	CreatedAt string `json:"createdAt" example:"2024-01-01 12:00:00"` // 创建时间
}

/**
 * @Description: 创建评论
 * @param req
 * @param userID
 * @return (*CommentResponse, error)
 */
func (s *CommentService) CreateComment(req *CreateCommentRequest, userID uint) (*CommentResponse, error) {
	var post models.Post
	if err := s.db.First(&post, req.PostID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文章不存在")
		}
		return nil, err
	}

	comment := &models.Comment{
		UserID:  userID,
		PostID:  req.PostID,
		Content: req.Content,
	}

	if err := s.db.Create(comment).Error; err != nil {
		return nil, err
	}

	return &CommentResponse{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

/**
 * @Description: 获取文章评论分页
 * @param req
 * @return ([]CommentWithUserResponse, int64, error)
 */
func (s *CommentService) GetCommentList(req *GetCommentListRequest) ([]CommentWithUserResponse, int64, error) {
	page := 1
	pageSize := 10

	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	var comments []models.Comment
	var total int64

	query := s.db.Model(&models.Comment{}).Where("post_id = ?", req.PostID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	commentResponses := make([]CommentWithUserResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = CommentWithUserResponse{
			ID:        comment.ID,
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			Username:  comment.User.Username,
			Nickname:  comment.User.Nickname,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return commentResponses, total, nil
}
