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

type SensitiveWord struct {
	ID    int    `gorm:"column:id;primary_key" json:"id"`
	Words string `gorm:"column:words" json:"words"`
}

// TableName sets the insert table name for this struct type
func (s *SensitiveWord) TableName() string {
	return "sensitive_words"
}
