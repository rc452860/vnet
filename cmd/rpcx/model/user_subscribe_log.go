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

type UserSubscribeLog struct {
	ID            int         `gorm:"column:id;primary_key" json:"id"`
	Sid           null.Int    `gorm:"column:sid" json:"sid"`
	RequestIP     null.String `gorm:"column:request_ip" json:"request_ip"`
	RequestTime   null.Time   `gorm:"column:request_time" json:"request_time"`
	RequestHeader null.String `gorm:"column:request_header" json:"request_header"`
}

// TableName sets the insert table name for this struct type
func (u *UserSubscribeLog) TableName() string {
	return "user_subscribe_log"
}
