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

type SsNodeLabel struct {
	ID      int `gorm:"column:id;primary_key" json:"id"`
	NodeID  int `gorm:"column:node_id" json:"node_id"`
	LabelID int `gorm:"column:label_id" json:"label_id"`
}

// TableName sets the insert table name for this struct type
func (s *SsNodeLabel) TableName() string {
	return "ss_node_label"
}
