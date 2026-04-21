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
	clerkClient := helper.NewClerkClient(cfg)

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
	}

	for _, user := range users {
		clerkID, err := clerkClient.FetchUserIDByEmail(user.Email)
		if err != nil {
			log.Fatal(err)
		}
		username, err := clerkClient.FetchUsername(clerkID)
		if err != nil {
			log.Fatal(err)
		}

		record := models.User{
			ClerkID: clerkID,
			Email:   user.Email,
			Username: func() *string {
				if username == "" {
					return nil
				}
				return &username
			}(),
			Role: user.Role,
		}
		if err := db.Where("clerk_id = ?", clerkID).Assign(record).FirstOrCreate(&record).Error; err != nil {
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
