package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"server/internal/database"
	"server/internal/handler"
	"server/internal/helper"
	"server/internal/middleware"
	"server/internal/repository"
	"server/internal/service"
	"server/internal/storage"
	"server/internal/websocket"
)

type App struct {
	router *gin.Engine
	server *http.Server
	hub    *websocket.Hub
}

func New() (*App, error) {
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

	clerkVerifier, err := helper.NewClerkVerifier(cfg)
	if err != nil {
		return nil, err
	}
	clerkClient := helper.NewClerkClient(cfg)

	userRepo := repository.NewUserRepository(db)
	friendRepo := repository.NewFriendRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	gameRepo := repository.NewGameRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	difficultyRepo := repository.NewDifficultyRepository(db)
	postRepo := repository.NewPostRepository(db)

	userService := service.NewUserService(userRepo, clerkClient)
	friendService := service.NewFriendService(friendRepo, userRepo)
	chatService := service.NewChatService(messageRepo, userRepo)
	gameService := service.NewGameService(gameRepo, userRepo, categoryRepo, difficultyRepo, minioClient, cfg)
	postService := service.NewPostService(postRepo)
	adminService := service.NewAdminService(userRepo, gameRepo, redisClient)

	friendHandler := handler.NewFriendHandler(friendService)
	chatHandler := handler.NewChatHandler(chatService)
	gameHandler := handler.NewGameHandler(gameService)
	postHandler := handler.NewPostHandler(postService)
	adminHandler := handler.NewAdminHandler(adminService, gameService)
	userHandler := handler.NewUserHandler()

	hub := websocket.NewHub(chatService)
	websocket.SetRedisClient(redisClient)
	go hub.Run()

	router := gin.Default()
	router.MaxMultipartMemory = 128 << 20

	corsMiddleware := middleware.NewCORSMiddleware(cfg.CORSOrigins)
	authMiddleware := middleware.NewClerkMiddleware(clerkVerifier, userService)
	roleMiddleware := middleware.NewRoleMiddleware()

	router.Use(corsMiddleware.Handle())
	router.GET("", func(c *gin.Context) {
		helper.Success(c, http.StatusOK, "welcome to Voca Hub API", nil)
	})

	router.GET("/health", func(c *gin.Context) {
		helper.Success(c, http.StatusOK, "ok", gin.H{
			"service": "server",
		})
	})

	router.GET("/play/:id", gameHandler.ServeGameFile)
	router.GET("/play/:id/*filepath", gameHandler.ServeGameFile)
	router.GET("/games/thumbnail/*filepath", gameHandler.ServeThumbnail)

	router.GET("/categories", gameHandler.ListCategories)

	difficulties := router.Group("/difficulties")
	{
		difficulties.GET("", gameHandler.ListDifficulties)
	}

	router.GET("/games", gameHandler.ListApprovedGames)

	router.GET("/posts", postHandler.List)

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			helper.Error(c, http.StatusNotFound, "endpoint tidak ditemukan")
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/play") {
			gameHandler.ServeRootAssetFallback(c)
			return
		}

		// Fallback default jika tidak cocok keduanya
		helper.Error(c, http.StatusNotFound, "endpoint tidak ditemukan")
	})

	api := router.Group("/api")
	api.Use(authMiddleware.Handle())
	{
		friend := api.Group("/friends")
		{
			friend.POST("/request", friendHandler.AddFriend)
			friend.POST("/:id/accept", friendHandler.AcceptFriend)
			friend.POST("/:id/reject", friendHandler.RejectFriend)
			friend.GET("", friendHandler.ListFriends)
			friend.GET("/pending", friendHandler.ListPendingRequests)
		}

		chat := api.Group("/chat")
		{
			chat.GET("/history/:user_id", chatHandler.GetHistory)
			chat.GET("/ws", func(c *gin.Context) {
				websocket.ServeWS(hub, c)
			})
		}

		users := api.Group("/users")
		{
			users.GET("/me", userHandler.Me)
		}

		games := api.Group("/games")
		{
			games.GET("/mine", gameHandler.ListMyGames)
			games.GET("/:id", gameHandler.GetApprovedGame)
			games.GET("/:id/play", gameHandler.PlayGame)
			games.Use(roleMiddleware.Require("USER", "DEVELOPER", "ADMIN"))
			games.POST("/upload", gameHandler.UploadGame)
			games.PUT("/:id", gameHandler.UpdateGame)
		}

		categories := api.Group("/categories")
		{
			categories.Use(roleMiddleware.Require("ADMIN"))
			categories.POST("", gameHandler.CreateCategory)
			categories.PUT("/:id", gameHandler.UpdateCategory)
			categories.DELETE("/:id", gameHandler.DeleteCategory)
		}

		posts := api.Group("/posts")
		{
			posts.GET("/mine", postHandler.ListMine)
			posts.GET("/:id", postHandler.GetByID)
			posts.POST("", postHandler.Create)
			posts.PUT("/:id", postHandler.Update)
			posts.DELETE("/:id", postHandler.Delete)
		}

		admin := api.Group("/admin")
		admin.Use(roleMiddleware.Require("ADMIN"))
		{
			admin.GET("/dashboard", adminHandler.Dashboard)
			admin.GET("/users", adminHandler.ListUsers)
			admin.GET("/games", adminHandler.ListGames)
			admin.POST("/games/:id/approve", adminHandler.ApproveGame)
			admin.POST("/games/:id/reject", adminHandler.RejectGame)
		}
	}

	addr := fmt.Sprintf(":%s", strings.TrimSpace(cfg.AppPort))
	server := &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return &App{
		router: router,
		server: server,
		hub:    hub,
	}, nil
}

func (a *App) Run() error {
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}

func init() {
	gin.SetMode(getGinMode())
}

func getGinMode() string {
	mode := strings.TrimSpace(os.Getenv("GIN_MODE"))
	if mode == "" {
		return gin.DebugMode
	}
	return mode
}
