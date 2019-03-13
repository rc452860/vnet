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

type UserTrafficLog struct {
	ID      int     `gorm:"column:id;primary_key" json:"id"`
	UserID  int     `gorm:"column:user_id" json:"user_id"`
	U       int     `gorm:"column:u" json:"u"`
	D       int     `gorm:"column:d" json:"d"`
	NodeID  int     `gorm:"column:node_id" json:"node_id"`
	Rate    float32 `gorm:"column:rate" json:"rate"`
	Traffic string  `gorm:"column:traffic" json:"traffic"`
	LogTime int     `gorm:"column:log_time" json:"log_time"`
}

// TableName sets the insert table name for this struct type
func (u *UserTrafficLog) TableName() string {
	return "user_traffic_log"
}
