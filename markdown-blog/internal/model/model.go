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
}

type Visit struct {
	gorm.Model
	PostSlug string `gorm:"index"`
}
