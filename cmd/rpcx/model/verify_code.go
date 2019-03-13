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

type VerifyCode struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Code      string    `gorm:"column:code" json:"code"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (v *VerifyCode) TableName() string {
	return "verify_code"
}
