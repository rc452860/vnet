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

type UserScoreLog struct {
	ID        int         `gorm:"column:id;primary_key" json:"id"`
	UserID    int         `gorm:"column:user_id" json:"user_id"`
	Before    int         `gorm:"column:before" json:"before"`
	After     int         `gorm:"column:after" json:"after"`
	Score     int         `gorm:"column:score" json:"score"`
	Desc      null.String `gorm:"column:desc" json:"desc"`
	CreatedAt null.Time   `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (u *UserScoreLog) TableName() string {
	return "user_score_log"
}
