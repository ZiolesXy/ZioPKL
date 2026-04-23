package main

import (
	"errors"
	"log"
	"mime"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gorm.io/gorm"

	"server/internal/database"
	"server/internal/domain/models"
	"server/internal/helper"
	"server/internal/storage"
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

	minioStorage, err := storage.NewMinIO(cfg)
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

	profileSelections, err := seedResolveProfileSelections(users)
	if err != nil {
		log.Fatal(err)
	}

	for index, user := range users {
		passwordHash, err := helper.HashPassword(cfg.SeedUserPassword)
		if err != nil {
			log.Fatal(err)
		}

		var existing models.User
		if err := db.Where("email = ?", user.Email).First(&existing).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Fatal(err)
		}

		profileURL, err := seedUploadProfile(minioStorage, user.Email, profileSelections[index], existing.ProfileURL)
		if err != nil {
			log.Fatal(err)
		}

		record := models.User{
			Email:        user.Email,
			PasswordHash: passwordHash,
			Role:         user.Role,
			ProfileURL:   profileURL,
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

func seedResolveProfileSelections(users []struct {
	Email string
	Role  string
}) ([]string, error) {
	files, err := seedListProfileFiles()
	if err != nil {
		return nil, err
	}

	selections := make([]string, len(users))
	if len(files) == 0 {
		return selections, nil
	}

	defaultFile := seedFindDefaultProfile(files)
	if defaultFile != "" {
		for index := range selections {
			selections[index] = defaultFile
		}
		return selections, nil
	}

	for index := range selections {
		if index < len(files) {
			selections[index] = files[index]
		}
	}

	return selections, nil
}

func seedListProfileFiles() ([]string, error) {
	profileDir := filepath.Join("seeders", "profile")
	entries, err := os.ReadDir(profileDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(profileDir, entry.Name())
		if !seedIsImageFile(path) {
			continue
		}

		files = append(files, path)
	}

	slices.Sort(files)
	return files, nil
}

func seedFindDefaultProfile(files []string) string {
	for _, file := range files {
		name := strings.ToLower(strings.TrimSpace(strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))))
		if strings.HasPrefix(name, "default") {
			return file
		}
	}
	return ""
}

func seedUploadProfile(minioStorage *storage.MinIOStorage, email string, sourcePath string, existingURL *string) (*string, error) {
	if strings.TrimSpace(sourcePath) == "" {
		if existingURL != nil {
			oldObject := minioStorage.ExtractProfileObjectName(*existingURL)
			if err := minioStorage.RemoveProfileObject(oldObject); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}

	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(sourcePath))
	objectName := "seed/users/" + seedSanitizeEmail(email) + ext
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if err := minioStorage.UploadFile(minioStorage.ProfileBucket(), objectName, file, info.Size(), contentType); err != nil {
		return nil, err
	}

	if existingURL != nil {
		oldObject := minioStorage.ExtractProfileObjectName(*existingURL)
		if oldObject != "" && oldObject != objectName {
			if err := minioStorage.RemoveProfileObject(oldObject); err != nil {
				return nil, err
			}
		}
	}

	return &objectName, nil
}

func seedSanitizeEmail(email string) string {
	replacer := strings.NewReplacer("@", "-at-", ".", "-", "_", "-", "+", "-plus-")
	return replacer.Replace(strings.ToLower(strings.TrimSpace(email)))
}

func seedIsImageFile(path string) bool {
	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(path)))
	return strings.HasPrefix(contentType, "image/")
}
