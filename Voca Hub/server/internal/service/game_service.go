package service

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"server/internal/domain/dto"
	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
	"server/internal/helper"
	"server/internal/storage"
)

type GameService struct {
	gameRepo       domainrepo.GameRepository
	userRepo       domainrepo.UserRepository
	categoryRepo   domainrepo.CategoryRepository
	difficultyRepo domainrepo.DifficultyRepository
	minioStorage   *storage.MinIOStorage
	cfg            helper.Config
}

func NewGameService(
	gameRepo domainrepo.GameRepository,
	userRepo domainrepo.UserRepository,
	categoryRepo domainrepo.CategoryRepository,
	difficultyRepo domainrepo.DifficultyRepository,
	minioStorage *storage.MinIOStorage,
	cfg helper.Config,
) *GameService {
	return &GameService{
		gameRepo:       gameRepo,
		userRepo:       userRepo,
		categoryRepo:   categoryRepo,
		difficultyRepo: difficultyRepo,
		minioStorage:   minioStorage,
		cfg:            cfg,
	}
}

func (s *GameService) UploadGame(developerID uint, title string, description string, categoryIDs []uint, difficultyID uint, fileHeader *multipart.FileHeader, thumbnailHeader *multipart.FileHeader) (*models.Game, error) {
	developer, err := s.userRepo.FindByID(developerID)
	if err != nil {
		return nil, err
	}
	if developer == nil {
		return nil, errors.New("developer not found")
	}
	// if developer.Role != "DEVELOPER" && developer.Role != "ADMIN" {
	// 	return nil, errors.New("forbidden")
	// }
	if strings.ToLower(filepath.Ext(fileHeader.Filename)) != ".zip" {
		return nil, errors.New("only zip file allowed")
	}

	categories, err := s.resolveCategories(categoryIDs)
	if err != nil {
		return nil, err
	}
	difficulty, err := s.resolveDifficulty(difficultyID)
	if err != nil {
		return nil, err
	}

	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	tempRoot, err := os.MkdirTemp("", "game-upload-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempRoot)

	zipPath := filepath.Join(tempRoot, fileHeader.Filename)
	if err := helper.CopyMultipartFile(src, zipPath); err != nil {
		return nil, err
	}

	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, errors.New("invalid zip file")
	}
	defer reader.Close()
	if len(reader.File) == 0 {
		return nil, errors.New("zip file is empty")
	}

	extractDir := filepath.Join(tempRoot, "extracted")
	if err := helper.ExtractZip(zipPath, extractDir); err != nil {
		return nil, err
	}
	if !hasRootIndexHTML(extractDir) {
		return nil, errors.New("zip must contain index.html in root")
	}
	if err := normalizeExtractedAssetPaths(extractDir); err != nil {
		return nil, err
	}

	objectPrefix := fmt.Sprintf("%d/%s", developerID, helper.Slugify(title))
	if err := s.minioStorage.UploadDirectory(objectPrefix, extractDir); err != nil {
		return nil, err
	}

	thumbnailPath, err := s.uploadThumbnail(title, thumbnailHeader)
	if err != nil {
		return nil, err
	}

	game := &models.Game{
		Title:         title,
		Description:   description,
		FileURL:       objectPrefix,
		ThumbnailPath: thumbnailPath,
		DeveloperID:   developerID,
		DifficultyID:  difficulty.ID,
		Difficulty:    *difficulty,
		Status:        "pending",
		Categories:    categories,
	}
	if err := s.gameRepo.Create(game); err != nil {
		return nil, err
	}

	return s.gameRepo.FindByID(game.ID)
}

func (s *GameService) UploadGameResponse(developerID uint, title string, description string, categoryIDs []uint, difficultyID uint, fileHeader *multipart.FileHeader, thumbnailHeader *multipart.FileHeader) (*dto.GameResponse, error) {
	game, err := s.UploadGame(developerID, title, description, categoryIDs, difficultyID, fileHeader, thumbnailHeader)
	if err != nil {
		return nil, err
	}

	response := dto.BuildGameResponse(game, s.BuildThumbnailURL)
	return &response, nil
}

