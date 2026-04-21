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
	clerkClient := helper.NewClerkClient(cfg)

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
}
