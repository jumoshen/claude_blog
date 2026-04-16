package dao

import (
	"gorm.io/gorm"

	"markdown-blog/internal/model"
)

// UserDAO 用户数据访问对象
type UserDAO struct {
	*DAO
}

// NewUserDAO 创建 UserDAO
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{NewDAO(db)}
}

// GetUserByGitHubID 根据 GitHub ID 获取用户
func (d *UserDAO) GetUserByGitHubID(githubID int64) (*model.User, error) {
	var user model.User
	if err := d.db.Where("git_hub_id = ?", githubID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID 根据 ID 获取用户
func (d *UserDAO) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := d.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (d *UserDAO) CreateUser(user *model.User) error {
	return d.db.Create(user).Error
}

// UpdateUser 更新用户
func (d *UserDAO) UpdateUser(user *model.User) error {
	return d.db.Save(user).Error
}

// GetAllUsers 获取所有用户
func (d *UserDAO) GetAllUsers() ([]model.User, error) {
	var users []model.User
	if err := d.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// BlogUserDAO 博客用户数据访问对象
type BlogUserDAO struct {
	*DAO
}

// NewBlogUserDAO 创建 BlogUserDAO
func NewBlogUserDAO(db *gorm.DB) *BlogUserDAO {
	return &BlogUserDAO{NewDAO(db)}
}

// GetBlogUserByUsername 根据用户名获取博客用户
func (d *BlogUserDAO) GetBlogUserByUsername(username string) (*model.BlogUser, error) {
	var user model.BlogUser
	if err := d.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetBlogUserByID 根据 ID 获取博客用户
func (d *BlogUserDAO) GetBlogUserByID(id uint) (*model.BlogUser, error) {
	var user model.BlogUser
	if err := d.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateBlogUser 创建博客用户
func (d *BlogUserDAO) CreateBlogUser(user *model.BlogUser) error {
	return d.db.Create(user).Error
}

// UpdateBlogUser 更新博客用户
func (d *BlogUserDAO) UpdateBlogUser(user *model.BlogUser) error {
	return d.db.Save(user).Error
}

// GetAllBlogUsers 获取所有博客用户
func (d *BlogUserDAO) GetAllBlogUsers() ([]model.BlogUser, error) {
	var users []model.BlogUser
	if err := d.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
