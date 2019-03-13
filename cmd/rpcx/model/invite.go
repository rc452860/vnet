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

type Invite struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	UID       int       `gorm:"column:uid" json:"uid"`
	Fuid      int       `gorm:"column:fuid" json:"fuid"`
	Code      string    `gorm:"column:code" json:"code"`
	Status    int       `gorm:"column:status" json:"status"`
	Dateline  null.Time `gorm:"column:dateline" json:"dateline"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt null.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (i *Invite) TableName() string {
	return "invite"
}
