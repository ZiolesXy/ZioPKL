package handlers

import (
	"main/seeders"
	"main/trash"
	"github.com/gin-gonic/gin"
)

func RunSeeder(c *gin.Context) {
	seeders.RunSeed(c)
}

func ClearSeeder(c *gin.Context) {
	trash.ClearData(c)
	c.JSON(200, gin.H{"messege": "Data berhasil terhapus"})
}