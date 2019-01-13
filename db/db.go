package db

import (
	"fmt"
	"strings"

	"github.com/rc452860/vnet/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	Id             int64  `gorm:"column:id;primary_key"`
	Port           int    `gorm:"column:port"`
	Enable         bool   `gorm:"column:enable"`
	Update         int64  `gorm:"column:u"`
	Download       int64  `gorm:"column:d"`
	TransferEnable int64  `gorm:"column:transfer_enable"`
	Method         string `gorm:"column:method"`
	Password       string `gorm:"column:passwd"`
	Limit          uint64 `gorm:"column:speed_limit_per_user"`
}

func (this User) TableName() string {
	return "user"
}

func Connect() (*gorm.DB, error) {

	var (
		username string
		password string
		database string
		host     string
		port     string
	)
	pattern := "username:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local"

	c := config.CurrentConfig()
	if c == nil {
		return nil, fmt.Errorf("config dbconfig is nil")
	}
	// log.Info("dbConfig: %s", c.DbConfig)
	username = c.DbConfig.User
	password = c.DbConfig.Passwd
	database = c.DbConfig.Database
	host = c.DbConfig.Host
	port = c.DbConfig.Port
	replacer := strings.NewReplacer(
		"username",
		username,
		"password",
		password,
		"host",
		host,
		"port",
		port,
		"database",
		database,
	)
	conntionStr := replacer.Replace(pattern)
	// log.Info("connection string: %s", conntionStr)
	db, err := gorm.Open("mysql", conntionStr)
	return db, err
}

func GetEnableUser() ([]User, error) {
	db, err := Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var userList []User
	db.Where("port != 0 AND enable = 1").Find(&userList)
	return userList, nil
}
