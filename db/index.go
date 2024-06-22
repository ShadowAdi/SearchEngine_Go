package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DBConn *gorm.DB

func InitDB() {

	dbUrl := os.Getenv("DATABASE_URL")
	var err error
	DBConn, err = gorm.Open(postgres.Open(dbUrl))
	if err != nil {
		fmt.Println("Failed to connect to a database")
		panic("Database Error")
	}

	if err := DBConn.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		fmt.Println("Error in database:", err)
		panic(err)
	}

	if err := DBConn.AutoMigrate(&User{}, &SearchSettings{}, &CrawledUrl{}); err != nil {
		panic("Migration Failed")
	}

}

func GetDB() *gorm.DB {
	return DBConn
}
