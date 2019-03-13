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

type Marketing struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	Type      int         `gorm:"column:type" json:"type"`
	Receiver  string      `gorm:"column:receiver" json:"receiver"`
	Title     string      `gorm:"column:title" json:"title"`
	Content   string      `gorm:"column:content" json:"content"`
	Error     null.String `gorm:"column:error" json:"error"`
	Status    int         `gorm:"column:status" json:"status"`
	CreatedAt time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (m *Marketing) TableName() string {
	return "marketing"
}
