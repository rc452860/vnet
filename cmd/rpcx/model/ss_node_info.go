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

type SsNodeInfo struct {
	ID      int     `gorm:"column:id;primary_key" json:"id"`
	NodeID  int     `gorm:"column:node_id" json:"node_id"`
	Uptime  float32 `gorm:"column:uptime" json:"uptime"`
	Load    string  `gorm:"column:load" json:"load"`
	LogTime int     `gorm:"column:log_time" json:"log_time"`
}

// TableName sets the insert table name for this struct type
func (s *SsNodeInfo) TableName() string {
	return "ss_node_info"
}
