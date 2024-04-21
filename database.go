package main

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func connectDB() error {
	dsn := "postgresql://matthieu:matthieu@localhost:5432/godecide?sslmode=disable"
	Database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err)
		os.Exit(2)
	}
	db = Database

	log.Println("Connected Successfully to Database")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	db.AutoMigrate(&User{}, &UserSession{}, &Workspace{})

	return nil
}
