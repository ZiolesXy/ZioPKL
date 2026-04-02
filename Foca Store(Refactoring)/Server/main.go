package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"voca-store/internal/database"
	"voca-store/internal/helper"
	"voca-store/internal/middleware"
	"voca-store/seeders"

	// New Repositories
	"voca-store/internal/repository"

	// New Services
	"voca-store/internal/service"

	// New Handlers
	"voca-store/internal/handler/address"
	"voca-store/internal/handler/auth"
	"voca-store/internal/handler/cart"
	"voca-store/internal/handler/category"
	"voca-store/internal/handler/checkout"
	"voca-store/internal/handler/coupon"
	"voca-store/internal/handler/product"
	"voca-store/internal/handler/profile"
	"voca-store/internal/handler/system"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// 2. Initialize Cloudinary
	if err := helper.InitCloudinary(); err != nil {
		log.Println("Warning: Cloudinary not initialized:", err)
	} else {
		log.Println("Cloudinary initialized successfully")
	}

	// 3. Initialize database
	db := database.InitDB()
	rdb := database.InitRedis()

	// 4. Initial Migration
	if err := seeders.MigrateAll(db); err != nil {
		panic("failed migrate: " + err.Error())
	}

	// 5. Initialize Repositories
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(rdb)
	cartRepo := repository.NewCartRepository(db)
	checkoutRepo := repository.NewCheckoutRepository(db)
	couponRepo := repository.NewCouponRepository(db)
	addressRepo := repository.NewAddressRepository(db)

	// 6. Initialize Services
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	authService := service.NewAuthService(userRepo, authRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	checkoutService := service.NewCheckoutService(checkoutRepo, userRepo)
	couponService := service.NewCouponService(couponRepo)
	addressService := service.NewAddressService(addressRepo)
	systemService := service.NewSystemService(db, rdb)

	// 7. Initialize Handlers
	categoryHandler := category.NewCategoryHandler(categoryService)
	productHandler := product.NewProductHandler(productService)
	authHandler := auth.NewAuthHandler(authService)
	profileHandler := profile.NewProfileHandler(authService)
	cartHandler := cart.NewCartHandler(cartService)
	checkoutHandler := checkout.NewCheckoutHandler(checkoutService)
	midtransHandler := checkout.NewMidtransHandler(checkoutService)
	couponHandler := coupon.NewCouponHandler(couponService)
	addressHandler := address.NewAddressHandler(addressService)
	systemHandler := system.NewSystemHandler(systemService)

	// 8. Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS configuration
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

	// ==========================================
	// ROUTES
	// ==========================================

	// Public routes
	r.GET("/", productHandler.GetAllProducts)
	r.GET("/password", systemHandler.GetNewSecret)
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/forgot-password", authHandler.ForgotPassword)
	r.POST("/verify-otp", authHandler.VerifyOTP)
	r.POST("/refresh", authHandler.RefreshToken)
	r.GET("/category/:slug", categoryHandler.GetCategoryBySlug)
	r.GET("/category", categoryHandler.GetAllCategory)
	r.GET("/product/:slug", productHandler.GetProductBySlug)
	r.GET("/products", productHandler.GetAllProducts)
	r.GET("/coupons", couponHandler.GetCoupons)
	r.POST("/midtrans/webhook", midtransHandler.MidtransWebhook)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.JWTAuth(db))
	{
		// Profile
		api.GET("/profile", profileHandler.GetProfile)
		api.PUT("/profile", profileHandler.UpdateProfile)
		api.PUT("/change-password", profileHandler.ChangePassword)

		// Cart
		api.GET("/cart", cartHandler.ViewCart)
		api.POST("/cart/items", cartHandler.AddToCart)
		api.DELETE("/cart/items/:id", cartHandler.RemoveCartItem)
		api.DELETE("/cart/items", cartHandler.RemoveCartItemMany)
		api.DELETE("/cart/items/all", cartHandler.ClearCart)

		// Checkout
		api.POST("/checkout", checkoutHandler.CreateCheckout)
		api.GET("/checkout/me", checkoutHandler.GetMyCheckouts)
		api.GET("/checkout/:uid", checkoutHandler.GetCheckoutByUID)
		api.DELETE("/checkout/:uid", checkoutHandler.DeleteCheckout)

		// Address
		api.POST("/addresses", addressHandler.CreateAddress)
		api.GET("/addresses", addressHandler.GetMyAddresses)
		api.GET("/addresses/:uid", addressHandler.GetAddressByUID)
		api.PUT("/addresses/:uid", addressHandler.UpdateAddress)
		api.DELETE("/addresses/:uid", addressHandler.DeleteAddress)

		// Coupon claims
		api.POST("/coupons/:id/claim", couponHandler.ClaimCoupon)
		api.GET("/coupons/me", couponHandler.GetMyCoupons)
		api.DELETE("/coupons/:id/remove", couponHandler.RemoveCoupon)

		api.POST("/logout", authHandler.Logout)

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.AdminOnly())
		{
			admin.POST("/category", categoryHandler.CreateCategory)
			admin.PUT("/category/:id", categoryHandler.UpdateCategory)
			admin.DELETE("/category/:id", categoryHandler.DeleteCategory)

			admin.POST("/products", productHandler.CreateProduct)
			admin.PUT("/products/:id", productHandler.UpdateProduct)
			admin.DELETE("/products/:id", productHandler.DeleteProduct)
			admin.DELETE("/products", productHandler.DeleteAllProducts)
			admin.DELETE("/products/assets", productHandler.DeleteAllProductImages)

			admin.POST("/coupons", couponHandler.CreateCoupon)
			admin.PUT("/coupon/:id", couponHandler.UpdateCoupon)
			admin.DELETE("/coupon/:id", couponHandler.DeleteCoupon)

			admin.GET("/checkout", checkoutHandler.GetCheckoutList)
			admin.PATCH("/checkout/:id/approve", checkoutHandler.ApproveCheckout)
			admin.PATCH("/checkout/:id/reject", checkoutHandler.RejectCheckout)
		}
	}

	// System & Maintenance
	sys := r.Group("/system")
	sys.Use(middleware.SystemAuth())
	{
		sys.POST("/reset", systemHandler.ResetDatabase)
		sys.POST("/reset/product", systemHandler.ResetDatabaseWithProducts)
		sys.POST("/reset/catalog", systemHandler.ResetDatabasePreserveCatalog)
		sys.POST("/migrate", systemHandler.Migrate)
		sys.DELETE("/reset/assets", systemHandler.DeleteAllCloudinaryAssets)
		sys.POST("/redis", systemHandler.ResetRedis)

		// Seeder endpoints
		seed := sys.Group("/seed")
		{
			seed.GET("/assets", systemHandler.SeedProductsFromAssets)
			seed.GET("/roles", systemHandler.SeedRoles)
			seed.GET("/admin", systemHandler.SeedAdmin)
			seed.GET("/users", systemHandler.SeedUsers)
			seed.GET("/products", systemHandler.SeedProducts)
			seed.GET("/coupons", systemHandler.SeedCoupons)
			seed.PUT("/sync", systemHandler.SyncAssetProducts)
			seed.GET("/all", systemHandler.SeedAll)
			seed.GET("/all-product", systemHandler.SeedAllWithProducts)
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
