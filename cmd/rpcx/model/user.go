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

type User struct {
	ID                int         `gorm:"column:id;primary_key" json:"id"`
	Username          string      `gorm:"column:username" json:"username"`
	Password          string      `gorm:"column:password" json:"password"`
	Port              int         `gorm:"column:port" json:"port"`
	Passwd            string      `gorm:"column:passwd" json:"passwd"`
	VmessID           string      `gorm:"column:vmess_id" json:"vmess_id"`
	TransferEnable    int64       `gorm:"column:transfer_enable" json:"transfer_enable"`
	U                 int64       `gorm:"column:u" json:"u"`
	D                 int64       `gorm:"column:d" json:"d"`
	T                 int         `gorm:"column:t" json:"t"`
	Enable            int         `gorm:"column:enable" json:"enable"`
	Method            string      `gorm:"column:method" json:"method"`
	Protocol          string      `gorm:"column:protocol" json:"protocol"`
	ProtocolParam     null.String `gorm:"column:protocol_param" json:"protocol_param"`
	Obfs              string      `gorm:"column:obfs" json:"obfs"`
	ObfsParam         null.String `gorm:"column:obfs_param" json:"obfs_param"`
	SpeedLimitPerCon  int64       `gorm:"column:speed_limit_per_con" json:"speed_limit_per_con"`
	SpeedLimitPerUser int64       `gorm:"column:speed_limit_per_user" json:"speed_limit_per_user"`
	Gender            int         `gorm:"column:gender" json:"gender"`
	Wechat            null.String `gorm:"column:wechat" json:"wechat"`
	Qq                null.String `gorm:"column:qq" json:"qq"`
	Usage             string      `gorm:"column:usage" json:"usage"`
	PayWay            int         `gorm:"column:pay_way" json:"pay_way"`
	Balance           int         `gorm:"column:balance" json:"balance"`
	Score             int         `gorm:"column:score" json:"score"`
	EnableTime        null.Time   `gorm:"column:enable_time" json:"enable_time"`
	ExpireTime        time.Time   `gorm:"column:expire_time" json:"expire_time"`
	BanTime           int         `gorm:"column:ban_time" json:"ban_time"`
	Remark            null.String `gorm:"column:remark" json:"remark"`
	Level             int         `gorm:"column:level" json:"level"`
	IsAdmin           int         `gorm:"column:is_admin" json:"is_admin"`
	RegIP             string      `gorm:"column:reg_ip" json:"reg_ip"`
	LastLogin         int         `gorm:"column:last_login" json:"last_login"`
	ReferralUID       int         `gorm:"column:referral_uid" json:"referral_uid"`
	TrafficResetDay   int         `gorm:"column:traffic_reset_day" json:"traffic_reset_day"`
	Status            int         `gorm:"column:status" json:"status"`
	RememberToken     null.String `gorm:"column:remember_token" json:"remember_token"`
	CreatedAt         null.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         null.Time   `gorm:"column:updated_at" json:"updated_at"`
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "user"
}
