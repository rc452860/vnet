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

type Coupon struct {
	ID             int       `gorm:"column:id;primary_key" json:"id"`
	Name           string    `gorm:"column:name" json:"name"`
	Logo           string    `gorm:"column:logo" json:"logo"`
	Sn             string    `gorm:"column:sn" json:"sn"`
	Type           int       `gorm:"column:type" json:"type"`
	Usage          int       `gorm:"column:usage" json:"usage"`
	Amount         int64     `gorm:"column:amount" json:"amount"`
	Discount       float64   `gorm:"column:discount" json:"discount"`
	AvailableStart int       `gorm:"column:available_start" json:"available_start"`
	AvailableEnd   int       `gorm:"column:available_end" json:"available_end"`
	IsDel          int       `gorm:"column:is_del" json:"is_del"`
	Status         int       `gorm:"column:status" json:"status"`
	CreatedAt      null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *Coupon) TableName() string {
	return "coupon"
}
