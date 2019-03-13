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

type PaymentCallback struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	ClientID  null.String `gorm:"column:client_id" json:"client_id"`
	YzID      null.String `gorm:"column:yz_id" json:"yz_id"`
	KdtID     null.String `gorm:"column:kdt_id" json:"kdt_id"`
	KdtName   null.String `gorm:"column:kdt_name" json:"kdt_name"`
	Mode      null.Int    `gorm:"column:mode" json:"mode"`
	Msg       null.String `gorm:"column:msg" json:"msg"`
	SendCount null.Int    `gorm:"column:sendCount" json:"sendCount"`
	Sign      null.String `gorm:"column:sign" json:"sign"`
	Status    null.String `gorm:"column:status" json:"status"`
	Test      null.Int    `gorm:"column:test" json:"test"`
	Type      null.String `gorm:"column:type" json:"type"`
	Version   null.String `gorm:"column:version" json:"version"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (p *PaymentCallback) TableName() string {
	return "payment_callback"
}
