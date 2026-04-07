package app

import (
	"voca-store/internal/handler/address"
	"voca-store/internal/handler/auth"
	"voca-store/internal/handler/cart"
	"voca-store/internal/handler/category"
	"voca-store/internal/handler/chat"
	"voca-store/internal/handler/checkout"
	"voca-store/internal/handler/coupon"
	"voca-store/internal/handler/product"
	"voca-store/internal/handler/profile"
	"voca-store/internal/handler/system"
	"voca-store/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, h map[string]interface{}, db *gorm.DB) {
	category := h["category"].(*category.CategoryHandler)
	product := h["product"].(*product.ProductHandler)
	auth := h["auth"].(*auth.AuthHandler)
	profile := h["profile"].(*profile.ProfileHandler)
	cart := h["cart"].(*cart.CartHandler)
	checkouts := h["checkout"].(*checkout.CheckoutHandler)
	midtrans := h["midtrans"].(*checkout.MidtransHandler)
	coupon := h["coupon"].(*coupon.CouponHandler)
	address := h["address"].(*address.AddressHandler)
	system := h["system"].(*system.SystemHandler)
	chatHandler := h["chat"].(*chat.ChatHandler)

	// PUBLIC ROUTES
	r.GET("/", product.GetAllProducts)
	r.GET("/password", system.GetNewSecret)

	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	r.POST("/forgot-password", auth.ForgotPassword)
	r.POST("/verify-otp", auth.VerifyOTP)
	r.POST("/refresh", auth.RefreshToken)

	r.GET("/category", category.GetAllCategory)
	r.GET("/category/:slug", category.GetCategoryBySlug)
	r.GET("/product/:slug", product.GetProductBySlug)
	r.GET("/products", product.GetAllProducts)

	r.GET("/coupons", coupon.GetCoupons)
	r.POST("/midtrans/webhook", midtrans.MidtransWebhook)

	// PROTECTED ROUTES
	api := r.Group("/api")
	api.Use(middleware.JWTAuth(db))

	// Profile
	api.GET("/profile", profile.GetProfile)
	api.PUT("/profile", profile.UpdateProfile)
	api.PUT("/change-password", profile.ChangePassword)

	// Cart
	api.GET("/cart", cart.ViewCart)
	api.POST("/cart/items", cart.AddToCart)
	api.DELETE("/cart/items/:id", cart.RemoveCartItem)
	api.DELETE("/cart/items", cart.RemoveCartItemMany)
	api.DELETE("/cart/items/all", cart.ClearCart)

	// Checkout
	api.POST("/checkout", checkouts.CreateCheckout)
	api.GET("/checkout/me", checkouts.GetMyCheckouts)
	api.GET("/checkout/:uid", checkouts.GetCheckoutByUID)
	api.DELETE("/checkout/:uid", checkouts.DeleteCheckout)

	// Address
	api.POST("/addresses", address.CreateAddress)
	api.GET("/addresses", address.GetMyAddresses)
	api.GET("/addresses/:uid", address.GetAddressByUID)
	api.PUT("/addresses/:uid", address.UpdateAddress)
	api.DELETE("/addresses/:uid", address.DeleteAddress)

	// Coupons
	api.POST("/coupons/:id/claim", coupon.ClaimCoupon)
	api.GET("/coupons/me", coupon.GetMyCoupons)
	api.DELETE("/coupons/:id/remove", coupon.RemoveCoupon)

	api.POST("/logout", auth.Logout)

	// Chat
	chats := api.Group("/chat")
	chats.GET("/status", chatHandler.GetActiveSession)
	chats.POST("/requests", chatHandler.CreateChatRequest)
	chats.GET("/sessions/:session_uid", chatHandler.GetSessionByUID)
	chats.GET("/sessions/:session_uid/messages", chatHandler.GetChatHistory)
	chats.POST("/sessions/:session_uid/read", chatHandler.MarkMessagesRead)
	chats.POST("/sessions/:session_uid/close", chatHandler.CloseSession)

	// ADMIN ROUTES
	admin := api.Group("/admin")
	admin.Use(middleware.AdminOnly())

	admin.POST("/category", category.CreateCategory)
	admin.PUT("/category/:id", category.UpdateCategory)
	admin.DELETE("/category/:id", category.DeleteCategory)

	admin.POST("/products", product.CreateProduct)
	admin.PUT("/products/:id", product.UpdateProduct)
	admin.DELETE("/products/:id", product.DeleteProduct)
	admin.DELETE("/products", product.DeleteAllProducts)
	admin.DELETE("/products/assets", product.DeleteAllProductImages)

	admin.POST("/coupons", coupon.CreateCoupon)
	admin.PUT("/coupon/:id", coupon.UpdateCoupon)
	admin.DELETE("/coupon/:id", coupon.DeleteCoupon)

	admin.GET("/checkout", checkouts.GetCheckoutList)
	admin.PATCH("/checkout/:id/approve", checkouts.ApproveCheckout)
	admin.PATCH("/checkout/:id/reject", checkouts.RejectCheckout)

	admin.GET("/chat/pending", chatHandler.GetPendingChatRequests)
	admin.GET("/chat/sessions", chatHandler.GetAllSessions)
	admin.POST("/chat/requests/:session_uid/accept", chatHandler.AcceptChatRequest)

	// SYSTEM ROUTES
	sys := r.Group("/system")
	sys.Use(middleware.SystemAuth())

	sys.POST("/reset", system.ResetDatabase)
	sys.POST("/reset/product", system.ResetDatabaseWithProducts)
	sys.POST("/reset/catalog", system.ResetDatabasePreserveCatalog)
	sys.POST("/migrate", system.Migrate)
	sys.DELETE("/reset/assets", system.DeleteAllCloudinaryAssets)
	sys.POST("/redis", system.ResetRedis)

	seed := sys.Group("/seed")
	seed.POST("/assets", system.SeedProductsFromAssets)
	seed.GET("/roles", system.SeedRoles)
	seed.GET("/admin", system.SeedAdmin)
	seed.GET("/users", system.SeedUsers)
	seed.POST("/products", system.SeedProducts)
	seed.GET("/coupons", system.SeedCoupons)
	seed.POST("/sync", system.SyncAssetProducts)
	seed.POST("/all", system.SeedAll)
	seed.POST("/all-product", system.SeedAllWithProducts)

	// Websocket
	ws := r.Group("/ws")
	ws.Use(middleware.WebSocketAuth(db))
	{
		ws.GET("/chat/:session_uid", chatHandler.WebSocketHandler)
	}
}