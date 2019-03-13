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

type ReferralApply struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Before    int       `gorm:"column:before" json:"before"`
	After     int       `gorm:"column:after" json:"after"`
	Amount    int       `gorm:"column:amount" json:"amount"`
	LinkLogs  string    `gorm:"column:link_logs" json:"link_logs"`
	Status    int       `gorm:"column:status" json:"status"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (r *ReferralApply) TableName() string {
	return "referral_apply"
}
