package service

import (
	"context"
	"voca-plane/internal/repository"
	"voca-plane/internal/seeders"

	"gorm.io/gorm"
)

type SystemService interface {
	Reset (ctx context.Context) error
	ResetAndSeed(ctx context.Context) error
}

type systemService struct {
	db         *gorm.DB
	systemRepo repository.SystemRepository
}

func NewSystemService(db *gorm.DB, systemRepo repository.SystemRepository) SystemService {
	return &systemService{
		db:         db,
		systemRepo: systemRepo,
	}
}

func (s *systemService) ResetAndSeed(ctx context.Context) error {
	if err := s.systemRepo.ResetDatabase(ctx); err != nil {
		return err
	}

	seeders.SeedAll(s.db)
	return nil
}

func (s *systemService) Reset(ctx context.Context) error {
	if err := s.systemRepo.ResetDatabase(ctx); err != nil {
		return err
	}

	seeders.SeedMain(s.db)
	return nil
}