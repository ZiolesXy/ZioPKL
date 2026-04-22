package service

import (
	"context"

	"github.com/redis/go-redis/v9"

	"server/internal/domain/dto"
	domainrepo "server/internal/domain/repository"
)

type AdminService struct {
	userRepo domainrepo.UserRepository
	gameRepo domainrepo.GameRepository
	redis    *redis.Client
}

func NewAdminService(userRepo domainrepo.UserRepository, gameRepo domainrepo.GameRepository, redisClient *redis.Client) *AdminService {
	return &AdminService{
		userRepo: userRepo,
		gameRepo: gameRepo,
		redis:    redisClient,
	}
}

func (s *AdminService) Dashboard() (*dto.DashboardStatsResponse, error) {
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

	return &dto.DashboardStatsResponse{
		TotalUsers:       totalUsers,
		TotalDevelopers:  totalDevelopers,
		TotalGames:       totalGames,
		TotalActiveChats: totalActiveChats,
	}, nil
}

func (s *AdminService) ListUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.List()
	if err != nil {
		return nil, err
	}
	return dto.BuildUserResponses(users), nil
}

func (s *AdminService) ListGames(buildThumbnailURL func(string) string) ([]dto.GameResponse, error) {
	games, err := s.gameRepo.ListAll()
	if err != nil {
		return nil, err
	}
	return dto.BuildGameResponses(games, buildThumbnailURL), nil
}
