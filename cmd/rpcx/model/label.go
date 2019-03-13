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

type Label struct {
	ID   int    `gorm:"column:id;primary_key" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Sort int    `gorm:"column:sort" json:"sort"`
}

// TableName sets the insert table name for this struct type
func (l *Label) TableName() string {
	return "label"
}
