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

type ReferralLog struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	RefUserID int       `gorm:"column:ref_user_id" json:"ref_user_id"`
	OrderID   int       `gorm:"column:order_id" json:"order_id"`
	Amount    int       `gorm:"column:amount" json:"amount"`
	RefAmount int       `gorm:"column:ref_amount" json:"ref_amount"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (r *ReferralLog) TableName() string {
	return "referral_log"
}
