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
	IsPinned   bool       `gorm:"default:false;comment:是否置顶"`
	IsFeatured bool       `gorm:"default:false;comment:是否推荐"`
	PinnedAt   *time.Time `gorm:"comment:置顶时间"`
	ScheduledAt *time.Time `gorm:"comment:定时发布时间"`
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
