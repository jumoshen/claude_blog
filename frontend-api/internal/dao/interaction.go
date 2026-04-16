package dao

import (
	"gorm.io/gorm"

	"markdown-blog/internal/model"
)

// LikeDAO 点赞数据访问对象
type LikeDAO struct {
	*DAO
}

// NewLikeDAO 创建 LikeDAO
func NewLikeDAO(db *gorm.DB) *LikeDAO {
	return &LikeDAO{NewDAO(db)}
}

// AddPostLike 添加点赞
func (d *LikeDAO) AddPostLike(postID, userID uint) error {
	like := &model.PostLike{
		PostID: postID,
		UserID: userID,
	}
	return d.db.Create(like).Error
}

// RemovePostLike 移除点赞
func (d *LikeDAO) RemovePostLike(postID, userID uint) error {
	return d.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostLike{}).Error
}

// GetPostLike 获取点赞记录
func (d *LikeDAO) GetPostLike(postID, userID uint) (*model.PostLike, error) {
	var like model.PostLike
	err := d.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &like, nil
}

// CountPostLikes 统计文章点赞数
func (d *LikeDAO) CountPostLikes(postID uint) (int64, error) {
	var count int64
	if err := d.db.Model(&model.PostLike{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FavoriteDAO 收藏数据访问对象
type FavoriteDAO struct {
	*DAO
}

// NewFavoriteDAO 创建 FavoriteDAO
func NewFavoriteDAO(db *gorm.DB) *FavoriteDAO {
	return &FavoriteDAO{NewDAO(db)}
}

// AddPostFavorite 添加收藏
func (d *FavoriteDAO) AddPostFavorite(postID, userID uint) error {
	fav := &model.PostFavorite{
		PostID: postID,
		UserID: userID,
	}
	return d.db.Create(fav).Error
}

// RemovePostFavorite 移除收藏
func (d *FavoriteDAO) RemovePostFavorite(postID, userID uint) error {
	return d.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostFavorite{}).Error
}

// GetPostFavorite 获取收藏记录
func (d *FavoriteDAO) GetPostFavorite(postID, userID uint) (*model.PostFavorite, error) {
	var fav model.PostFavorite
	err := d.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&fav).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &fav, nil
}

// ListUserFavorites 获取用户的所有收藏
func (d *FavoriteDAO) ListUserFavorites(userID uint) ([]model.PostFavorite, error) {
	var favorites []model.PostFavorite
	if err := d.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, err
	}
	return favorites, nil
}
