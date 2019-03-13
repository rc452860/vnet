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

type UserBalanceLog struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	UserID    int         `gorm:"column:user_id" json:"user_id"`
	OrderID   int         `gorm:"column:order_id" json:"order_id"`
	Before    int         `gorm:"column:before" json:"before"`
	After     int         `gorm:"column:after" json:"after"`
	Amount    int         `gorm:"column:amount" json:"amount"`
	Desc      null.String `gorm:"column:desc" json:"desc"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserBalanceLog) TableName() string {
	return "user_balance_log"
}
