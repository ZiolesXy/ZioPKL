package app

import (
	"github.com/redis/go-redis/v9"

	"server/internal/database"
	"server/internal/handler"
	"server/internal/helper"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"
	"server/internal/storage"
)

type Container struct {
	Config helper.Config

	RedisClient *redis.Client

	TokenManager *helper.TokenManager

	UserService   *service.UserService
	TokenStore    *service.TokenStoreService
	AuthService   *service.AuthService
	FriendService *service.FriendService
	ChatService   *service.ChatService
	GameService   *service.GameService
	PostService   *service.PostService
	AdminService  *service.AdminService
	SystemService service.SystemService

	AuthHandler   *handler.AuthHandler
	UserHandler   *handler.UserHandler
	FriendHandler *handler.FriendHandler
	ChatHandler   *handler.ChatHandler
	GameHandler   *handler.GameHandler
	PostHandler   *handler.PostHandler
	AdminHandler  *handler.AdminHandler
	SystemHandler *handler.SystemHandler

	CORSMiddleware *middleware.CORSMiddleware
	AuthMiddleware *middleware.AuthMiddleware
	RoleMiddleware *middleware.RoleMiddleware
}

func NewContainer() (*Container, error) {
	cfg := helper.LoadConfig()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	redisClient, err := database.NewRedis(cfg)
	if err != nil {
		return nil, err
	}

	minioClient, err := storage.NewMinIO(cfg)
	if err != nil {
		return nil, err
	}

	tokenManager, err := helper.NewTokenManager(cfg)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	friendRepo := repository.NewFriendRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	gameRepo := repository.NewGameRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	difficultyRepo := repository.NewDifficultyRepository(db)
	postRepo := repository.NewPostRepository(db)

	userService := service.NewUserService(userRepo)
	tokenStore := service.NewTokenStoreService(redisClient)
	authService := service.NewAuthService(userRepo, tokenManager, tokenStore)
	friendService := service.NewFriendService(friendRepo, userRepo)
	chatService := service.NewChatService(messageRepo, userRepo)
	gameService := service.NewGameService(gameRepo, userRepo, categoryRepo, difficultyRepo, minioClient, cfg)
	postService := service.NewPostService(postRepo)
	adminService := service.NewAdminService(userRepo, gameRepo, redisClient)
	systemService := service.NewSystemService()

	return &Container{
		Config: cfg,

		RedisClient: redisClient,

		TokenManager: tokenManager,

		UserService:   userService,
		TokenStore:    tokenStore,
		AuthService:   authService,
		FriendService: friendService,
		ChatService:   chatService,
		GameService:   gameService,
		PostService:   postService,
		AdminService:  adminService,
		SystemService: systemService,

		AuthHandler:   handler.NewAuthHandler(authService),
		UserHandler:   handler.NewUserHandler(),
		FriendHandler: handler.NewFriendHandler(friendService),
		ChatHandler:   handler.NewChatHandler(chatService),
		GameHandler:   handler.NewGameHandler(gameService),
		PostHandler:   handler.NewPostHandler(postService),
		AdminHandler:  handler.NewAdminHandler(adminService, gameService),
		SystemHandler: handler.NewSystemHandler(systemService),

		CORSMiddleware: middleware.NewCORSMiddleware(cfg.CORSOrigins),
		AuthMiddleware: middleware.NewAuthMiddleware(tokenManager, tokenStore, userService),
		RoleMiddleware: middleware.NewRoleMiddleware(),
	}, nil
}
