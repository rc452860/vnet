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

type EmailLog struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	UserID    int         `gorm:"column:user_id" json:"user_id"`
	Title     null.String `gorm:"column:title" json:"title"`
	Content   null.String `gorm:"column:content" json:"content"`
	Status    int         `gorm:"column:status" json:"status"`
	Error     null.String `gorm:"column:error" json:"error"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	Type      null.Int    `gorm:"column:type" json:"type"`
	Address   null.String `gorm:"column:address" json:"address"`
}

// TableName sets the insert table name for this struct type
func (e *EmailLog) TableName() string {
	return "email_log"
}
