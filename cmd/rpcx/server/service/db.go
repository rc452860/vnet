package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var ds string

// init database
func InitDS(datasource string) error {
	_, err := gorm.Open("mysql", datasource)
	if err != nil {
		return err
	}
	ds = datasource
	return nil
}

// get database connection
func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", ds)
	if err != nil {
		return nil, err
	}
	// db.LogMode(true)
	return db, nil
}
