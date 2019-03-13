package service

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var ds string

// 初始化数据库
func InitDS(datasource string) error {
	_, err := gorm.Open("mysql", datasource)
	if err != nil {
		return err
	}
	ds = datasource
	return nil
}

// 获得数据库连接
func GetDB() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", ds)
	if err != nil {
		return nil, err
	}
	return db, nil
}
