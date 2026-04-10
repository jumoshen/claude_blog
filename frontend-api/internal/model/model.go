package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GitHubID  int64  `gorm:"uniqueIndex"`
	Login     string
	Name      string
	AvatarURL string
	Email     string
}

type Post struct {
	gorm.Model
	Slug       string `gorm:"uniqueIndex;size:255"`
	Title      string `gorm:"size:255"`
	Date       time.Time
	Tags       string `gorm:"type:text"`
	Categories string `gorm:"type:text"`
	Summary    string `gorm:"type:text"`
	Content    string `gorm:"type:longtext"`
	Views      int64  `gorm:"default:0"`
	Status     int    `gorm:"default:1;index"`       // 0=草稿 1=已发布 2=下架 3=待发布
	CategoryID *uint  `gorm:"index;comment:分类ID"` // 关联分类
	IsPinned    bool        `gorm:"default:false;comment:是否置顶"`
	IsFeatured bool        `gorm:"default:false;comment:是否推荐"`
	PinnedAt   *time.Time  `gorm:"comment:置顶时间"`
	ScheduledAt *time.Time  `gorm:"comment:定时发布时间"`
	PasswordHash string     `gorm:"size:255;default:'';comment:文章密码保护"`
}

type Visit struct {
	gorm.Model
	PostSlug  string `gorm:"size:200;index"`
	UserID    int64  `gorm:"index"`
	IP        string `gorm:"size:50"`
	UserAgent string `gorm:"size:500"`
}

type Comment struct {
	gorm.Model
	PostSlug  string `gorm:"size:200;index"`
	UserID    int64  `gorm:"index"`
	Nickname  string `gorm:"size:50"`
	Content   string `gorm:"type:text"`
	IP        string `gorm:"size:50;index"`
	DeviceID  string `gorm:"size:64;index"`
	UserAgent string `gorm:"size:500"`
	Status    int    `gorm:"default:1;index"` // 1=正常 0=待审核 -1=违规
	IsFake    bool   `gorm:"default:false"`   // 是否为假数据
}

type SensitiveWord struct {
	gorm.Model
	Word  string `gorm:"size:100;uniqueIndex"`
	Level int    `gorm:"default:1"`
}

// Category 文章分类表
type Category struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;comment:分类名称"`
	Slug        string `gorm:"size:100;not null;uniqueIndex;comment:分类别名"`
	Description string `gorm:"size:500;comment:分类描述"`
	Sort        int    `gorm:"default:0;comment:排序权重"`
}

// Tag 文章标签表
type Tag struct {
	gorm.Model
	Name  string `gorm:"size:100;not null;comment:标签名称"`
	Slug  string `gorm:"size:100;not null;uniqueIndex;comment:标签别名"`
	Color string `gorm:"size:20;default:'';comment:标签颜色"`
}

// PostTag 文章标签关联表
type PostTag struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"not null;comment:文章ID;index:idx_tag_id"`
	TagID     uint      `gorm:"not null;comment:标签ID;index:idx_tag_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// PostCategory 文章分类关联表
type PostCategory struct {
	ID         uint      `gorm:"primaryKey"`
	PostID     uint      `gorm:"not null;comment:文章ID;index:idx_pc_post"`
	CategoryID uint      `gorm:"not null;comment:分类ID;index:idx_pc_cat"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// BlogUser 博客用户表（C端用户）- 使用users表
type BlogUser struct {
	ID           uint           `gorm:"primaryKey"`
	Username     string         `gorm:"column:username;size:50;not null;uniqueIndex;comment:用户名"`
	Email        string         `gorm:"column:email;size:100;not null;uniqueIndex;comment:邮箱"`
	PasswordHash string         `gorm:"column:password;size:255;not null;comment:密码哈希"`
	Nickname     string         `gorm:"column:nickname;size:100;default:'';comment:昵称"`
	AvatarURL    string         `gorm:"column:avatar;size:500;default:'';comment:头像"`
	Status       int            `gorm:"column:status;default:1;comment:状态 1=正常 0=禁用"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (BlogUser) TableName() string {
	return "users"
}

// PostLike 文章点赞表
type PostLike struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"not null;uniqueIndex:uk_post_user;comment:文章ID"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_post_user;comment:用户ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// PostFavorite 文章收藏表
type PostFavorite struct {
	ID        uint      `gorm:"primaryKey"`
	PostID    uint      `gorm:"not null;uniqueIndex:uk_post_user;comment:文章ID"`
	UserID    uint      `gorm:"not null;uniqueIndex:uk_post_user;comment:用户ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

// AdminLog 管理后台操作日志
type AdminLog struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"not null;comment:管理员ID"`
	Username   string    `gorm:"size:50;not null;comment:管理员用户名"`
	Action     string    `gorm:"size:100;not null;comment:操作类型"`
	TargetType string    `gorm:"size:50;not null;comment:目标类型"`
	TargetID   *uint     `gorm:"comment:目标ID"`
	TargetName string    `gorm:"size:255;default:'';comment:目标名称"`
	Details    string    `gorm:"type:text;comment:操作详情"`
	IP         string    `gorm:"size:50;default:'';comment:IP地址"`
	UserAgent  string    `gorm:"size:500;default:'';comment:User-Agent"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// AdminRole 管理员角色
type AdminRole struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:50;not null;comment:角色名称"`
	Permissions string    `gorm:"type:text;comment:权限列表JSON"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
