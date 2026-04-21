package main

import (
	"archive/zip"
	"errors"
	"fmt"
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

var defaultCategories = []string{
	"Action",
	"Adventure",
	"Arcade",
	"Puzzle",
	"Strategy",
}

var defaultDifficulties = []string{
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
		&models.Friend{},
		&models.Message{},
		&models.Difficulty{},
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
	difficulties, err := ensureDifficulties(db)
	if err != nil {
		log.Fatal(err)
	}
	if len(categories) < 5 {
		log.Fatal("at least 5 categories are required")
	}
	if len(difficulties) != len(defaultDifficulties) {
		log.Fatal("exactly 3 difficulties are required")
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

	thumbnailPaths, err := resolveSeedThumbnailPaths(minioStorage, zipFiles)
	if err != nil {
		log.Fatal(err)
	}

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
			ThumbnailPath: func() string {
				if index < len(thumbnailPaths) {
					return thumbnailPaths[index]
				}
				return ""
			}(),
			DeveloperID:  developerID,
			DifficultyID: difficulties[index%len(difficulties)].ID,
			Difficulty:   difficulties[index%len(difficulties)],
			Status:       "approved",
			Categories:   selectedCategories,
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

func ensureDifficulties(db *gorm.DB) ([]models.Difficulty, error) {
	for index, name := range defaultDifficulties {
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
	return difficulties, nil
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

func resolveSeedThumbnailPaths(minioStorage *storage.MinIOStorage, zipFiles []string) ([]string, error) {
	assignments := make([]string, len(zipFiles))

	thumbnailDir := filepath.Join("seeders", "thumbnail")
	info, err := os.Stat(thumbnailDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return assignments, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", thumbnailDir)
	}

	entries, err := os.ReadDir(thumbnailDir)
	if err != nil {
		return nil, err
	}

	var defaultFile string
	var imageFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		fullPath := filepath.Join(thumbnailDir, name)
		if !isSeedImageFile(name) {
			continue
		}

		baseName := strings.TrimSuffix(strings.ToLower(name), strings.ToLower(filepath.Ext(name)))
		if baseName == "default" && defaultFile == "" {
			defaultFile = fullPath
			continue
		}

		imageFiles = append(imageFiles, fullPath)
	}

	if defaultFile != "" {
		objectName, err := uploadSeedThumbnailFile(minioStorage, defaultFile)
		if err != nil {
			return nil, err
		}
		for i := range assignments {
			assignments[i] = objectName
		}
		return assignments, nil
	}

	if len(imageFiles) == 0 {
		return assignments, nil
	}

	slices.Sort(imageFiles)

	uploaded := make([]string, 0, len(imageFiles))
	for _, imageFile := range imageFiles {
		objectName, err := uploadSeedThumbnailFile(minioStorage, imageFile)
		if err != nil {
			return nil, err
		}
		uploaded = append(uploaded, objectName)
	}

	for i := range assignments {
		assignments[i] = uploaded[i%len(uploaded)]
	}

	return assignments, nil
}

func uploadSeedThumbnailFile(minioStorage *storage.MinIOStorage, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	contentType := mime.TypeByExtension(strings.ToLower(filepath.Ext(filePath)))
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("thumbnail %s must be an image", filepath.Base(filePath))
	}

	objectName := filepath.ToSlash(filepath.Join("seed", filepath.Base(filePath)))
	if err := minioStorage.UploadFile(minioStorage.ThumbnailBucket(), objectName, file, info.Size(), contentType); err != nil {
		return "", err
	}

	return objectName, nil
}

func isSeedImageFile(name string) bool {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg":
		return true
	default:
		return false
	}
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
