package handlers

import (
	"main/database"
	"main/models"
	"main/respons"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserBarangs(c *gin.Context) {
	var users []models.User

	if err := database.DB.Preload("Role").Preload("Barangs").Order("id ASC").Find(&users).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed get users barangs", nil),
		)
		return
	}

	userResponses := []respons.UserWithBarang{}

	for _, user := range users {
		barangResp := []respons.Barang{}
		for _, barang := range user.Barangs {
			barangResp = append(barangResp, respons.Barang{
				ID:   barang.ID,
				Name: barang.Name,
			})
		}

		roleResp := respons.Role{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		}

		userResponses = append(userResponses, respons.UserWithBarang{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Role:    roleResp,
			Barangs: barangResp,
		})
	}

	entries := respons.NewEntries(userResponses)
	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", entries),
	)
}

func GetUserBarangByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.Preload("Role").Preload("Barangs").First(&user, id).Error; err != nil {
		c.JSON(
			http.StatusNotFound,
			respons.NewJsonResponse("User not found", nil),
		)
		return
	}

	barangResp := []respons.Barang{}
	for _, barang := range user.Barangs {
		barangResp = append(barangResp, respons.Barang{
			ID:   barang.ID,
			Name: barang.Name,
		})
	}

	roleResp := respons.Role{
		ID:   user.Role.ID,
		Name: user.Role.Name,
	}

	userResp := respons.UserWithBarang{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Role:    roleResp,
		Barangs: barangResp,
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Success", userResp),
	)
}

func GetUserBarangPivot(c *gin.Context) {
	var results []respons.OwnersipDTO
	limit := 20
	offset := 0

	// Query builder GORM (Tanpa Exec)
	query := database.DB.Table("user_barangs").
		Select("users.name AS user_name, barangs.name AS barang_name").
		Joins("JOIN users ON users.id = user_barangs.user_id").
		Joins("JOIN barangs ON barangs.id = user_barangs.barang_id").
		Order("users.name ASC").
		Limit(limit).
		Offset(offset)

	// Eksekusi query dan masukkan hasilnya ke slice results
	if err := query.Scan(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, respons.NewJsonResponse("Gagal ambil data", err.Error()))
		return
	}

	entries := respons.NewEntries(results)
	c.JSON(http.StatusOK, respons.NewJsonResponse("Success", entries))
}

func AssignBarang(c *gin.Context) {
	userID := c.Param("user_id")
	barangID := c.Param("barang_id")

	var user models.User
	var barang models.Barang

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, respons.NewJsonResponse("User not found", nil))
		return
	}

	if err := database.DB.First(&barang, barangID).Error; err != nil {
		c.JSON(http.StatusNotFound, respons.NewJsonResponse("Barang not found", nil))
		return
	}

	if err := database.DB.Model(&user).Association("Barangs").Append(&barang); err != nil {
		c.JSON(http.StatusInternalServerError, respons.NewJsonResponse("Failed to assign barang", nil))
		return
	}

	result := map[string]interface{}{
		"user_name":   user.Name,
		"barang_name": barang.Name,
	}

	c.JSON(http.StatusOK, respons.NewJsonResponse("Barang added", result))
}

func RemoveBarang(c *gin.Context) {
	userID := c.Param("id")
	barangID := c.Param("barang_id")

	var user models.User
	var barang models.Barang

	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"Error": "User tidak ditemukan"})
		return
	}

	if err := database.DB.First(&barang, barangID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	if err := database.DB.Model(&user).Association("Barangs").Delete(&barang); err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"Messege": "Barang berhasil dilepas dari user"})

}
