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

type Country struct {
	ID          int    `gorm:"column:id;primary_key" json:"id"`
	CountryName string `gorm:"column:country_name" json:"country_name"`
	CountryCode string `gorm:"column:country_code" json:"country_code"`
}

// TableName sets the insert table name for this struct type
func (c *Country) TableName() string {
	return "country"
}
