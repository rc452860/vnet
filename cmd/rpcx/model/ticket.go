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

type Ticket struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Title     string    `gorm:"column:title" json:"title"`
	Content   string    `gorm:"column:content" json:"content"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (t *Ticket) TableName() string {
	return "ticket"
}
