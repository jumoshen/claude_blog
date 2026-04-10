package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"markdown-blog/internal/config"
	"markdown-blog/internal/model"
)

type Repository struct {
	db    *gorm.DB
	redis *redis.Client
}

func New(cfg *config.Config) (*Repository, error) {
	// Connect to MySQL
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate (only for basic tables, use migrations/ for schema changes)
	// db.AutoMigrate(&model.Post{}, &model.User{}, &model.Visit{})

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetRedisAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test Redis connection
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		// Log warning but continue - Redis is optional for basic operations
	}

	return &Repository{
		db:    db,
		redis: rdb,
	}, nil
}

func (r *Repository) GetDB() *gorm.DB {
	return r.db
}

func (r *Repository) GetRedis() *redis.Client {
	return r.redis
}

// Post operations
func (r *Repository) ListPosts() ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.Where("status = 1").Order("date desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// ListPostsPaginated 分页获取已发布文章
func (r *Repository) ListPostsPaginated(tag string, category string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("status = 1")

	// 如果按标签筛选，使用新的 normalized 表
	if tag != "" {
		var tagPostIDs []uint
		if err := r.db.Table("post_tags").
			Select("post_tags.post_id").
			Joins("JOIN tags ON post_tags.tag_id = tags.id").
			Joins("JOIN posts ON post_tags.post_id = posts.id").
			Where("tags.slug = ? AND posts.status = 1", tag).
			Find(&tagPostIDs).Error; err != nil {
			return nil, 0, err
		}
		if len(tagPostIDs) == 0 {
			return posts, 0, nil
		}
		query = query.Where("id IN ?", tagPostIDs)
	}

	// 如果按分类筛选
	if category != "" {
		var catPostIDs []uint
		if err := r.db.Table("post_categories").
			Select("post_categories.post_id").
			Joins("JOIN categories ON post_categories.category_id = categories.id").
			Joins("JOIN posts ON post_categories.post_id = posts.id").
			Where("categories.slug = ? AND posts.status = 1", category).
			Find(&catPostIDs).Error; err != nil {
			return nil, 0, err
		}
		if len(catPostIDs) == 0 {
			return posts, 0, nil
		}
		query = query.Where("id IN ?", catPostIDs)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.Order("date desc").Offset(offset).Limit(pageSize).Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *Repository) GetPostBySlug(slug string) (*model.Post, error) {
	var post model.Post
	if err := r.db.Where("slug = ?", slug).First(&post).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *Repository) CreatePost(post *model.Post) error {
	return r.db.Create(post).Error
}

func (r *Repository) UpdatePost(post *model.Post) error {
	return r.db.Save(post).Error
}

func (r *Repository) UpsertPost(post *model.Post) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "slug"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "date", "tags", "categories", "summary", "content", "updated_at"}),
	}).Create(post).Error
}

