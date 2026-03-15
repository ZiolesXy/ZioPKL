package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"voca-plane/config"
	"voca-plane/internal/handler"
	"voca-plane/internal/repository"
	"voca-plane/internal/seeders"
	"voca-plane/internal/service"
	"voca-plane/middleware"
	"voca-plane/pkg/helper"
	"voca-plane/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("⚙️  Starting Server...")
	godotenv.Load()
	cfg := config.LoadConfig()
    if cfg.GinMode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

	midtransClient := helper.NewMidTransClient(
		cfg.MidtransServerKey,
		cfg.MidtransIsProd,
	)
	
	db := config.PostgresDatabase(cfg)

	seeders.InitSeeders(db)

	userRepo := repository.NewUserRepository(db)
	flightRepo := repository.NewFlightRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	promoRepo := repository.NewPromoRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	airlineRepo := repository.NewAirlineRepository(db)
	airportRepo := repository.NewAirportRepository(db)
	systemRepo := repository.NewSystemRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.AccessTokenExpiry, cfg.RefreshTokenExpiry)
	flightService := service.NewFlightService(flightRepo)
	transactionService := service.NewTransactionService(transactionRepo, flightRepo, promoRepo, db, midtransClient)
	userService := service.NewUserService(userRepo)
	adminService := service.NewAdminService(adminRepo, userRepo, flightRepo, airlineRepo, airportRepo, promoRepo, db)
	systemService := service.NewSystemService(db, systemRepo)

	authHandler := handler.NewAuthHandler(authService)
	flightHandler := handler.NewFlightHandler(flightService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	userHandler := handler.NewUserHandler(userService)
	adminHandler := handler.NewAdminHandler(adminService)
	systemHandler := handler.NewSystemHandler(systemService)

	r := gin.Default()
	r.Use(middleware.Logger())

	routes.SetUpRoutes(r, authHandler, flightHandler, transactionHandler, userHandler, adminHandler, systemHandler, userRepo, cfg.JWTSecret, cfg.AllowedOrigins, cfg.AppPassword)

	srv := &http.Server{
		Addr: ":" + cfg.AppPort,
		Handler: r,
	}

	go func(){
		log.Printf("🚀 Server starting on port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("❌ Server forced to shutdown: ", err)
	}
	log.Println("✅ Server exiting")
}