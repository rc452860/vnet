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

type CouponLog struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	CouponID  int       `gorm:"column:coupon_id" json:"coupon_id"`
	GoodsID   int       `gorm:"column:goods_id" json:"goods_id"`
	OrderID   int       `gorm:"column:order_id" json:"order_id"`
	Desc      string    `gorm:"column:desc" json:"desc"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (c *CouponLog) TableName() string {
	return "coupon_log"
}
