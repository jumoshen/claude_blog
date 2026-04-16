package dao

import (
	"markdown-blog/internal/repository"
)

// DAOS 数据访问对象集合
type DAOS struct {
	Post     *PostDAO
	User     *UserDAO
	BlogUser *BlogUserDAO
	Comment  *CommentDAO
	Tag      *TagDAO
	Category *CategoryDAO
	Visit    *VisitDAO
	Like     *LikeDAO
	Favorite *FavoriteDAO
}

// NewDAOS 创建 DAOS 实例
func NewDAOS(repo *repository.Repository) *DAOS {
	db := repo.GetDB()
	return &DAOS{
		Post:     NewPostDAO(db),
		User:     NewUserDAO(db),
		BlogUser: NewBlogUserDAO(db),
		Comment:  NewCommentDAO(db),
		Tag:      NewTagDAO(db),
		Category: NewCategoryDAO(db),
		Visit:    NewVisitDAO(db),
		Like:     NewLikeDAO(db),
		Favorite: NewFavoriteDAO(db),
	}
}
