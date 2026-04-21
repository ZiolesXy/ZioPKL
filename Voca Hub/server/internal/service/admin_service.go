package service

import (
	"context"

	"github.com/redis/go-redis/v9"

	"server/internal/domain/models"
	domainrepo "server/internal/domain/repository"
)

type AdminService struct {
	userRepo domainrepo.UserRepository
	gameRepo domainrepo.GameRepository
	redis    *redis.Client
}

type DashboardStats struct {
	TotalUsers       int64 `json:"total_users"`
	TotalDevelopers  int64 `json:"total_developers"`
	TotalGames       int64 `json:"total_games"`
	TotalActiveChats int64 `json:"total_active_chats"`
}

func NewAdminService(userRepo domainrepo.UserRepository, gameRepo domainrepo.GameRepository, redisClient *redis.Client) *AdminService {
	return &AdminService{
		userRepo: userRepo,
		gameRepo: gameRepo,
		redis:    redisClient,
	}
}

func (s *AdminService) Dashboard() (*DashboardStats, error) {
	totalUsers, err := s.userRepo.CountAll()
	if err != nil {
		return nil, err
	}
	totalDevelopers, err := s.userRepo.CountByRole("DEVELOPER")
	if err != nil {
		return nil, err
	}
	totalGames, err := s.gameRepo.CountAll()
	if err != nil {
		return nil, err
	}
	totalActiveChats, err := s.redis.Get(context.Background(), "active_chats").Int64()
	if err == redis.Nil {
		totalActiveChats = 0
	} else if err != nil {
		return nil, err
	}

	return &DashboardStats{
		TotalUsers:       totalUsers,
		TotalDevelopers:  totalDevelopers,
		TotalGames:       totalGames,
		TotalActiveChats: totalActiveChats,
	}, nil
}

func (s *AdminService) ListUsers() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *AdminService) ListGames() ([]models.Game, error) {
	return s.gameRepo.ListAll()
}