func (s *GameService) UpdateGame(gameID uint, actor *models.User, title string, description string, categoryIDs []uint, difficultyID *uint, fileHeader *multipart.FileHeader, thumbnailHeader *multipart.FileHeader) (*models.Game, error) {
	game, err := s.gameRepo.FindByID(gameID)
	if err != nil {
		return nil, err
	}
	if game == nil {
		return nil, errors.New("game not found")
	}
	if actor.Role != "ADMIN" && game.DeveloperID != actor.ID {
		return nil, errors.New("forbidden")
	}

	if strings.TrimSpace(title) != "" {
		game.Title = title
	}
	if description != "" {
		game.Description = description
	}
	if categoryIDs != nil {
		categories, err := s.resolveCategories(categoryIDs)
		if err != nil {
			return nil, err
		}
		game.Categories = categories
	}
	if difficultyID != nil {
		difficulty, err := s.resolveDifficulty(*difficultyID)
		if err != nil {
			return nil, err
		}
		game.DifficultyID = difficulty.ID
		game.Difficulty = *difficulty
	}
	if fileHeader != nil {
		if err := s.replaceGameArchive(game, fileHeader); err != nil {
			return nil, err
		}
	}
	if thumbnailHeader != nil {
		thumbnailPath, err := s.uploadThumbnail(game.Title, thumbnailHeader)
		if err != nil {
			return nil, err
		}
		game.ThumbnailPath = thumbnailPath
	}

	if err := s.gameRepo.Update(game); err != nil {
		return nil, err
	}

	return s.gameRepo.FindByID(game.ID)
}

func (s *GameService) UpdateGameResponse(gameID uint, actor *models.User, title string, description string, categoryIDs []uint, difficultyID *uint, fileHeader *multipart.FileHeader, thumbnailHeader *multipart.FileHeader) (*dto.GameResponse, error) {
	game, err := s.UpdateGame(gameID, actor, title, description, categoryIDs, difficultyID, fileHeader, thumbnailHeader)
	if err != nil {
		return nil, err
	}

	response := dto.BuildGameResponse(game, s.BuildThumbnailURL)
	return &response, nil
}

func (s *GameService) ListApprovedGames() ([]models.Game, error) {
	return s.gameRepo.ListApproved()
}

func (s *GameService) ListApprovedGameResponses() ([]dto.GameResponse, error) {
	games, err := s.ListApprovedGames()
	if err != nil {
		return nil, err
	}
	return dto.BuildGameResponses(games, s.BuildThumbnailURL), nil
}

func (s *GameService) GetApprovedGame(id uint) (*models.Game, error) {
	game, err := s.gameRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if game == nil || game.Status != "approved" {
		return nil, errors.New("game not found")
	}
	return game, nil
}

func (s *GameService) GetApprovedGameResponse(id uint) (*dto.GameResponse, error) {
	game, err := s.GetApprovedGame(id)
	if err != nil {
		return nil, err
	}

	response := dto.BuildGameResponse(game, s.BuildThumbnailURL)
	return &response, nil
}

func (s *GameService) ListGamesByDeveloperID(developerID uint) ([]models.Game, error) {
	return s.gameRepo.ListByDeveloperID(developerID)
}

func (s *GameService) ListGameResponsesByDeveloperID(developerID uint) ([]dto.GameResponse, error) {
	games, err := s.ListGamesByDeveloperID(developerID)
	if err != nil {
		return nil, err
	}
	return dto.BuildGameResponses(games, s.BuildThumbnailURL), nil
}

func (s *GameService) PlayGame(id uint) (string, error) {
	game, err := s.GetApprovedGame(id)
	if err != nil {
		return "", err
	}

	objectPrefix := s.minioStorage.ExtractObjectPrefix(game.FileURL)
	if objectPrefix == "" {
		return "", errors.New("game file is missing")
	}

	return fmt.Sprintf("/play/%d/", game.ID), nil
}

func (s *GameService) PlayGameResponse(id uint) (*dto.FileURLResponse, error) {
	url, err := s.PlayGame(id)
	if err != nil {
		return nil, err
	}
	return &dto.FileURLResponse{FileURL: url}, nil
}

func (s *GameService) ListAllGames() ([]models.Game, error) {
	return s.gameRepo.ListAll()
}

func (s *GameService) ApproveGame(id uint) error {
	game, err := s.gameRepo.FindByID(id)
	if err != nil {
		return err
	}
	if game == nil {
		return errors.New("game not found")
	}
	game.Status = "approved"
	return s.gameRepo.Update(game)
}

func (s *GameService) RejectGame(id uint) error {
	game, err := s.gameRepo.FindByID(id)
	if err != nil {
		return err
	}
	if game == nil {
		return errors.New("game not found")
	}
	game.Status = "rejected"
	return s.gameRepo.Update(game)
}

func (s *GameService) ListCategories() ([]models.Category, error) {
	return s.categoryRepo.ListAll()
}

func (s *GameService) ListCategoryResponses() ([]dto.CategoryResponse, error) {
	categories, err := s.ListCategories()
	if err != nil {
		return nil, err
	}
	return dto.BuildCategoryResponses(categories), nil
}

func (s *GameService) ListDifficulties() ([]models.Difficulty, error) {
	return s.difficultyRepo.ListAll()
}

func (s *GameService) ListDifficultyResponses() ([]dto.DifficultyResponse, error) {
	difficulties, err := s.ListDifficulties()
	if err != nil {
		return nil, err
	}
	return dto.BuildDifficultyResponses(difficulties), nil
}

