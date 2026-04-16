package dao

import (
	"gorm.io/gorm"
)

// DAO 数据访问对象基类
type DAO struct {
	db *gorm.DB
}

// NewDAO 创建 DAO 实例
func NewDAO(db *gorm.DB) *DAO {
	return &DAO{db: db}
}

// GetDB 获取数据库连接
func (d *DAO) GetDB() *gorm.DB {
	return d.db
}
