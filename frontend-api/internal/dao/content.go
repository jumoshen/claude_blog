package dao

import (
	"gorm.io/gorm"

	"markdown-blog/internal/model"
)

// TagDAO 标签数据访问对象
type TagDAO struct {
	*DAO
}

// NewTagDAO 创建 TagDAO
func NewTagDAO(db *gorm.DB) *TagDAO {
	return &TagDAO{NewDAO(db)}
}

// GetTagByName 根据名称获取标签
func (d *TagDAO) GetTagByName(name string) (*model.Tag, error) {
	var tag model.Tag
	if err := d.db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// CreateTag 创建标签
func (d *TagDAO) CreateTag(tag *model.Tag) error {
	return d.db.Create(tag).Error
}

// GetOrCreateTag 获取或创建标签
func (d *TagDAO) GetOrCreateTag(name string) (*model.Tag, error) {
	var tag model.Tag
	err := d.db.Where("name = ?", name).First(&tag).Error
	if err == nil {
		return &tag, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	tag.Name = name
	if err := d.db.Create(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

// GetAllTags 获取所有标签
func (d *TagDAO) GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := d.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// CategoryDAO 分类数据访问对象
type CategoryDAO struct {
	*DAO
}

// NewCategoryDAO 创建 CategoryDAO
func NewCategoryDAO(db *gorm.DB) *CategoryDAO {
	return &CategoryDAO{NewDAO(db)}
}

// GetCategoryByName 根据名称获取分类
func (d *CategoryDAO) GetCategoryByName(name string) (*model.Category, error) {
	var category model.Category
	if err := d.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// CreateCategory 创建分类
func (d *CategoryDAO) CreateCategory(category *model.Category) error {
	return d.db.Create(category).Error
}

// GetOrCreateCategory 获取或创建分类
func (d *CategoryDAO) GetOrCreateCategory(name string) (*model.Category, error) {
	var category model.Category
	err := d.db.Where("name = ?", name).First(&category).Error
	if err == nil {
		return &category, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	category.Name = name
	if err := d.db.Create(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// GetAllCategories 获取所有分类
func (d *CategoryDAO) GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	if err := d.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// VisitDAO 访问记录数据访问对象
type VisitDAO struct {
	*DAO
}

// NewVisitDAO 创建 VisitDAO
func NewVisitDAO(db *gorm.DB) *VisitDAO {
	return &VisitDAO{NewDAO(db)}
}

// CreateVisit 创建访问记录
func (d *VisitDAO) CreateVisit(visit *model.Visit) error {
	return d.db.Create(visit).Error
}

// GetVisitsByPostSlug 获取文章的所有访问记录
func (d *VisitDAO) GetVisitsByPostSlug(slug string) ([]model.Visit, error) {
	var visits []model.Visit
	if err := d.db.Where("post_slug = ?", slug).Order("created_at DESC").Find(&visits).Error; err != nil {
		return nil, err
	}
	return visits, nil
}
