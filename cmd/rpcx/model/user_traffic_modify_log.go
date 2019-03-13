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

type UserTrafficModifyLog struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	OrderID   int       `gorm:"column:order_id" json:"order_id"`
	Before    int64     `gorm:"column:before" json:"before"`
	After     int64     `gorm:"column:after" json:"after"`
	Desc      string    `gorm:"column:desc" json:"desc"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserTrafficModifyLog) TableName() string {
	return "user_traffic_modify_log"
}
