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

type UserSubscribe struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	UserID    int         `gorm:"column:user_id" json:"user_id"`
	Code      null.String `gorm:"column:code" json:"code"`
	Times     int         `gorm:"column:times" json:"times"`
	Status    int         `gorm:"column:status" json:"status"`
	BanTime   int         `gorm:"column:ban_time" json:"ban_time"`
	BanDesc   string      `gorm:"column:ban_desc" json:"ban_desc"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserSubscribe) TableName() string {
	return "user_subscribe"
}
