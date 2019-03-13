package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
)

type Article struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	Title     string      `gorm:"column:title" json:"title"`
	Author    null.String `gorm:"column:author" json:"author"`
	Summary   null.String `gorm:"column:summary" json:"summary"`
	Content   null.String `gorm:"column:content" json:"content"`
	Type      null.Int    `gorm:"column:type" json:"type"`
	IsDel     int         `gorm:"column:is_del" json:"is_del"`
	Sort      int         `gorm:"column:sort" json:"sort"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (a *Article) TableName() string {
	return "article"
}
