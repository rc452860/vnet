package db

import (
	"strings"

	"github.com/rc452860/vnet/utils"

	"github.com/rc452860/vnet/log"

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
}

func (this User) TableName() string {
	return "user"
}

func Connect() (*gorm.DB, error) {
	logging := log.GetLogger("root")
	defer func() {
		if err := recover(); err != nil {
			logging.Error("%v", err)
		}
	}()

	var (
		username string
		password string
		database string
		host     string
		port     string
	)
	pattern := "username:password@tcp(host:port)/database?charset=utf8&parseTime=True&loc=Local"

	config := utils.ConfigFactory("config.json")

	dbConfig := config.Map.GetConfigMap("dbconfig")
	logging.Info("dbConfig: %s", dbConfig)
	username = dbConfig.GetString("username")
	password = dbConfig.GetString("password")
	database = dbConfig.GetString("database")
	host = dbConfig.GetString("host")
	port = dbConfig.GetString("port")
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
	logging.Info("connection string: %s", conntionStr)
	db, err := gorm.Open("mysql", conntionStr)
	return db, err
}

func GetEnableUser() ([]User, error) {
	db, err := Connect()
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var userList []User
	db.Where("port != 0 AND enable = 1").Find(&userList)
	return userList, nil
}
