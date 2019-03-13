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

type Config struct {
	ID    int         `gorm:"column:id;primary_key" json:"id"`
	Name  string      `gorm:"column:name" json:"name"`
	Value null.String `gorm:"column:value" json:"value"`
}

// TableName sets the insert table name for this struct type
func (c *Config) TableName() string {
	return "config"
}