func (s *GameService) CreateCategory(name string) (*models.Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	category := &models.Category{Name: name}
	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *GameService) CreateCategoryResponse(name string) (*dto.CategoryResponse, error) {
	category, err := s.CreateCategory(name)
	if err != nil {
		return nil, err
	}

	response := dto.BuildCategoryResponse(*category)
	return &response, nil
}

func (s *GameService) UpdateCategory(id uint, name string) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}

	category.Name = name
	if err := s.categoryRepo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *GameService) UpdateCategoryResponse(id uint, name string) (*dto.CategoryResponse, error) {
	category, err := s.UpdateCategory(id, name)
	if err != nil {
		return nil, err
	}

	response := dto.BuildCategoryResponse(*category)
	return &response, nil
}

func (s *GameService) DeleteCategory(id uint) error {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}
	return s.categoryRepo.Delete(id)
}

func (s *GameService) BuildThumbnailURL(objectName string) string {
	if strings.TrimSpace(objectName) == "" {
		return ""
	}
	return "/games/thumbnail/" + strings.TrimLeft(objectName, "/")
}

func (s *GameService) OpenThumbnailAsset(objectName string) (io.ReadCloser, string, error) {
	cleanObjectName := strings.TrimPrefix(filepath.ToSlash(filepath.Clean("/"+objectName)), "/")
	if cleanObjectName == "" || cleanObjectName == "." {
		return nil, "", errors.New("thumbnail not found")
	}
	for _, segment := range strings.Split(cleanObjectName, "/") {
		if segment == ".." {
			return nil, "", errors.New("invalid thumbnail path")
		}
	}

	object, info, err := s.minioStorage.GetThumbnailObject(cleanObjectName)
	if err != nil {
		return nil, "", err
	}
	return object, info.ContentType, nil
}

func (s *GameService) OpenGameAsset(id uint, assetPath string) (io.ReadCloser, string, error) {
	game, err := s.GetApprovedGame(id)
	if err != nil {
		return nil, "", err
	}

	objectPrefix := s.minioStorage.ExtractObjectPrefix(game.FileURL)
	if objectPrefix == "" {
		return nil, "", errors.New("game file is missing")
	}

	for _, segment := range strings.Split(filepath.ToSlash(assetPath), "/") {
		if segment == ".." {
			return nil, "", errors.New("invalid asset path")
		}
	}

	cleanAssetPath := strings.TrimPrefix(filepath.ToSlash(filepath.Clean("/"+assetPath)), "/")
	if cleanAssetPath == "" {
		cleanAssetPath = "index.html"
	}

	objectName := fmt.Sprintf("%s/%s", strings.Trim(objectPrefix, "/"), cleanAssetPath)
	object, info, err := s.minioStorage.GetObject(objectName)
	if err != nil {
		return nil, "", err
	}

	return object, info.ContentType, nil
}

func hasRootIndexHTML(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "index.html"))
	return err == nil && !info.IsDir()
}

func (s *GameService) resolveCategories(categoryIDs []uint) ([]models.Category, error) {
	if len(categoryIDs) == 0 {
		return []models.Category{}, nil
	}

	deduped := make([]uint, 0, len(categoryIDs))
	seen := make(map[uint]struct{}, len(categoryIDs))
	for _, id := range categoryIDs {
		if id == 0 {
			return nil, errors.New("invalid category id")
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		deduped = append(deduped, id)
	}

	categories, err := s.categoryRepo.FindByIDs(deduped)
	if err != nil {
		return nil, err
	}
	if len(categories) != len(deduped) {
		foundIDs := make([]uint, 0, len(categories))
		for _, category := range categories {
			foundIDs = append(foundIDs, category.ID)
		}
		for _, id := range deduped {
			if !slices.Contains(foundIDs, id) {
				return nil, fmt.Errorf("category %d not found", id)
			}
		}
	}

	return categories, nil
}

func (s *GameService) resolveDifficulty(difficultyID uint) (*models.Difficulty, error) {
	if difficultyID == 0 {
		return nil, errors.New("invalid difficulty id")
	}

	difficulty, err := s.difficultyRepo.FindByID(difficultyID)
	if err != nil {
		return nil, err
	}
	if difficulty == nil {
		return nil, fmt.Errorf("difficulty %d not found", difficultyID)
	}

	return difficulty, nil
}

func (s *GameService) replaceGameArchive(game *models.Game, fileHeader *multipart.FileHeader) error {
	if strings.ToLower(filepath.Ext(fileHeader.Filename)) != ".zip" {
		return errors.New("only zip file allowed")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	tempRoot, err := os.MkdirTemp("", "game-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempRoot)

	zipPath := filepath.Join(tempRoot, fileHeader.Filename)
	if err := helper.CopyMultipartFile(src, zipPath); err != nil {
		return err
	}
	if err := validateAndExtractGameArchive(zipPath, filepath.Join(tempRoot, "extracted")); err != nil {
		return err
	}

	objectPrefix := fmt.Sprintf("%d/%s", game.DeveloperID, helper.Slugify(game.Title))
	if err := s.minioStorage.UploadDirectory(objectPrefix, filepath.Join(tempRoot, "extracted")); err != nil {
		return err
	}
	game.FileURL = objectPrefix
	return nil
}

func (s *GameService) uploadThumbnail(title string, thumbnailHeader *multipart.FileHeader) (string, error) {
	if thumbnailHeader == nil {
		return "", nil
	}

	src, err := thumbnailHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	contentType := thumbnailHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(thumbnailHeader.Filename)))
	}
	if !strings.HasPrefix(contentType, "image/") {
		return "", errors.New("thumbnail must be an image")
	}

	objectName := fmt.Sprintf("%s%s", helper.Slugify(title), strings.ToLower(filepath.Ext(thumbnailHeader.Filename)))
	if err := s.minioStorage.UploadFile(s.minioStorage.ThumbnailBucket(), objectName, src, thumbnailHeader.Size, contentType); err != nil {
		return "", err
	}
	return objectName, nil
}

