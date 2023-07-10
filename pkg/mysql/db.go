package mysql

import (
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO put db config in separate config file
const DB_USERNAME = "root"
const DB_NAME = "landd"
const DB_HOST = "localhost"
const DB_PORT = "3306"

var DB_PASSWORD = ""
var db *gorm.DB

func Init() {
	if DB_PASSWORD == "" {
		DB_PASSWORD = os.Getenv("LANDD_MYSQL_PASSWORD")
	}

	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@(" + DB_HOST + ":" + DB_PORT + ")/" +
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
