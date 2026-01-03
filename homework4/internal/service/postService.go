package service

import (
	"errors"
	"homework4/internal/app/mysql"
	"homework4/internal/models"

	"time"

	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService() *PostService {
	return &PostService{db: mysql.DB}
}

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=20" example:"我的第一篇文章"`      // 标题
	Content string `json:"content" binding:"required,min=1,max=200" example:"这是文章内容..."` // 内容
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	PostID  uint   `json:"postId" binding:"required" example:"1"`                         // 文章ID
	Title   string `json:"title" binding:"omitempty,min=1,max=20" example:"更新后的标题"`       // 标题
	Content string `json:"content" binding:"omitempty,min=1,max=200" example:"更新后的内容..."` // 内容
}

// DeletePostRequest 删除文章请求
type DeletePostRequest struct {
	PostID uint `json:"postId" binding:"required" example:"1"` // 文章ID
}

// GetPostListRequest 获取文章列表请求
type GetPostListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1" example:"1"`              // 页码
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100" example:"10"` // 每页数量
}

// PostResponse 文章响应
type PostResponse struct {
	ID        uint   `json:"id" example:"1"`                          // 文章ID
	UserID    uint   `json:"userId" example:"1"`                      // 用户ID
	Title     string `json:"title" example:"我的第一篇文章"`                 // 标题
	Content   string `json:"content" example:"这是文章内容..."`             // 内容
	CreatedAt string `json:"createdAt" example:"2024-01-01 12:00:00"` // 创建时间
	UpdatedAt string `json:"updatedAt" example:"2024-01-01 12:00:00"` // 更新时间
}

/**
 * @Description: 创建文章
 * @param req
 * @param userID
 * @return (*PostResponse, error)
 */
func (s *PostService) CreatePost(req *CreatePostRequest, userID uint) (*PostResponse, error) {
	post := &models.Post{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := s.db.Create(post).Error; err != nil {
		return nil, err
	}

	return &PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		CreatedAt: post.CreatedAt.Format(time.DateTime),
		UpdatedAt: post.UpdatedAt.Format(time.DateTime),
	}, nil
}

/**
 * @Description: 更新文章
 * @param req
 * @param userID
 * @return (bool, error)
 */
func (s *PostService) UpdatePost(req *UpdatePostRequest, userID uint) (bool, error) {
	var post models.Post
	if err := s.db.First(&post, req.PostID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("文章不存在")
		}
		return false, err
	}

	if post.UserID != userID {
		return false, errors.New("无权修改此文章")
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}

	if len(updates) == 0 {
		return false, nil
	}
	if err := s.db.Model(&post).Updates(updates).Error; err != nil {
		return false, err
	}
	return true, nil
}

/**
 * @Description: 删除文章
 * @param req
 * @param userID
 * @return error
 */
func (s *PostService) DeletePost(req *DeletePostRequest, userID uint) error {
	var post models.Post
	if err := s.db.First(&post, req.PostID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	if post.UserID != userID {
		return errors.New("无权删除此文章")
	}

	if err := s.db.Delete(&post).Error; err != nil {
		return err
	}

	return nil
}

/**
 * @Description: 获取文章分页
 * @param req
 * @return ([]PostResponse, int64, error)
 */
func (s *PostService) GetPostList(req *GetPostListRequest) ([]PostResponse, int64, error) {
	page := 1
	pageSize := 10

	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	var posts []models.Post
	var total int64

	if err := s.db.Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := s.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	postResponses := make([]PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = PostResponse{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return postResponses, total, nil
}
