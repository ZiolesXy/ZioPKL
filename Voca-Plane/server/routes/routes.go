package routes

import (
	"net/http"
	"voca-plane/internal/handler"
	"voca-plane/internal/repository"
	"voca-plane/middleware"
	globalResp "voca-plane/pkg/response"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine,
	authHandler *handler.AuthHandler,
	flightHandler *handler.FlightHandler,
	transactionHandler *handler.TransactionHandler,
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	systemHandler *handler.SystemHandler,
	userRepo repository.UserRepository,
	jwtSecret string,
	allowedOrigins string,
	appPassword string) {
	r.Use(middleware.CORS(allowedOrigins))

	v1 := r.Group("/api/v1")
	{
		v1.GET("", func(c *gin.Context) {
			globalResp.Success(c, http.StatusOK, "welcome to the voca plane", nil)
		})
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		v1.GET("/info", userHandler.GetDeviceInfo)

		// System Routes
		sys := v1.Group("/system")
		sys.Use(middleware.JWTAuth(jwtSecret, userRepo))
		sys.Use(middleware.RequireSuperAdmin())
		sys.Use(middleware.AppPassword(appPassword))
		{
			sys.POST("/seed", systemHandler.Seed)
		}

		// Public Auth
		v1.POST("/auth/register", authHandler.Register)
		v1.POST("/auth/login", authHandler.Login)
		v1.POST("/auth/refresh", authHandler.RefreshToken)
		v1.POST("/transactions/midtrans/callback", transactionHandler.MidtransCallback)

		// Public Flight Search
		v1.GET("/flights", flightHandler.GetAll)
		v1.GET("/flights/all", flightHandler.GetAllFull)
		v1.GET("/flights/search", flightHandler.Search)
		v1.GET("/flights/:id", flightHandler.GetByID)
		v1.GET("/flights/:id/seats", flightHandler.GetSeats)
		v1.GET("/airports", adminHandler.GetAirports)
		v1.GET("airlines", adminHandler.GetAirlines)

		// Protected User Routes
		userProtected := v1.Group("")
		userProtected.Use(middleware.JWTAuth(jwtSecret, userRepo))
		{
			userProtected.GET("/user/profile", userHandler.GetProfile)
			userProtected.PATCH("/user/profile", userHandler.UpdateProfile)

			userProtected.GET("/transactions", transactionHandler.GetListAll)
			userProtected.GET("/transactions/:code", transactionHandler.GetByCode)
			userProtected.POST("/transactions", transactionHandler.Create)
			// userProtected.PATCH("/transactions/:code/pay", transactionHandler.Pay)
			userProtected.DELETE("/transactions/:code", transactionHandler.Cancel)
		}

		// Protected Admin Routes
		adminProtected := v1.Group("/admin")
		adminProtected.Use(middleware.JWTAuth(jwtSecret, userRepo))
		adminProtected.Use(middleware.RequireAdmin())
		{
			adminProtected.GET("/dashboard", adminHandler.GetDashboard)

			adminProtected.GET("/users", adminHandler.GetUsers)
			adminProtected.PATCH("/users/:id/role", adminHandler.UpdateUserRole)
			adminProtected.DELETE("/users/:id", adminHandler.DeleteUser)
			adminProtected.PATCH("/users/:id/restore", adminHandler.RestoreUser)
			adminProtected.PATCH("/users/:id/ban", adminHandler.BanUser)
			adminProtected.PATCH("/users/:id/unban", adminHandler.UnbanUser)

			adminProtected.GET("/transactions", adminHandler.GetTransactions)

			adminProtected.GET("/flights", adminHandler.GetFlights)
			adminProtected.POST("/flights", adminHandler.CreateFlight)
			adminProtected.PUT("/flights/:id", adminHandler.UpdateFlight)
			adminProtected.DELETE("/flights/:id", adminHandler.DeleteFlight)

			adminProtected.GET("/airlines", adminHandler.GetAirlines)
			adminProtected.POST("/airlines", adminHandler.CreateAirline)
			adminProtected.PUT("/airlines/:id", adminHandler.UpdateAirline)
			adminProtected.DELETE("/airlines/:id", adminHandler.DeleteAirline)

			adminProtected.GET("/airports", adminHandler.GetAirports)
			adminProtected.POST("/airports", adminHandler.CreateAirport)
			adminProtected.PUT("/airports/:id", adminHandler.UpdateAirport)
			adminProtected.DELETE("/airports/:id", adminHandler.DeleteAirport)

			adminProtected.GET("/promos", adminHandler.GetPromos)
			adminProtected.POST("/promos", adminHandler.CreatePromo)
			adminProtected.PUT("/promos/:id", adminHandler.UpdatePromo)
			adminProtected.DELETE("/promos/:id", adminHandler.DeletePromo)
		}
	}
}
