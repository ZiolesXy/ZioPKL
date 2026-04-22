package main

import (
	"errors"
	"log"

	"gorm.io/gorm"

	"server/internal/database"
	"server/internal/domain/models"
	"server/internal/helper"
)

var seedDefaultDifficulties = []string{
	"Mudah",
	"Sedang",
	"Sulit",
}

func main() {
	cfg := helper.LoadConfig()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Difficulty{},
	); err != nil {
		log.Fatal(err)
	}

	users := []struct {
		Email string
		Role  string
	}{
		{Email: "eaglegaming3605@gmail.com", Role: "ADMIN"},
		{Email: "developer@example.com", Role: "DEVELOPER"},
		{Email: "pashaprabasakti@gmail.com", Role: "USER"},
		{Email: "abrilliantp738@gmail.com", Role: "USER"},
	}

	for _, user := range users {
		passwordHash, err := helper.HashPassword(cfg.SeedUserPassword)
		if err != nil {
			log.Fatal(err)
		}

		record := models.User{
			Email:        user.Email,
			PasswordHash: passwordHash,
			Role:         user.Role,
		}
		if err := db.Where("email = ?", user.Email).Assign(record).FirstOrCreate(&record).Error; err != nil {
			log.Fatal(err)
		}
	}

	if _, err := seedEnsureDifficulties(db); err != nil {
		log.Fatal(err)
	}
}

func seedEnsureDifficulties(db *gorm.DB) ([]models.Difficulty, error) {
	for index, name := range seedDefaultDifficulties {
		difficulty := models.Difficulty{
			ID:   uint(index + 1),
			Name: name,
		}
		if err := db.Where("id = ?", difficulty.ID).Assign(models.Difficulty{Name: difficulty.Name}).FirstOrCreate(&difficulty).Error; err != nil {
			return nil, err
		}
	}

	var difficulties []models.Difficulty
	if err := db.Order("id asc").Find(&difficulties).Error; err != nil {
		return nil, err
	}
	if len(difficulties) != len(seedDefaultDifficulties) {
		return nil, errors.New("exactly 3 difficulties are required")
	}
	return difficulties, nil
}
