package dao

import (
	"gorm.io/gorm"

	"markdown-blog/internal/model"
)

// PostDAO 文章数据访问对象
type PostDAO struct {
	*DAO
}

// NewPostDAO 创建 PostDAO
func NewPostDAO(db *gorm.DB) *PostDAO {
	return &PostDAO{NewDAO(db)}
}

// ListPosts 获取所有已发布文章
func (d *PostDAO) ListPosts() ([]model.Post, error) {
	var posts []model.Post
	if err := d.db.Where("status = 1").Order("date desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// ListPostsPaginated 分页获取已发布文章
func (d *PostDAO) ListPostsPaginated(tag string, category string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := d.db.Model(&model.Post{}).Where("status = 1")

	if tag != "" {
		var tagPostIDs []uint
		if err := d.db.Table("post_tags").
			Select("post_tags.post_id").
			Joins("JOIN tags ON post_tags.tag_id = tags.id").
			Where("tags.name = ?", tag).
			Pluck("post_tags.post_id", &tagPostIDs).Error; err != nil {
			return nil, 0, err
		}
		if len(tagPostIDs) > 0 {
			query = query.Where("id IN ?", tagPostIDs)
		} else {
			return []model.Post{}, 0, nil
		}
	}

	if category != "" {
		var categoryPostIDs []uint
		if err := d.db.Table("post_categories").
			Select("post_categories.post_id").
			Joins("JOIN categories ON post_categories.category_id = categories.id").
			Where("categories.name = ?", category).
			Pluck("post_categories.post_id", &categoryPostIDs).Error; err != nil {
			return nil, 0, err
		}
		if len(categoryPostIDs) > 0 {
			query = query.Where("id IN ?", categoryPostIDs)
		} else {
			return []model.Post{}, 0, nil
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("is_pinned DESC, date DESC").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetPostBySlug 根据slug获取文章
func (d *PostDAO) GetPostBySlug(slug string) (*model.Post, error) {
	var post model.Post
	if err := d.db.Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostByID 根据ID获取文章
func (d *PostDAO) GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := d.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostTags 获取文章标签名称列表
func (d *PostDAO) GetPostTags(postIDs []uint) (map[uint][]string, error) {
	if len(postIDs) == 0 {
		return make(map[uint][]string), nil
	}

	var results []struct {
		PostID uint
		TagID  uint
		TagName string
	}
	if err := d.db.Table("post_tags").
		Select("post_tags.post_id, post_tags.tag_id, tags.name as tag_name").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id IN ?", postIDs).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	result := make(map[uint][]string)
	for _, r := range results {
		result[r.PostID] = append(result[r.PostID], r.TagName)
	}
	return result, nil
}

// GetPostCategories 获取文章分类名称列表
func (d *PostDAO) GetPostCategories(postIDs []uint) (map[uint][]string, error) {
	if len(postIDs) == 0 {
		return make(map[uint][]string), nil
	}

	var results []struct {
		PostID      uint
		CategoryID  uint
		CategoryName string
	}
	if err := d.db.Table("post_categories").
		Select("post_categories.post_id, post_categories.category_id, categories.name as category_name").
		Joins("JOIN categories ON post_categories.category_id = categories.id").
		Where("post_categories.post_id IN ?", postIDs).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	result := make(map[uint][]string)
	for _, r := range results {
		result[r.PostID] = append(result[r.PostID], r.CategoryName)
	}
	return result, nil
}

// CreatePost 创建文章
func (d *PostDAO) CreatePost(post *model.Post) error {
	return d.db.Create(post).Error
}

// UpdatePost 更新文章
func (d *PostDAO) UpdatePost(post *model.Post) error {
	return d.db.Save(post).Error
}

// DeletePost 删除文章
func (d *PostDAO) DeletePost(id uint) error {
	return d.db.Delete(&model.Post{}, id).Error
}

// IncrementViews 增加阅读量
func (d *PostDAO) IncrementViews(slug string) error {
	return d.db.Model(&model.Post{}).Where("slug = ?", slug).UpdateColumn("views", gorm.Expr("views + 1")).Error
}

// ListFeaturedPosts 获取推荐文章
func (d *PostDAO) ListFeaturedPosts() ([]model.Post, error) {
	var posts []model.Post
	if err := d.db.Where("status = 1 AND is_featured = ?", true).Order("date DESC").Limit(5).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// ListPopularPosts 获取热门文章
func (d *PostDAO) ListPopularPosts(limit int) ([]model.Post, error) {
	var posts []model.Post
	if err := d.db.Where("status = 1").Order("views DESC").Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// SearchPosts 搜索文章
func (d *PostDAO) SearchPosts(keyword string) ([]model.Post, error) {
	var posts []model.Post
	search := "%" + keyword + "%"
	if err := d.db.Where("status = 1 AND (title LIKE ? OR content LIKE ? OR summary LIKE ?)", search, search, search).Order("date DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// GetArchives 获取归档数据
func (d *PostDAO) GetArchives() (map[string][]model.Post, error) {
	var posts []model.Post
	if err := d.db.Where("status = 1").Order("date DESC").Find(&posts).Error; err != nil {
		return nil, err
	}

	archives := make(map[string][]model.Post)
	for _, post := range posts {
		year := post.Date.Format("2006")
		archives[year] = append(archives[year], post)
	}
	return archives, nil
}

// SetPostTags 设置文章标签
func (d *PostDAO) SetPostTags(postID uint, tagIDs []uint) error {
	// 删除现有标签
	if err := d.db.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
		return err
	}
	// 添加新标签
	for _, tagID := range tagIDs {
		if err := d.db.Create(&model.PostTag{PostID: postID, TagID: tagID}).Error; err != nil {
			return err
		}
	}
	return nil
}

// SetPostCategories 设置文章分类
func (d *PostDAO) SetPostCategories(postID uint, categoryID *uint) error {
	// 删除现有分类
	if err := d.db.Where("post_id = ?", postID).Delete(&model.PostCategory{}).Error; err != nil {
		return err
	}
	// 添加新分类
	if categoryID != nil {
		return d.db.Create(&model.PostCategory{PostID: postID, CategoryID: *categoryID}).Error
	}
	return nil
}

// GetAllTags 获取所有标签
func (d *PostDAO) GetAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := d.db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetAllCategories 获取所有分类
func (d *PostDAO) GetAllCategories() ([]model.Category, error) {
	var categories []model.Category
	if err := d.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetTagsByPostID 获取文章标签
func (d *PostDAO) GetTagsByPostID(postID uint) ([]model.Tag, error) {
	var tags []model.Tag
	if err := d.db.Table("tags").
		Joins("JOIN post_tags ON tags.id = post_tags.tag_id").
		Where("post_tags.post_id = ?", postID).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetCategoriesByPostID 获取文章分类
func (d *PostDAO) GetCategoriesByPostID(postID uint) ([]model.Category, error) {
	var categories []model.Category
	if err := d.db.Table("categories").
		Joins("JOIN post_categories ON categories.id = post_categories.category_id").
		Where("post_categories.post_id = ?", postID).
		Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetPrevNextPost 获取上一篇和下一篇文章
func (d *PostDAO) GetPrevNextPost(id uint) (*model.Post, *model.Post, error) {
	var posts []model.Post
	if err := d.db.Where("status = 1 AND id > ?", id).Order("id ASC").Limit(1).Find(&posts).Error; err != nil {
		return nil, nil, err
	}
	var next *model.Post
	if len(posts) > 0 {
		next = &posts[0]
	}

	var prevPosts []model.Post
	if err := d.db.Where("status = 1 AND id < ?", id).Order("id DESC").Limit(1).Find(&prevPosts).Error; err != nil {
		return nil, nil, err
	}
	var prev *model.Post
	if len(prevPosts) > 0 {
		prev = &prevPosts[0]
	}

	return prev, next, nil
}

// ListRelatedPosts 获取相关文章
func (d *PostDAO) ListRelatedPosts(postID uint, limit int) ([]model.Post, error) {
	var posts []model.Post
	if err := d.db.Where("status = 1 AND id != ?", postID).Order("date DESC").Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}
