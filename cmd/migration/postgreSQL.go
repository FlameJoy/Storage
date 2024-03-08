package main

import (
	"fmt"
	"log"
	"os"
	"testezhik/cmd/api/initialization"
	"testezhik/cmd/api/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	initialization.LoadEnv("../../.env")
}

func main() {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect PostgreSQL DB: " + err.Error())
	}
	fmt.Println("Successfully connected to PostgreSQL DB")
	// Migrate the schemas
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Migration error: %s", err)
		return
	}
	log.Println("Migration successfully completed")
}
