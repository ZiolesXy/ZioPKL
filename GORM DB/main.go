package main

import (
	"main/database"
	"main/handlers"
	"main/middlewares"
	"main/models"
	"main/trash"
	"net/http"
	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        
        if origin == "http://localhost:3001" {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Vary", "Origin")
            c.Header("Access-Control-Allow-Credentials", "true") 
        }

        // PERHATIKAN BARIS INI: Tambahkan 'Authorization' ke dalam list
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

        if c.Request.Method == http.MethodOptions {
            // Sangat penting agar header di atas ikut terkirim saat status No Content
            c.AbortWithStatus(http.StatusNoContent)
            return
        }

        c.Next()
    }
}

func main() {
	database.ConnedtDB()

	database.DB.AutoMigrate(
		&models.Role{},
		&models.User{},
		&models.Barang{},
		&models.RefreshToken{},
	)

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/trash", trash.SayHello)

	//public
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/refresh-token", handlers.RefreshToken)
	r.POST("/logout", handlers.Logout)

	r.POST("/dummy", handlers.RunSeeder)

	//auth
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	{
		r.GET("/", handlers.GetUserBarangs)
		auth.GET("/me", handlers.GetProfile)
		auth.GET("/profile/:id", handlers.GetUserByID)

		auth.GET("/users", handlers.GetUsers)
		auth.GET("/user/:id", handlers.GetUserBarangByID)

		auth.GET("/roles", handlers.GetRole)
		auth.GET("/role/:id", handlers.GetRoleByID)

		auth.GET("/barangs", handlers.GetBarangs)
		auth.GET("/barang/:id", handlers.GetBarangByID)

		auth.GET("/users/barangs", handlers.GetUserBarangs)
		auth.GET("/user/barang", handlers.GetUserBarangPivot)
	}

	//admin only
	admin := r.Group("/")
	admin.Use(
		middlewares.AuthMiddleware(),
		middlewares.OnlyAdmin(),
	)
	{
		admin.POST("/roles", handlers.CreateRole)
		admin.PUT("/role/:id", handlers.UpdateRole)
		admin.DELETE("/role/:id", handlers.DelRole)
		
		admin.PUT("/user/:id", handlers.UpdateUsers)
		admin.DELETE("/user/:id", handlers.DelUser)

		admin.POST("/barangs", handlers.CreateBarang)
		admin.PUT("/barang/:id", handlers.UpdateBarang)
		admin.DELETE("/barang/:id", handlers.DelBarang)

		admin.POST("/user/:user_id/barang/:barang_id", handlers.AssignBarang)
		admin.DELETE("/user/:id/barang/:barang_id", handlers.RemoveBarang)

		admin.DELETE("/logouts", handlers.LogoutAll)
	}

	r.Run(":3605")
}