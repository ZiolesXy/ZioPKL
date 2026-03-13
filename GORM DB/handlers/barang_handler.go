package handlers

import (
	"main/database"
	"main/models"
	"main/request"
	"main/respons"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateBarang(c *gin.Context) {
	var req request.BarangPost

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
		request.NewJsonResponse("Invalid request", err.Error()))
		return
	}

	barang := models.Barang{Name: req.Name}

	if err := database.DB.Create(&barang).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
		request.NewJsonResponse("Error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, request.NewJsonResponse("Succes", respons.Barang{
		ID: barang.ID,
		Name: barang.Name,
	}))
}

func GetBarangs(c *gin.Context) {
	var barangs []models.Barang

	if err := database.DB.Order("id ASC").Find(&barangs).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed get barangs", nil),
		)
		return
	}

	barangResponse := []respons.Barang{}

	for _, barang := range barangs {
		barangResponse = append(barangResponse, respons.Barang{
			ID: barang.ID,
			Name: barang.Name,
		})
	}


	entries := respons.NewEntries(barangResponse)
	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", entries),
	)
}

func GetBarangByID(c *gin.Context) {
	id := c.Param("id")

	var barang models.Barang

	if err := database.DB.Preload("Users").First(&barang, id).Error; err != nil {
		c.JSON(
			http.StatusNotFound,
			respons.NewJsonResponse("Barang not found", nil),
		)
		return
	}

	userResps := []respons.User{}
	for _, user := range barang.Users {
		userResps = append(userResps, respons.User{
			ID: user.ID,
			Name: user.Name,
		})
	}

	barangResp := respons.Barang {
		ID: barang.ID,
		Name: barang.Name,
	}

	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", barangResp),
	)
}

func UpdateBarang(c *gin.Context) {
	id := c.Param("id")

	var barang models.Barang
	if err := database.DB.First(&barang, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
		request.NewJsonResponse("Barang not found", nil))
		return
	}

	var req request.BarangPut
	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest,
		request.NewJsonResponse("Error", err))
		return
	}

	if req.Name != nil {
		barang.Name = *req.Name
	}

	if err := database.DB.Save(&barang).Error; err != nil {
		c.JSON(http.StatusBadRequest,
		request.NewJsonResponse("Error", err.Error()))
		return
	}

	c.JSON(http.StatusOK,
	request.NewJsonResponse("Barang updated", respons.Barang{
		ID: barang.ID,
		Name: barang.Name,
	}))
}

func DelBarang(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Barang{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			request.NewJsonResponse("Barang not found", nil))
		return
	}

	c.JSON(http.StatusOK,
		request.NewJsonResponse("Barang deleted", nil))
}