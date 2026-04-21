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

	if err := db.Migrator().DropTable(
		"game_categories",
		&models.Post{},
		&models.Category{},
		&models.Game{},
		&models.Difficulty{},
		&models.Message{},
		&models.Friend{},
		&models.User{},
	); err != nil {
		log.Fatal(err)
	}
}
