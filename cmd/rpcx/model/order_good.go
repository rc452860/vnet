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

type OrderGood struct {
	ID          int       `gorm:"column:id;primary_key" json:"id"`
	Oid         int       `gorm:"column:oid" json:"oid"`
	OrderSn     string    `gorm:"column:order_sn" json:"order_sn"`
	UserID      int       `gorm:"column:user_id" json:"user_id"`
	GoodsID     int       `gorm:"column:goods_id" json:"goods_id"`
	Num         int       `gorm:"column:num" json:"num"`
	OriginPrice int       `gorm:"column:origin_price" json:"origin_price"`
	Price       int       `gorm:"column:price" json:"price"`
	IsExpire    int       `gorm:"column:is_expire" json:"is_expire"`
	CreatedAt   null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (o *OrderGood) TableName() string {
	return "order_goods"
}
