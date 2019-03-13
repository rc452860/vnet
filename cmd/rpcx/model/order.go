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

type Order struct {
}

// TableName sets the insert table name for this struct type
func (o *Order) TableName() string {
	return "order"
}
