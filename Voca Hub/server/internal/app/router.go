package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"server/internal/domain/dto"
	"server/internal/helper"
	"server/internal/websocket"
)

func SetupRouter(container *Container, hub *websocket.Hub) *gin.Engine {
	router := gin.Default()
	router.MaxMultipartMemory = 128 << 20

	router.Use(container.CORSMiddleware.Handle())

	router.GET("", func(c *gin.Context) {
		helper.Success(c, http.StatusOK, "welcome to Voca Hub API", nil)
	})

	router.GET("/health", func(c *gin.Context) {
		helper.Success(c, http.StatusOK, "ok", dto.HealthResponse{Service: "server"})
	})

	router.GET("/password", container.SystemHandler.GetNewSecret)

	router.GET("/play/:id", container.GameHandler.ServeGameFile)
	router.GET("/play/:id/*filepath", container.GameHandler.ServeGameFile)
	router.GET("/games/thumbnail/*filepath", container.GameHandler.ServeThumbnail)
	router.GET("/users/profile/*filepath", container.UserHandler.ServeProfile)

	router.GET("/categories", container.GameHandler.ListCategories)

	difficulties := router.Group("/difficulties")
	{
		difficulties.GET("", container.GameHandler.ListDifficulties)
	}

	router.GET("/games", container.GameHandler.ListApprovedGames)
	router.GET("/posts", container.PostHandler.List)

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			helper.Error(c, http.StatusNotFound, "endpoint tidak ditemukan")
			return
		}

		if strings.HasPrefix(c.Request.URL.Path, "/play") {
			container.GameHandler.ServeRootAssetFallback(c)
			return
		}

		helper.Error(c, http.StatusNotFound, "endpoint tidak ditemukan")
	})

	api := router.Group("/api")
	publicAuth := api.Group("/auth")
	{
		publicAuth.POST("/register", container.AuthHandler.Register)
		publicAuth.POST("/login", container.AuthHandler.Login)
		publicAuth.POST("/refresh", container.AuthHandler.Refresh)
	}

	protected := router.Group("/api")
	protected.Use(container.AuthMiddleware.Handle())
	{
		protectedAuth := protected.Group("/auth")
		{
			protectedAuth.POST("/logout", container.AuthHandler.Logout)
		}

		friend := protected.Group("/friends")
		{
			friend.POST("/request", container.FriendHandler.AddFriend)
			friend.POST("/:id/accept", container.FriendHandler.AcceptFriend)
			friend.POST("/:id/reject", container.FriendHandler.RejectFriend)
			friend.GET("", container.FriendHandler.ListFriends)
			friend.GET("/pending", container.FriendHandler.ListPendingRequests)
		}

		chat := protected.Group("/chat")
		{
			chat.GET("/history/:user_id", container.ChatHandler.GetHistory)
			chat.GET("/ws", func(c *gin.Context) {
				websocket.ServeWS(hub, c)
			})
		}

		users := protected.Group("/users")
		{
			users.GET("/me", container.UserHandler.Me)
			users.PUT("/me", container.UserHandler.UpdateMe)
		}

		games := protected.Group("/games")
		{
			games.GET("/mine", container.GameHandler.ListMyGames)
			games.GET("/:id", container.GameHandler.GetApprovedGame)
			games.GET("/:id/play", container.GameHandler.PlayGame)
			games.Use(container.RoleMiddleware.Require("USER", "DEVELOPER", "ADMIN"))
			games.POST("/upload", container.GameHandler.UploadGame)
			games.PUT("/:id", container.GameHandler.UpdateGame)
		}

		categories := protected.Group("/categories")
		{
			categories.Use(container.RoleMiddleware.Require("ADMIN"))
			categories.POST("", container.GameHandler.CreateCategory)
			categories.PUT("/:id", container.GameHandler.UpdateCategory)
			categories.DELETE("/:id", container.GameHandler.DeleteCategory)
		}

		posts := protected.Group("/posts")
		{
			posts.GET("/mine", container.PostHandler.ListMine)
			posts.GET("/:id", container.PostHandler.GetByID)
			posts.POST("", container.PostHandler.Create)
			posts.PUT("/:id", container.PostHandler.Update)
			posts.DELETE("/:id", container.PostHandler.Delete)
		}

		admin := protected.Group("/admin")
		admin.Use(container.RoleMiddleware.Require("ADMIN"))
		{
			admin.GET("/dashboard", container.AdminHandler.Dashboard)
			admin.GET("/users", container.AdminHandler.ListUsers)
			admin.GET("/games", container.AdminHandler.ListGames)
			admin.POST("/games/:id/approve", container.AdminHandler.ApproveGame)
			admin.POST("/games/:id/reject", container.AdminHandler.RejectGame)
		}
	}

	return router
}
