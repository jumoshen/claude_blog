package repository

import (
	"context"
	"fmt"
	"strings"
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
	db.AutoMigrate(&model.Post{}, &model.User{}, &model.Visit{})

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
func (r *Repository) ListPostsPaginated(tag string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	query := r.db.Model(&model.Post{}).Where("status = 1")

	if tag != "" {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
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

// GetAllTags 获取所有标签及其文章数量（仅已发布文章，使用 SQL 统计）
func (r *Repository) GetAllTags() (map[string]int, error) {
	var posts []model.Post
	if err := r.db.Select("tags").Where("status = 1").Find(&posts).Error; err != nil {
		return nil, err
	}

	tagCount := make(map[string]int)
	for _, p := range posts {
		tags := strings.Split(p.Tags, ",")
		for _, tag := range tags {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				tagCount[tag]++
			}
		}
	}
	return tagCount, nil
}

// GetAllCategories 获取所有分类及其文章数量（仅已发布文章）
func (r *Repository) GetAllCategories() (map[string]int, error) {
	var posts []model.Post
	if err := r.db.Select("categories").Where("status = 1").Find(&posts).Error; err != nil {
		return nil, err
	}

	catCount := make(map[string]int)
	for _, p := range posts {
		cats := strings.Split(p.Categories, ",")
		for _, cat := range cats {
			cat = strings.TrimSpace(cat)
			if cat != "" {
				catCount[cat]++
			}
		}
	}
	return catCount, nil
}

// CreateVisit 创建访问记录
func (r *Repository) CreateVisit(visit *model.Visit) error {
	return r.db.Create(visit).Error
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
