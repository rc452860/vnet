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

type UserBanLog struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Minutes   int       `gorm:"column:minutes" json:"minutes"`
	Desc      string    `gorm:"column:desc" json:"desc"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserBanLog) TableName() string {
	return "user_ban_log"
}
