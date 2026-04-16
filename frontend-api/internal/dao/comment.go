package dao

import (
	"gorm.io/gorm"

	"markdown-blog/internal/model"
)

// CommentDAO 评论数据访问对象
type CommentDAO struct {
	*DAO
}

// NewCommentDAO 创建 CommentDAO
func NewCommentDAO(db *gorm.DB) *CommentDAO {
	return &CommentDAO{NewDAO(db)}
}

// GetCommentsBySlug 获取文章的所有评论
func (d *CommentDAO) GetCommentsBySlug(slug string) ([]model.Comment, error) {
	var comments []model.Comment
	if err := d.db.Where("post_slug = ?", slug).Order("created_at DESC").Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// CreateComment 创建评论
func (d *CommentDAO) CreateComment(comment *model.Comment) error {
	return d.db.Create(comment).Error
}

// DeleteComment 删除评论
func (d *CommentDAO) DeleteComment(id uint) error {
	return d.db.Delete(&model.Comment{}, id).Error
}

// GetCommentByID 根据 ID 获取评论
func (d *CommentDAO) GetCommentByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	if err := d.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}
