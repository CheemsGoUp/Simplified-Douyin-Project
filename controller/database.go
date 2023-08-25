package controller

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDB() (db *gorm.DB, err error) {
	dsn := "root:sqlofmine@(127.0.0.1:3306)/SIMPLIFIED_DOUYIN_PROJECT?charset=utf8mb4&parseTime=True&loc=Local"

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return db, err
	}

	return db, nil
}
