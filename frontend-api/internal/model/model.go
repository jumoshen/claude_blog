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
	Status     int    `gorm:"default:1;index"` // 0=草稿 1=已发布 2=下架
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
}

type SensitiveWord struct {
	gorm.Model
	Word  string `gorm:"size:100;uniqueIndex"`
	Level int    `gorm:"default:1"`
}
