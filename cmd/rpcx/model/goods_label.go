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

type GoodsLabel struct {
	ID      int `gorm:"column:id;primary_key" json:"id"`
	GoodsID int `gorm:"column:goods_id" json:"goods_id"`
	LabelID int `gorm:"column:label_id" json:"label_id"`
}

// TableName sets the insert table name for this struct type
func (g *GoodsLabel) TableName() string {
	return "goods_label"
}
