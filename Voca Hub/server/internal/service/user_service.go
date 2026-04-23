package service

import (
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"

	"server/internal/domain/dto"
	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
	"server/internal/storage"
)

type UserService struct {
	userRepo     domainrepo.UserRepository
	minioStorage *storage.MinIOStorage
}

func NewUserService(userRepo domainrepo.UserRepository, minioStorage *storage.MinIOStorage) *UserService {
	return &UserService{
		userRepo:     userRepo,
		minioStorage: minioStorage,
	}
}

func (s *UserService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) UpdateProfile(userID uint, username *string, profileHeader *multipart.FileHeader) (*models.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if username != nil {
		user.Username = normalizeOptionalString(username)
	}

	if profileHeader != nil {
		profilePath, err := s.replaceProfileImage(user, profileHeader)
		if err != nil {
			return nil, err
		}
		user.ProfileURL = &profilePath
	}

	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(userID)
}

func (s *UserService) UpdateProfileResponse(userID uint, username *string, profileHeader *multipart.FileHeader) (*dto.UserResponse, error) {
	user, err := s.UpdateProfile(userID, username, profileHeader)
	if err != nil {
		return nil, err
	}

	response := dto.BuildUserResponse(*user)
	return &response, nil
}

func (s *UserService) replaceProfileImage(user *models.User, profileHeader *multipart.FileHeader) (string, error) {
	src, err := profileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	contentType := profileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(profileHeader.Filename)))
	}
	if !strings.HasPrefix(strings.ToLower(contentType), "image/") {
		return "", errors.New("profile image must be an image")
	}

	objectName := filepath.ToSlash(filepath.Join("users", normalizeProfileSegment(user.Email), "profile"+strings.ToLower(filepath.Ext(profileHeader.Filename))))
	if err := s.minioStorage.UploadFile(s.minioStorage.ProfileBucket(), objectName, src, profileHeader.Size, contentType); err != nil {
		return "", err
	}

	if user.ProfileURL != nil {
		oldObject := s.minioStorage.ExtractProfileObjectName(*user.ProfileURL)
		if oldObject != "" && oldObject != objectName {
			if err := s.minioStorage.RemoveProfileObject(oldObject); err != nil {
				return "", err
			}
		}
	}

	return objectName, nil
}

func (s *UserService) OpenProfileAsset(objectName string) (io.ReadCloser, string, error) {
	cleanObjectName := strings.TrimPrefix(filepath.ToSlash(filepath.Clean("/"+objectName)), "/")
	if cleanObjectName == "" || cleanObjectName == "." {
		return nil, "", errors.New("profile not found")
	}
	for _, segment := range strings.Split(cleanObjectName, "/") {
		if segment == ".." {
			return nil, "", errors.New("invalid profile path")
		}
	}

	object, info, err := s.minioStorage.GetProfileObject(cleanObjectName)
	if err != nil {
		return nil, "", err
	}
	return object, info.ContentType, nil
}

func normalizeProfileSegment(value string) string {
	replacer := strings.NewReplacer("@", "-at-", ".", "-", "_", "-", "+", "-plus-", " ", "-")
	return replacer.Replace(strings.ToLower(strings.TrimSpace(value)))
}