func validateAndExtractGameArchive(zipPath string, extractDir string) error {
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
	if !hasRootIndexHTML(extractDir) {
		return errors.New("zip must contain index.html in root")
	}
	return normalizeExtractedAssetPaths(extractDir)
}

var cssURLPattern = regexp.MustCompile(`url\(\s*[^)\s]+\s*\)`)

func normalizeExtractedAssetPaths(root string) error {
	return filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !isNormalizableTextAsset(filePath) {
			return nil
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		if strings.IndexByte(string(content), 0) >= 0 {
			return nil
		}

		updated := normalizeQuotedAssetPaths(string(content), '\'')
		updated = normalizeQuotedAssetPaths(updated, '"')
		updated = normalizeQuotedAssetPaths(updated, '`')
		updated = normalizeCSSURLAssetPaths(updated)

		if updated == string(content) {
			return nil
		}

		return os.WriteFile(filePath, []byte(updated), info.Mode())
	})
}

func isNormalizableTextAsset(filePath string) bool {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".html", ".htm", ".css", ".js", ".mjs":
		return true
	default:
		return false
	}
}

func normalizeQuotedAssetPaths(content string, quote rune) string {
	pattern := regexp.MustCompile(regexp.QuoteMeta(string(quote)) + `((?:\./|\.\./)+[^` + regexp.QuoteMeta(string(quote)) + "\r\n]+)" + regexp.QuoteMeta(string(quote)))

	return pattern.ReplaceAllStringFunc(content, func(match string) string {
		rawPath := match[1 : len(match)-1]
		normalized, ok := normalizeAssetReference(rawPath)
		if !ok {
			return match
		}
		return string(quote) + normalized + string(quote)
	})
}

func normalizeCSSURLAssetPaths(content string) string {
	return cssURLPattern.ReplaceAllStringFunc(content, func(match string) string {
		inner := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(match, "url("), ")"))
		if inner == "" {
			return match
		}

		quote := ""
		rawPath := inner
		if strings.HasPrefix(inner, "\"") && strings.HasSuffix(inner, "\"") && len(inner) >= 2 {
			quote = `"`
			rawPath = inner[1 : len(inner)-1]
		} else if strings.HasPrefix(inner, "'") && strings.HasSuffix(inner, "'") && len(inner) >= 2 {
			quote = "'"
			rawPath = inner[1 : len(inner)-1]
		}

		normalized, ok := normalizeAssetReference(rawPath)
		if !ok {
			return match
		}

		return "url(" + quote + normalized + quote + ")"
	})
}

func normalizeAssetReference(value string) (string, bool) {
	raw := strings.TrimSpace(value)
	if !strings.HasPrefix(raw, "./") && !strings.HasPrefix(raw, "../") {
		return "", false
	}
	if !hasNormalizableAssetExtension(raw) {
		return "", false
	}

	cleaned := strings.TrimPrefix(path.Clean("/"+raw), "/")
	if cleaned == "" || cleaned == "." {
		return "", false
	}

	return cleaned, true
}

func hasNormalizableAssetExtension(value string) bool {
	trimmed := strings.ToLower(strings.TrimSpace(value))
	trimmed = strings.Split(trimmed, "?")[0]
	trimmed = strings.Split(trimmed, "#")[0]

	switch path.Ext(trimmed) {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico", ".bmp",
		".mp3", ".wav", ".ogg", ".m4a",
		".mp4", ".webm",
		".woff", ".woff2", ".ttf", ".otf", ".eot",
		".json", ".txt":
		return true
	default:
		return false
	}
}
