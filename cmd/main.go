package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type HelloWorld struct {
	ID    uint
	Hello string
	World string
}

func readFirstRecord(db *gorm.DB) HelloWorld {
	var helloWorld HelloWorld
	db.First(&helloWorld)
	return helloWorld
}

func main() {
	dsn := fmt.Sprintf("root:%s@tcp(localhost:3306)/landd?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("LANDD_MYSQL_PASSWORD"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect database:", err)
	}

	data := readFirstRecord(db)

	s, _ := json.Marshal(data)

	router := gin.Default()

	router.GET("/hello-world", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": s,
		})
	})

	router.Run(":8080")
}
