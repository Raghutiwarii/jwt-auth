package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jwt-auth/models"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	dsn := "user=postgres password=root dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	DB = db
	db.AutoMigrate(&models.User{})

	return db, nil

}
