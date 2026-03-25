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

	// Auto migrate
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
	if err := r.db.Order("date desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
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

// GetAllTags 获取所有标签及其文章数量
func (r *Repository) GetAllTags() (map[string]int, error) {
	var posts []model.Post
	if err := r.db.Order("date desc").Find(&posts).Error; err != nil {
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

// GetAllCategories 获取所有分类及其文章数量
func (r *Repository) GetAllCategories() (map[string]int, error) {
	var posts []model.Post
	if err := r.db.Order("date desc").Find(&posts).Error; err != nil {
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
