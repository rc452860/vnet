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

type TicketReply struct {
	ID        int       `gorm:"column:id;primary_key" json:"id"`
	TicketID  int       `gorm:"column:ticket_id" json:"ticket_id"`
	UserID    int       `gorm:"column:user_id" json:"user_id"`
	Content   string    `gorm:"column:content" json:"content"`
	CreatedAt null.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName sets the insert table name for this struct type
func (t *TicketReply) TableName() string {
	return "ticket_reply"
}
