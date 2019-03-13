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

type Verify struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Type      int       `gorm:"column:type" json:"type"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Token     string    `gorm:"column:token" json:"token"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *Verify) TableName() string {
	return "verify"
}
