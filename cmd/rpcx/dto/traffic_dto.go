package dto

type TrafficDto struct {
	UserID int   `gorm:"column:user_id" json:"user_id"`
	U      int64 `gorm:"column:u" json:"u"`
	D      int64 `gorm:"column:d" json:"d"`
	Port   int   `gorm:"column:port" json:"port"`
}
