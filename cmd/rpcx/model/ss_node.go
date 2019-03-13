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

type SsNode struct {
	ID             int         `gorm:"column:id;primary_key" json:"id"`
	Type           int         `gorm:"column:type" json:"type"`
	Name           string      `gorm:"column:name" json:"name"`
	GroupID        int         `gorm:"column:group_id" json:"group_id"`
	CountryCode    null.String `gorm:"column:country_code" json:"country_code"`
	Server         null.String `gorm:"column:server" json:"server"`
	IP             null.String `gorm:"column:ip" json:"ip"`
	Ipv6           null.String `gorm:"column:ipv6" json:"ipv6"`
	Desc           null.String `gorm:"column:desc" json:"desc"`
	Method         string      `gorm:"column:method" json:"method"`
	Protocol       string      `gorm:"column:protocol" json:"protocol"`
	ProtocolParam  null.String `gorm:"column:protocol_param" json:"protocol_param"`
	Obfs           string      `gorm:"column:obfs" json:"obfs"`
	ObfsParam      null.String `gorm:"column:obfs_param" json:"obfs_param"`
	TrafficRate    float32     `gorm:"column:traffic_rate" json:"traffic_rate"`
	Bandwidth      int         `gorm:"column:bandwidth" json:"bandwidth"`
	Traffic        int64       `gorm:"column:traffic" json:"traffic"`
	MonitorURL     null.String `gorm:"column:monitor_url" json:"monitor_url"`
	IsSubscribe    null.Int    `gorm:"column:is_subscribe" json:"is_subscribe"`
	SSHPort        int         `gorm:"column:ssh_port" json:"ssh_port"`
	IsTcpCheck     int         `gorm:"column:is_tcp_check" json:"is_tcp_check"`
	Icmp           int         `gorm:"column:icmp" json:"icmp"`
	Tcp            int         `gorm:"column:tcp" json:"tcp"`
	Udp            int         `gorm:"column:udp" json:"udp"`
	Compatible     null.Int    `gorm:"column:compatible" json:"compatible"`
	Single         null.Int    `gorm:"column:single" json:"single"`
	SingleForce    null.Int    `gorm:"column:single_force" json:"single_force"`
	SinglePort     null.String `gorm:"column:single_port" json:"single_port"`
	SinglePasswd   null.String `gorm:"column:single_passwd" json:"single_passwd"`
	SingleMethod   null.String `gorm:"column:single_method" json:"single_method"`
	SingleProtocol string      `gorm:"column:single_protocol" json:"single_protocol"`
	SingleObfs     string      `gorm:"column:single_obfs" json:"single_obfs"`
	Sort           int         `gorm:"column:sort" json:"sort"`
	Status         int         `gorm:"column:status" json:"status"`
	V2AlterID      int         `gorm:"column:v2_alter_id" json:"v2_alter_id"`
	V2Port         int         `gorm:"column:v2_port" json:"v2_port"`
	V2Net          string      `gorm:"column:v2_net" json:"v2_net"`
	V2Type         string      `gorm:"column:v2_type" json:"v2_type"`
	V2Host         string      `gorm:"column:v2_host" json:"v2_host"`
	V2Path         string      `gorm:"column:v2_path" json:"v2_path"`
	V2TLS          int         `gorm:"column:v2_tls" json:"v2_tls"`
	CreatedAt      time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"column:updated_at" json:"updated_at"`
	Token          string      `gorm:"column:token" json:"token"`
}

// TableName sets the insert table name for this struct type
func (s *SsNode) TableName() string {
	return "ss_node"
}