// User operations
func (r *Repository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *Repository) UpdateUser(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *Repository) GetUserByGitHubID(id int64) (*model.User, error) {
	var user model.User
	if err := r.db.Where("git_hub_id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Blacklist operations
func (r *Repository) AddToBlacklist(ctx context.Context, jti string, exp time.Duration) error {
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	return r.redis.Set(ctx, key, "1", exp).Err()
}

func (r *Repository) IsBlacklisted(ctx context.Context, jti string) bool {
	key := fmt.Sprintf("jwt:blacklist:%s", jti)
	result, err := r.redis.Exists(ctx, key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

// TagWithCount 标签及其数量
type TagWithCount struct {
	Slug  string `json:"slug"`
	Count int    `json:"count"`
	Color string `json:"color"`
}

// GetAllTags 获取所有标签及其文章数量（从 normalized 表查询）
func (r *Repository) GetAllTags() ([]TagWithCount, error) {
	var results []TagWithCount
	if err := r.db.Table("tags").
		Select("tags.slug, COUNT(post_tags.post_id) as count, tags.color").
		Joins("LEFT JOIN post_tags ON tags.id = post_tags.tag_id").
		Joins("LEFT JOIN posts ON post_tags.post_id = posts.id AND posts.status = 1").
		Group("tags.id").
		Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetAllCategories 获取所有分类及其文章数量（从 normalized 表查询）
func (r *Repository) GetAllCategories() (map[string]int, error) {
	var results []struct {
		Slug string
		Count int
	}
	if err := r.db.Table("categories").
		Select("categories.slug, COUNT(post_categories.post_id) as count").
		Joins("LEFT JOIN post_categories ON categories.id = post_categories.category_id").
		Joins("LEFT JOIN posts ON post_categories.post_id = posts.id AND posts.status = 1").
		Group("categories.id").
		Find(&results).Error; err != nil {
		return nil, err
	}

	catCount := make(map[string]int)
	for _, c := range results {
		if c.Slug != "" {
			catCount[c.Slug] = c.Count
		}
	}
	return catCount, nil
}

// GetPostTags 获取指定文章的标签slug列表
func (r *Repository) GetPostTags(postIDs []uint) (map[uint][]string, error) {
	if len(postIDs) == 0 {
		return make(map[uint][]string), nil
	}

	var results []struct {
		PostID uint
		Slug   string
	}
	if err := r.db.Table("post_tags").
		Select("post_tags.post_id, tags.slug").
		Joins("JOIN tags ON post_tags.tag_id = tags.id").
		Where("post_tags.post_id IN ?", postIDs).
		Find(&results).Error; err != nil {
		return nil, err
	}

	tagMap := make(map[uint][]string)
	for _, r := range results {
		tagMap[r.PostID] = append(tagMap[r.PostID], r.Slug)
	}
	return tagMap, nil
}

// GetPostCategories 获取指定文章的分类slug列表
func (r *Repository) GetPostCategories(postIDs []uint) (map[uint][]string, error) {
	if len(postIDs) == 0 {
		return make(map[uint][]string), nil
	}

	var results []struct {
		PostID uint
		Slug   string
	}
	if err := r.db.Table("post_categories").
		Select("post_categories.post_id, categories.slug").
		Joins("JOIN categories ON post_categories.category_id = categories.id").
		Where("post_categories.post_id IN ?", postIDs).
		Find(&results).Error; err != nil {
		return nil, err
	}

	catMap := make(map[uint][]string)
	for _, r := range results {
		catMap[r.PostID] = append(catMap[r.PostID], r.Slug)
	}
	return catMap, nil
}

// CreateVisit 创建访问记录
func (r *Repository) CreateVisit(visit *model.Visit) error {
	return r.db.Create(visit).Error
}

// CanRecordVisit 检查是否可以记录访问（防刷）
func (r *Repository) CanRecordVisit(ctx context.Context, postSlug string, ip string, userID int64) (bool, error) {
	key := fmt.Sprintf("visit:%s:%s:%d", postSlug, ip, userID)

	// 使用 SETNX，1小时内不重复记录
	exists, err := r.redis.SetNX(ctx, key, "1", time.Hour).Result()
	if err != nil {
		// Redis 错误时允许记录
		return true, nil
	}
	return exists, nil
}

// Comment operations
func (r *Repository) CreateComment(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *Repository) GetCommentsByPostSlug(postSlug string, limit int) ([]model.Comment, error) {
	var comments []model.Comment
	if err := r.db.Where("post_slug = ? AND status = 1", postSlug).
		Order("created_at DESC").
		Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

// IncrCommentRateLimit 增加评论频率限制计数
func (r *Repository) IncrCommentRateLimit(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := r.redis.Pipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}

// GetCommentRateLimit 获取评论频率限制计数
func (r *Repository) GetCommentRateLimit(ctx context.Context, key string) (int64, error) {
	val, err := r.redis.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return val, err
}

// SensitiveWord operations
func (r *Repository) GetAllSensitiveWords() ([]model.SensitiveWord, error) {
	var words []model.SensitiveWord
	if err := r.db.Where("deleted_at IS NULL").Find(&words).Error; err != nil {
		return nil, err
	}
	return words, nil
}

func (r *Repository) SetSensitiveWordsCache(ctx context.Context, words []string) error {
	if len(words) == 0 {
		return nil
	}
	key := "sensitive_words:cache"

	// 使用Hash存储，每个词一个field，避免大key问题
	pipe := r.redis.Pipeline()
	pipe.Del(ctx, key) // 先删除旧数据
	for _, word := range words {
		pipe.HSet(ctx, key, word, 1)
	}
	pipe.Expire(ctx, key, time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *Repository) GetSensitiveWordsCache(ctx context.Context) ([]string, error) {
	key := "sensitive_words:cache"
	result, err := r.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// 如果缓存不存在或为空，返回nil让调用方去数据库读取
	if len(result) == 0 {
		return nil, nil
	}

	words := make([]string, 0, len(result))
	for word := range result {
		words = append(words, word)
	}
	return words, nil
}

func (r *Repository) InvalidateSensitiveWordsCache(ctx context.Context) error {
	key := "sensitive_words:cache"
	return r.redis.Del(ctx, key).Err()
}

// Category operations
func (r *Repository) ListCategories(page, pageSize int) ([]model.Category, int64, error) {
	var categories []model.Category
	var total int64

	query := r.db.Model(&model.Category{}).Where("deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("sort ASC, id ASC").Offset(offset).Limit(pageSize).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *Repository) ListAllCategories() ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.Where("deleted_at IS NULL").Order("sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *Repository) GetCategoryByID(id uint) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) GetCategoryBySlug(slug string) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("slug = ? AND deleted_at IS NULL", slug).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) CreateCategory(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *Repository) UpdateCategory(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *Repository) DeleteCategory(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *Repository) CountPostsByCategory(categoryID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Post{}).Where("category_id = ? AND deleted_at IS NULL", categoryID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Tag operations
func (r *Repository) ListTags(page, pageSize int) ([]model.Tag, int64, error) {
	var tags []model.Tag
	var total int64

	query := r.db.Model(&model.Tag{}).Where("deleted_at IS NULL")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	return tags, total, nil
}

func (r *Repository) ListAllTags() ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.Where("deleted_at IS NULL").Order("id ASC").Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *Repository) GetTagByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.Where("id = ? AND deleted_at IS NULL", id).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *Repository) GetTagBySlug(slug string) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.Where("slug = ? AND deleted_at IS NULL", slug).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *Repository) CreateTag(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *Repository) UpdateTag(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

func (r *Repository) DeleteTag(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

func (r *Repository) CountPostsByTag(tagID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.PostTag{}).Where("tag_id = ?", tagID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) GetTagIDsByPostID(postID uint) ([]uint, error) {
	var postTags []model.PostTag
	if err := r.db.Where("post_id = ?", postID).Find(&postTags).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(postTags))
	for _, pt := range postTags {
		ids = append(ids, pt.TagID)
	}
	return ids, nil
}

func (r *Repository) SetPostTags(postID uint, tagIDs []uint) error {
	// 删除旧的关联
	if err := r.db.Where("post_id = ?", postID).Delete(&model.PostTag{}).Error; err != nil {
		return err
	}
	// 创建新的关联
	for _, tagID := range tagIDs {
		pt := model.PostTag{PostID: postID, TagID: tagID}
		if err := r.db.Create(&pt).Error; err != nil {
			return err
		}
	}
	return nil
}

// SearchPosts 搜索文章（标题和内容）
func (r *Repository) SearchPosts(keyword string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("status = 1")

	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("title LIKE ? OR content LIKE ?", keyword, keyword)
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

// ListPostsPinned 获取置顶文章
func (r *Repository) ListPostsPinned() ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.Where("status = 1 AND is_pinned = ?", true).Order("pinned_at DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// ListPostsFeatured 获取推荐文章
func (r *Repository) ListPostsFeatured() ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.Where("status = 1 AND is_featured = ?", true).Order("date DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// SetPostPinned 设置/取消置顶
func (r *Repository) SetPostPinned(id uint, pinned bool) error {
	var pinnedAt *time.Time
	if pinned {
		now := time.Now()
		pinnedAt = &now
	}
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_pinned": pinned,
		"pinned_at": pinnedAt,
	}).Error
}

// SetPostFeatured 设置/取消推荐
func (r *Repository) SetPostFeatured(id uint, featured bool) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Update("is_featured", featured).Error
}

// SchedulePost 设置定时发布
func (r *Repository) SchedulePost(id uint, scheduledAt time.Time) error {
	return r.db.Model(&model.Post{}).Where("id = ?", id).Updates(map[string]interface{}{
		"scheduled_at": scheduledAt,
		"status":        3, // 待发布
	}).Error
}

// PublishScheduledPosts 发布已到时的定时文章
func (r *Repository) PublishScheduledPosts() (int64, error) {
	result := r.db.Model(&model.Post{}).
		Where("status = ? AND scheduled_at <= ?", 3, time.Now()).
		Updates(map[string]interface{}{
			"status":        1, // 已发布
			"scheduled_at":  nil,
		})
	return result.RowsAffected, result.Error
}

// GetPostByID 根据ID获取文章
func (r *Repository) GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := r.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPrevPost 获取上一篇文章（按日期排序，当前文章之前的）
func (r *Repository) GetPrevPost(currentDate time.Time, currentSlug string) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("status = 1 AND date < ? AND slug != ?", currentDate, currentSlug).
		Order("date DESC").First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &post, err
}

// GetNextPost 获取下一篇文章（按日期排序，当前文章之后的）
func (r *Repository) GetNextPost(currentDate time.Time, currentSlug string) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("status = 1 AND date > ? AND slug != ?", currentDate, currentSlug).
		Order("date ASC").First(&post).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &post, err
}

// ListPopularPosts 获取热门文章（按浏览量排序）
func (r *Repository) ListPopularPosts(limit int) ([]model.Post, error) {
	var posts []model.Post
	if err := r.db.Where("status = 1").Order("views DESC").Limit(limit).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// ListRelatedPosts 获取相关文章（共享标签）
func (r *Repository) ListRelatedPosts(currentSlug string, tags []string, limit int) ([]model.Post, error) {
	var posts []model.Post

	// 如果没有标签，返回空
	if len(tags) == 0 {
		return posts, nil
	}

	// 从 post_tags 表查询共享标签的文章
	var relatedPostIDs []uint
	if err := r.db.Table("post_tags").
		Select("DISTINCT post_id").
		Where("post_id IN (SELECT post_id FROM post_tags WHERE tag_id IN (SELECT id FROM tags WHERE slug IN ?))", tags).
		Where("post_id NOT IN (SELECT id FROM posts WHERE slug = ?)", currentSlug).
		Limit(limit).
		Find(&relatedPostIDs).Error; err != nil {
		return nil, err
	}

	if len(relatedPostIDs) == 0 {
		return posts, nil
	}

	if err := r.db.Where("id IN ?", relatedPostIDs).Order("views DESC").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// User operations for C端用户
func (r *Repository) CreateBlogUser(user *model.BlogUser) error {
	return r.db.Create(user).Error
}

func (r *Repository) GetBlogUserByUsername(username string) (*model.BlogUser, error) {
	var user model.BlogUser
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetBlogUserByEmail(email string) (*model.BlogUser, error) {
	var user model.BlogUser
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetBlogUserByID(id uint) (*model.BlogUser, error) {
	var user model.BlogUser
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// PostLike operations
func (r *Repository) AddPostLike(postID, userID uint) error {
	like := &model.PostLike{
		PostID: postID,
		UserID: userID,
	}
	return r.db.Create(like).Error
}

func (r *Repository) RemovePostLike(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostLike{}).Error
}

func (r *Repository) GetPostLike(postID, userID uint) (*model.PostLike, error) {
	var like model.PostLike
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &like, err
}

func (r *Repository) CountPostLikes(postID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.PostLike{}).Where("post_id = ?", postID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// PostFavorite operations
func (r *Repository) AddPostFavorite(postID, userID uint) error {
	fav := &model.PostFavorite{
		PostID: postID,
		UserID: userID,
	}
	return r.db.Create(fav).Error
}

func (r *Repository) RemovePostFavorite(postID, userID uint) error {
	return r.db.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&model.PostFavorite{}).Error
}

func (r *Repository) GetPostFavorite(postID, userID uint) (*model.PostFavorite, error) {
	var fav model.PostFavorite
	err := r.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&fav).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &fav, err
}

func (r *Repository) ListUserFavorites(userID uint) ([]model.Post, error) {
	var favorites []model.PostFavorite
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		return nil, err
	}

	var postIDs []uint
	for _, f := range favorites {
		postIDs = append(postIDs, f.PostID)
	}

	if len(postIDs) == 0 {
		return []model.Post{}, nil
	}

	var posts []model.Post
	if err := r.db.Where("id IN ? AND status = 1", postIDs).Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// HasPassword 检查文章是否有密码保护
func (r *Repository) HasPassword(slug string) (bool, error) {
	var post model.Post
	if err := r.db.Select("password_hash").Where("slug = ?", slug).First(&post).Error; err != nil {
		return false, err
	}
	return post.PasswordHash != "", nil
}

// AdminLog operations
func (r *Repository) CreateAdminLog(log *model.AdminLog) error {
	return r.db.Create(log).Error
}

func (r *Repository) ListAdminLogs(page, pageSize int) ([]model.AdminLog, int64, error) {
	var logs []model.AdminLog
	var total int64

	query := r.db.Model(&model.AdminLog{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
