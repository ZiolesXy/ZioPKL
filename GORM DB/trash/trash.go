package trash

import (
	"errors"
	"main/database"
	"main/helpers"
	"main/models"
	"main/request"
	"main/respons"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SayHello(c *gin.Context){
	c.JSON(200, respons.NewJsonResponse("Success", "Hallo Mok"))
}

func GetUsers(c *gin.Context) {
	var users []models.User

	if err := database.DB.Preload("Role").Preload("Barangs").Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if err := database.DB.Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		c.JSON(400, gin.H{"Error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetBarangs(c *gin.Context) {
	var barangs []models.Barang
	
	if err := database.DB.Find(&barangs).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, barangs)
}

func GetUserBarangs(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.Preload("Barangs").First(&user, userID).Error;err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetRole(c *gin.Context) {
	var role []models.Role
	if err := database.DB.Find(&role).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, &user)
}

func CreateBarang(c *gin.Context) {
	var barang models.Barang

	if err := c.ShouldBindJSON(&barang); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Create(&barang).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, barang)
}

func CreateRole(c *gin.Context) {
	var role models.Role

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Create(&role).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

func UpdateUsers(c *gin.Context) {
	id := c.Param("id")

	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, user)
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Role tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, role)
}

func UpdateBarang(c *gin.Context) {
	id := c.Param("id")
	var barang models.Barang

	if err := database.DB.First(&barang, id).Error; err != nil {
		c.JSON(400, gin.H{"Error": "Barang tidak ditemukan"})
		return
	}

	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(404, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Save(&barang).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, barang)
}

func DelUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	if err := database.DB.Model(&user).Association("Barangs").Clear(); err != nil{
		c.JSON(400, gin.H{"Error": err})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"Messege" : "User berhasil dihapus"})
}

func DelRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	
	if erro := database.DB.Delete(&role).Error; erro != nil {
		c.JSON(400, gin.H{"Error": erro.Error()})
		return
	}
	c.JSON(200, gin.H{"Messege" : "User berhasil dihapus"})
}

func DelBarang(c *gin.Context) {
	id := c.Param("id")
	var barang models.Barang

	if err := database.DB.First(&barang, id).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Barang tidak ditemukan"})
		return
	}

	if err := database.DB.Delete(&barang).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"messege": "Barang berhasil dihapus"})
}

func AssignBarang(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	var barang models.Barang
	
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&barang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.First(&barang, barang.ID).Error; err != nil{
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := database.DB.Model(&user).Association("Barangs").Append(&barang); err != nil{
		c.JSON(400, gin.H{"Error": err.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Messege": "Barang berhasil di assign"})
}

func ClearData(c *gin.Context) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}

	for _, user := range users {
		if err := database.DB.Model(&user).Association("Barangs").Clear(); err != nil {
			c.JSON(400, gin.H{"Error1": err.Error()})
			return
		}
		if err := database.DB.Unscoped().Delete(&user).Error; err != nil {
			c.JSON(400, gin.H{"Error2": err.Error()})
			return
		}
		// if err := database.DB.Unscoped().Delete(&models.Role{}).Error; err != nil {
		// 	c.JSON(400, gin.H{"Error3": err.Error()})
		// 	return
		// }
	}

	// if err := database.DB.Unscoped().Delete(&models.User{}).Error; err != nil {
	// 	c.JSON(400, gin.H{"Error2": err.Error()})
	// 	return
	// }

	// if err := database.DB.Unscoped().Delete(&models.Role{}).Error; err != nil {
	// 	c.JSON(400, gin.H{"Error3": err.Error()})
	// 	return
	// }

	// c.JSON(200, gin.H{
	// 	"messege": "All table data cleared",
	// })
}

func GetUserBarangPivot (c *gin.Context) {
	var results []map[string]interface{}

	if err := database.DB.Table("user_barangs").Select("user_id, barang_id").Order("user_id ASC").Find(&results).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed to fetch pivot data", err.Error()),
		)
		return
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", results),
	)
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("No token provided", nil),
			)
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// ✅ pastikan algoritma HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return helpers.ACCESS_SECRET, nil
		})

		// 🔥 BEDAKAN ERROR JWT
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.AbortWithStatusJSON(
					401,
					respons.NewJsonResponse("Token expired", nil),
				)
				return
			}

			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token", nil),
			)
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token", nil),
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token claims", nil),
			)
			return
		}

		// ✅ pastikan token access
		if claims["type"] != "access" {
			c.AbortWithStatusJSON(
				401,
				respons.NewJsonResponse("Invalid token type", nil),
			)
			return
		}

		// ✅ simpan data ke context
		c.Set("user_id", claims["user_id"])
		c.Set("role", claims["role"])

		c.Next()
	}
}

func CreateUser1(c *gin.Context) {
	var req request.UserPost

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
			request.NewJsonResponse("Invalid request", err.Error()))
		return
	}

	var role models.Role
	if err := database.DB.First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest,
			request.NewJsonResponse("Role not found", nil))
		return
	}

	user := models.User{
		Name:   req.Name,
		RoleID: req.RoleID,
		Role:   role,
	}

	if err := database.DB.Preload("Roles").Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
			request.NewJsonResponse("Failed Create user", nil))
		return
	}

	c.JSON(http.StatusCreated,
		request.NewJsonResponse("User created", respons.User{
			ID:   user.ID,
			Name: user.Name,
			Role: respons.Role{
				ID:   user.Role.ID,
				Name: user.Role.Name,
			},
		}))
}

func GetProfile(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Authorization header missing", nil),
		)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid authorization format", nil),
		)
		return
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return helpers.ACCESS_SECRET, nil
	})

	if err != nil || !token.Valid {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid or expired token", nil),
		)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid token claims", nil),
		)
		return
	}

	if claims["type"] != "access" {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid token type", nil),
		)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			respons.NewJsonResponse("Invalid user id", nil),
		)
		return
	}

	userID := uint(userIDFloat)

	var user models.User
	if err := database.DB.Preload("Role").Preload("Barangs").First(&user, userID).Error; err != nil {
		c.JSON(
			http.StatusNotFound,
			respons.NewJsonResponse("User not found", nil),
		)
		return
	}

	userResp := respons.Profile{
		ID: userID,
		Name: user.Name,
		Email: user.Email,
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Success", userResp),
	)
}