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

type UserLabel struct {
	ID      int `gorm:"column:id;primary_key" json:"id"`
	UserID  int `gorm:"column:user_id" json:"user_id"`
	LabelID int `gorm:"column:label_id" json:"label_id"`
}

// TableName sets the insert table name for this struct type
func (u *UserLabel) TableName() string {
	return "user_label"
}
