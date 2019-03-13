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

type UserLoginLog struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	IP        string    `gorm:"column:ip" json:"ip"`
	Country   string    `gorm:"column:country" json:"country"`
	Province  string    `gorm:"column:province" json:"province"`
	City      string    `gorm:"column:city" json:"city"`
	County    string    `gorm:"column:county" json:"county"`
	Isp       string    `gorm:"column:isp" json:"isp"`
	Area      string    `gorm:"column:area" json:"area"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserLoginLog) TableName() string {
	return "user_login_log"
}
