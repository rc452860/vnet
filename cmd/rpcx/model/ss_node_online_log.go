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

type SsNodeOnlineLog struct {
	ID         int `gorm:"column:id;primary_key" json:"id"`
	NodeID     int `gorm:"column:node_id" json:"node_id"`
	OnlineUser int `gorm:"column:online_user" json:"online_user"`
	LogTime    int `gorm:"column:log_time" json:"log_time"`
}

// TableName sets the insert table name for this struct type
func (s *SsNodeOnlineLog) TableName() string {
	return "ss_node_online_log"
}
