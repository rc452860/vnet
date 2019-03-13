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

type SsGroupNode struct {
	ID      int `gorm:"column:id;primary_key" json:"id"`
	GroupID int `gorm:"column:group_id" json:"group_id"`
	NodeID  int `gorm:"column:node_id" json:"node_id"`
}

// TableName sets the insert table name for this struct type
func (s *SsGroupNode) TableName() string {
	return "ss_group_node"
}
