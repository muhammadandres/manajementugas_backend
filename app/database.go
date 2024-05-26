package app

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"manajemen_tugas_master/model/domain"
	"os"
)

func ConnectDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %v", err)
	}

	log.Println("Running migrations")
	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Task{},
	); err != nil {
		return nil, err
	}

	return db, err
}
