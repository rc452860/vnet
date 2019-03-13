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

type Good struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	Sku       string      `gorm:"column:sku" json:"sku"`
	Name      string      `gorm:"column:name" json:"name"`
	Logo      string      `gorm:"column:logo" json:"logo"`
	Traffic   int64       `gorm:"column:traffic" json:"traffic"`
	Score     int         `gorm:"column:score" json:"score"`
	Type      int         `gorm:"column:type" json:"type"`
	Price     int         `gorm:"column:price" json:"price"`
	Desc      null.String `gorm:"column:desc" json:"desc"`
	Days      int         `gorm:"column:days" json:"days"`
	Color     string      `gorm:"column:color" json:"color"`
	Sort      int         `gorm:"column:sort" json:"sort"`
	IsLimit   int         `gorm:"column:is_limit" json:"is_limit"`
	IsHot     int         `gorm:"column:is_hot" json:"is_hot"`
	IsDel     int         `gorm:"column:is_del" json:"is_del"`
	Status    int         `gorm:"column:status" json:"status"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (g *Good) TableName() string {
	return "goods"
}
