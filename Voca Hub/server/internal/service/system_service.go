package service

import (
	"server/internal/helper"
)

type SystemService interface {
	GenerateSecret() (string, error)
}

type systemService struct {
}

func NewSystemService() SystemService {
	return &systemService{}
}

func (s *systemService) GenerateSecret() (string, error) {
	return helper.GenerateRandomSecret(32)
}