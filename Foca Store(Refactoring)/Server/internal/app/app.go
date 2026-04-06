package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"voca-store/internal/database"
	"voca-store/internal/handler/address"
	"voca-store/internal/handler/auth"
	"voca-store/internal/handler/cart"
	"voca-store/internal/handler/category"
	"voca-store/internal/handler/checkout"
	"voca-store/internal/handler/coupon"
	"voca-store/internal/handler/product"
	"voca-store/internal/handler/profile"
	"voca-store/internal/handler/system"
	"voca-store/internal/helper"
	"voca-store/internal/repository"
	"voca-store/internal/service"
	"voca-store/seeders"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
    Redis  *redis.Client
}

func New() *App {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Init Cloudinary
	if err := helper.InitCloudinary(); err != nil {
		log.Println("Cloudinary Warning: ",err)
	}

	// Init DB & REDIS
	db := database.InitDB()
	rdb := database.InitRedis()

	// Run Migrations
	if err := seeders.MigrateAll(db); err != nil {
		panic(err)
	}

	// Repository
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(rdb)
	cartRepo := repository.NewCartRepository(db)
	checkoutRepo := repository.NewCheckoutRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	addressRepo := repository.NewAddressRepository(db)

	// Service
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	authService := service.NewAuthService(userRepo, authRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	checkoutService := service.NewCheckoutService(checkoutRepo, userRepo)
	couponService := service.NewCouponService(couponRepo)
	addressService := service.NewAddressService(addressRepo)
	systemService := service.NewSystemService(db, rdb)

	// Handler
	handlers := map[string]interface{}{
		"category": category.NewCategoryHandler(categoryService),
		"product":  product.NewProductHandler(productService),
		"auth":     auth.NewAuthHandler(authService),
		"profile":  profile.NewProfileHandler(authService),
		"cart":     cart.NewCartHandler(cartService),
		"checkout": checkout.NewCheckoutHandler(checkoutService),
		"midtrans": checkout.NewMidtransHandler(checkoutService),
		"coupon":   coupon.NewCouponHandler(couponService),
		"address":  address.NewAddressHandler(addressService),
		"system":   system.NewSystemHandler(systemService),
	}

	r := SetUpServer()
	RegisterRoutes(r, handlers, db)
	return &App{Router: r}
}

func (a *App) Run() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr: ":" + port,
		Handler: a.Router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func ()  {
		fmt.Printf("Server running on port %s\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	database.CloseDB()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
	return nil
}