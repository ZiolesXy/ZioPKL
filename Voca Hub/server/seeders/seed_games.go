package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"log"
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

var defaultCategories = []string{
	"Action",
	"Adventure",
	"Arcade",
	"Puzzle",
	"Strategy",
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
		&models.Friend{},
		&models.Message{},
		&models.Game{},
		&models.Category{},
		&models.Post{},
	); err != nil {
		log.Fatal(err)
	}

	developerID, err := resolveSeederDeveloperID(db)
	if err != nil {
		log.Fatal(err)
	}

	categories, err := ensureCategories(db)
	if err != nil {
		log.Fatal(err)
	}
	if len(categories) < 5 {
		log.Fatal("at least 5 categories are required")
	}

	gameDir, err := resolveGameSeedDir()
	if err != nil {
		log.Fatal(err)
	}

	zipFiles, err := filepath.Glob(filepath.Join(gameDir, "*.zip"))
	if err != nil {
		log.Fatal(err)
	}
	if len(zipFiles) == 0 {
		log.Fatalf("no zip files found in %s", gameDir)
	}
	slices.Sort(zipFiles)

	for index, zipPath := range zipFiles {
		title := titleFromFilename(filepath.Base(zipPath))
		description := fmt.Sprintf("Placeholder description for %s.", title)

		objectPrefix, err := uploadSeedGameArchive(minioStorage, developerID, title, zipPath)
		if err != nil {
			log.Fatalf("failed seeding %s: %v", filepath.Base(zipPath), err)
		}

		selectedCategories := []models.Category{
			categories[index%len(categories)],
			categories[(index+1)%len(categories)],
		}

		game := models.Game{
			Title:       title,
			Description: description,
			FileURL:     objectPrefix,
			DeveloperID: developerID,
			Status:      "approved",
			Categories:  selectedCategories,
		}

		if err := db.Where("developer_id = ? AND title = ?", developerID, title).Assign(game).FirstOrCreate(&game).Error; err != nil {
			log.Fatalf("failed persisting game %s: %v", title, err)
		}

		log.Printf("seeded game: %s", title)
	}
}

func ensureCategories(db *gorm.DB) ([]models.Category, error) {
	for _, name := range defaultCategories {
		category := models.Category{Name: name}
		if err := db.Where("name = ?", name).FirstOrCreate(&category).Error; err != nil {
			return nil, err
		}
	}

	var categories []models.Category
	if err := db.Order("name asc").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func resolveSeederDeveloperID(db *gorm.DB) (uint, error) {
	var user models.User

	for _, role := range []string{"DEVELOPER", "ADMIN", "USER"} {
		if err := db.Where("role = ?", role).First(&user).Error; err == nil {
			return user.ID, nil
		}
	}

	if err := db.First(&user).Error; err == nil {
		return user.ID, nil
	}

	return 0, errors.New("no users found; run user seeder first")
}

func resolveGameSeedDir() (string, error) {
	candidates := []string{
		filepath.Join("seeders", "games"),
		filepath.Join("seeders", "Games"),
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", errors.New("game seed directory not found; expected seeders/games")
}

func titleFromFilename(name string) string {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	base = strings.ReplaceAll(base, "_", " ")
	base = strings.ReplaceAll(base, "-", " ")
	return strings.TrimSpace(strings.Join(strings.Fields(base), " "))
}

func uploadSeedGameArchive(minioStorage *storage.MinIOStorage, developerID uint, title string, zipPath string) (string, error) {
	tempRoot, err := os.MkdirTemp("", "seed-game-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tempRoot)

	extractDir := filepath.Join(tempRoot, "extracted")
	if err := validateAndExtractSeedArchive(zipPath, extractDir); err != nil {
		return "", err
	}

	objectPrefix := fmt.Sprintf("%d/%s", developerID, helper.Slugify(title))
	if err := minioStorage.UploadDirectory(objectPrefix, extractDir); err != nil {
		return "", err
	}

	return objectPrefix, nil
}

func validateAndExtractSeedArchive(zipPath string, extractDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return errors.New("invalid zip file")
	}
	defer reader.Close()

	if len(reader.File) == 0 {
		return errors.New("zip file is empty")
	}

	if err := helper.ExtractZip(zipPath, extractDir); err != nil {
		return err
	}

	indexInfo, err := os.Stat(filepath.Join(extractDir, "index.html"))
	if err != nil || indexInfo.IsDir() {
		return errors.New("zip must contain index.html in root")
	}

	return nil
}
