package main

import (
	"log"

	"server/internal/database"
	"server/internal/domain/models"
	"server/internal/helper"
)

func main() {
	cfg := helper.LoadConfig()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Friend{},
		&models.Message{},
		&models.Game{},
		&models.Category{},
		&models.Post{},
	); err != nil {
		log.Fatal(err)
	}
}
