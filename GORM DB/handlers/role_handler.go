package handlers

import (
	"main/database"
	"main/models"
	"main/request"
	"main/respons"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRole(c *gin.Context) {
	var req request.BarangPost

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest,
		request.NewJsonResponse("Invalid request", err.Error()))
		return
	}

	role := models.Role{
		Name: req.Name,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError,
		request.NewJsonResponse("Failed created role", nil))
		return
	}

	c.JSON(http.StatusCreated,
	request.NewJsonResponse("Role created", respons.Barang{
		ID: role.ID,
		Name: role.Name,
	}))
}

func GetRole(c *gin.Context) {
	var roles []models.Role

	if err := database.DB.Order("id ASC").Find(&roles).Error; err != nil {
		c.JSON(
			http.StatusInternalServerError,
			respons.NewJsonResponse("Failed get roles", nil),
		)
		return
	}

	roleResponses := []respons.Role{}

	for _, role := range roles {
		roleResponses = append(roleResponses, respons.Role{
			ID: role.ID,
			Name: role.Name,
		})
	}

	entries := respons.NewEntries(roleResponses)
	c.JSON(
		http.StatusOK,
		respons.NewJsonResponse("Succes", entries),
	)
}

func GetRoleByID(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.First(&role, id).Order("id ASC").Error; err != nil{
		c.JSON(
			404,
			respons.NewJsonResponse("Role not found", nil),
		)
		return
	}

	roleResp := respons.Role{
		ID: role.ID,
		Name: role.Name,
	}

	c.JSON(http.StatusOK, respons.NewJsonResponse("Succes", roleResp))
}

func UpdateRole(c *gin.Context)  {
	id := c.Param("id")

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
		request.NewJsonResponse("Role not found", nil))
		return
	}

	var req request.RolePut
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, request.NewJsonResponse("Error", nil))
		return
	}

	if req.Name != nil {
		role.Name = *req.Name
	}

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusBadRequest, request.NewJsonResponse("Error", nil))
		return
	}

	c.JSON(http.StatusOK,
	request.NewJsonResponse("Role update", respons.Role{
		ID: role.ID,
		Name: role.Name,
	}))
}

func DelRole(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Role{}, id).Error; err != nil {
		c.JSON(http.StatusNotFound,
			request.NewJsonResponse("Role not found", nil))
		return
	}

	c.JSON(http.StatusOK,
		request.NewJsonResponse("Role deleted", nil))
}