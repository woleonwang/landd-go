package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO put db config in separate config file
const DB_USERNAME = "landd"
const DB_PASSWORD = ""
const DB_NAME = "landd"
const DB_HOST = ""
const DB_PORT = "3306"

var db *gorm.DB

func Init() {
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" +
		DB_NAME + "?" + "parseTime=true&loc=Local"

	mysqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error connecting to DB")
	}
	db = mysqlDB
}

func GetDB() *gorm.DB {
	return db
}
