package service

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
	"server/internal/helper"
	"server/internal/storage"
)

type GameService struct {
	gameRepo     domainrepo.GameRepository
	userRepo     domainrepo.UserRepository
	minioStorage *storage.MinIOStorage
	cfg          helper.Config
}

func NewGameService(
	gameRepo domainrepo.GameRepository,
	userRepo domainrepo.UserRepository,
	minioStorage *storage.MinIOStorage,
	cfg helper.Config,
) *GameService {
	return &GameService{
		gameRepo:     gameRepo,
		userRepo:     userRepo,
		minioStorage: minioStorage,
		cfg:          cfg,
	}
}

func (s *GameService) UploadGame(developerID uint, title string, description string, fileHeader *multipart.FileHeader) (*models.Game, error) {
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

	objectPrefix := fmt.Sprintf("%d/%s", developerID, helper.Slugify(title))
	if err := s.minioStorage.UploadDirectory(objectPrefix, extractDir); err != nil {
		return nil, err
	}

	game := &models.Game{
		Title:       title,
		Description: description,
		FileURL:     objectPrefix,
		DeveloperID: developerID,
		Status:      "pending",
	}
	if err := s.gameRepo.Create(game); err != nil {
		return nil, err
	}

	return s.gameRepo.FindByID(game.ID)
}

func (s *GameService) ListApprovedGames() ([]models.Game, error) {
	return s.gameRepo.ListApproved()
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

func (s *GameService) ListGamesByDeveloperID(developerID uint) ([]models.Game, error) {
	return s.gameRepo.ListByDeveloperID(developerID)
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
