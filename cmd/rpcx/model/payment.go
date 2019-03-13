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

type Payment struct {
	ID         int         `gorm:"column:id;primary_key" json:"id"`
	Sn         null.String `gorm:"column:sn" json:"sn"`
	UserID     int         `gorm:"column:user_id" json:"user_id"`
	Oid        null.Int    `gorm:"column:oid" json:"oid"`
	OrderSn    null.String `gorm:"column:order_sn" json:"order_sn"`
	PayWay     int         `gorm:"column:pay_way" json:"pay_way"`
	Amount     int         `gorm:"column:amount" json:"amount"`
	QrID       int         `gorm:"column:qr_id" json:"qr_id"`
	QrURL      null.String `gorm:"column:qr_url" json:"qr_url"`
	QrCode     null.String `gorm:"column:qr_code" json:"qr_code"`
	QrLocalURL null.String `gorm:"column:qr_local_url" json:"qr_local_url"`
	Status     int         `gorm:"column:status" json:"status"`
	CreatedAt  time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (p *Payment) TableName() string {
	return "payment"
}
