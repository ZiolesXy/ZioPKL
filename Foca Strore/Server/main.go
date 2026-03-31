package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"voca-store/database"
	"voca-store/handlers"
	"voca-store/helper"
	"voca-store/middleware"
	"voca-store/seeders"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Cloudinary
	log.Println("starting cloudinary initialiized")
	if err := helper.InitCloudinary(); err != nil {
		log.Println("Warning: Cloudinary not initialized:", err)
	} else {
		log.Println("Cloudinary initialized successfully")
	}

	// Initialize database
	db := database.InitDB()
	rdb := database.InitRedis()

	if err := seeders.MigrateAll(db); err != nil {
		panic("failed migrate")
	}

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//cors set
	originEnv := os.Getenv("ALLOW_ORIGINS")
	allowedOrigins := []string{"http://localhost:3000"}
	if originEnv != "" {
		allowedOrigins = strings.Split(originEnv, ",")
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	authHandler := handlers.NewAuthHandler(db)
	r.GET("/", handlers.GetAllProducts(db))
	r.GET("/password", handlers.GetNewSecret)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/forgot-password", authHandler.ForgotPassword)
	r.POST("/verify-otp", authHandler.VerifyOTP)
	r.POST("/refresh", authHandler.RefreshToken)
	r.GET("/category/:slug", handlers.GetCategoryBySlug(db))
	r.GET("/category", handlers.GetAllCategory(db))
	r.GET("/product/:slug", handlers.GetProductBySlug(db))
	r.GET("/products", handlers.GetAllProducts(db))
	r.GET("/coupons", handlers.GetCoupons(db))
	r.POST("/midtrans/webhook", handlers.MidtransWebhook(db))

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTAuth(db))
	{
		// User routes
		protected.GET("/profile", handlers.GetProfile(db))
		protected.PUT("/profile", handlers.UpdateProfile(db))
		protected.PUT("/change-password", handlers.ChangePassword(db))
		protected.GET("/cart", handlers.ViewCart(db))
		protected.POST("/cart/items", handlers.AddToCart(db))
		protected.DELETE("/cart/items/:id", handlers.RemoveCartItem(db))
		protected.DELETE("/cart/items", handlers.RemoveCartItemMany(db))
		protected.DELETE("/cart/items/all", handlers.ClearCart(db))
		protected.POST("/checkout", handlers.Checkout(db))
		protected.GET("/checkout/me", handlers.GetMyCheckout(db))
		protected.GET("/checkout/:uid", handlers.GetCheckoutByUID(db))
		protected.DELETE("/checkout/:uid", handlers.DeleteMyCheckout(db))
		protected.POST("/logout", authHandler.Logout)

		// Address routes
		protected.POST("/addresses", handlers.CreateAddress(db))
		protected.GET("/addresses", handlers.GetMyAddresses(db))
		protected.GET("/addresses/:uid", handlers.GetAddressByUID(db))
		protected.PUT("/addresses/:uid", handlers.UpdateAddress(db))
		protected.DELETE("/addresses/:uid", handlers.DeleteAddress(db))

		// Coupon user routes
		protected.POST("/coupons/:id/claim", handlers.ClaimCoupon(db))
		protected.GET("/coupons/me", handlers.GetMyCoupons(db))
		protected.DELETE("/coupons/:id/remove", handlers.RemoveCoupon(db))

		// Admin routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/category", handlers.CreateCategory(db))
			admin.PUT("/category/:id", handlers.UpdateCategory(db))
			admin.DELETE("/category/:id", handlers.DeleteCategory(db))
			admin.POST("/products", handlers.CreateProduct(db))
			admin.PUT("/products/:id", handlers.UpdateProduct(db))
			admin.DELETE("/products/:id", handlers.DeleteProduct(db))
			admin.DELETE("/products", handlers.DeleteAllProducts(db))
			admin.DELETE("/products/assets", handlers.DeleteAllProductImages(db))
			admin.POST("/coupons", handlers.CreateCoupon(db))
			admin.PUT("/coupon/:id", handlers.UpdateCoupon(db))
			admin.DELETE("/coupon/:id", handlers.DeleteCoupon(db))
			admin.GET("/checkout", handlers.GetCheckout(db))
			admin.PATCH("/checkout/:id/approve", handlers.ApproveCheckout(db))
			admin.PATCH("/checkout/:id/reject", handlers.RejectCheckout(db))
		}
	}

	system := r.Group("/system")
	system.Use(middleware.SystemAuth())
	{
		system.POST("/reset", handlers.ResetDatabaseHandler(db, rdb))
		system.POST("/reset/product", handlers.ResetDatabaseWithProductsHandler(db, rdb))
		system.POST("/reset/catalog", handlers.ResetDatabasePreserveProductsAndCategoriesHandler(db, rdb))
		system.POST("/migrate", handlers.MigrateHandler(db))
		system.DELETE("/reset/assets", handlers.DeleteAllCloudinaryAssets())

		system.POST("/redis", handlers.ResetRedis(rdb))
		// Seeder endpoint
		seed := system.Group("/seed")
		{
			seed.GET("/assets", handlers.SeedProductsFromAssetsHandler(db))
			seed.GET("/roles", handlers.SeedRoleHandler(db))
			seed.GET("/admin", handlers.SeedAdminHandler(db))
			seed.GET("/users", handlers.SeedUsersHandler(db))
			seed.GET("/products", handlers.SeedProductsHandler(db))
			seed.GET("/coupons", handlers.SeedCouponHandler(db))
			seed.PUT("/sync", handlers.SyncAssetProductsHandler(db))
			seed.GET("/all", handlers.SeedAllHandler(db))
			seed.GET("/all-product", handlers.SeedAllWithProductnonAssetsHandler(db))
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server running on port %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}