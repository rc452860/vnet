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

type SsConfig struct {
	ID        int    `gorm:"column:id;primary_key" json:"id"`
	Name      string `gorm:"column:name" json:"name"`
	Type      int    `gorm:"column:type" json:"type"`
	IsDefault int    `gorm:"column:is_default" json:"is_default"`
	Sort      int    `gorm:"column:sort" json:"sort"`
}

// TableName sets the insert table name for this struct type
func (s *SsConfig) TableName() string {
	return "ss_config"
}
