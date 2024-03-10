package utils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := "user=chatuser password=yoyoyoyoyo dbname=chatapp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestDBConnection() {
	db, err := ConnectDB()
	if err != nil {
		fmt.Println("Failed to connect", err)
		return
	}

	sqlDB, _ := db.DB()

	defer sqlDB.Close()

	fmt.Println("Connected to database successfully!")
}
