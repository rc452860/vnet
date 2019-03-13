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

type SsNodeIP struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	NodeID    int         `gorm:"column:node_id" json:"node_id"`
	Port      int         `gorm:"column:port" json:"port"`
	Type      string      `gorm:"column:type" json:"type"`
	IP        null.String `gorm:"column:ip" json:"ip"`
	CreatedAt int         `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (s *SsNodeIP) TableName() string {
	return "ss_node_ip"
}
