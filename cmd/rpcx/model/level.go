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

type Level struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Level     int       `gorm:"column:level" json:"level"`
	LevelName string    `gorm:"column:level_name" json:"level_name"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (l *Level) TableName() string {
	return "level"
}
