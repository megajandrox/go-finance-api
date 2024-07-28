package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/models"
	"github.com/megajandrox/go-finance-api/pkg/repository"
)

func GetAllAssets(repo repository.AssetRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		assets, err := repo.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"assets": assets})
	}
}

func CreateAsset(repo repository.AssetRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var asset models.Asset
		if err := c.ShouldBindJSON(&asset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.Create(&asset); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Asset created successfully", "asset": asset})
	}
}

func UpdateAsset(repo repository.AssetRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var asset models.Asset
		if err := c.ShouldBindJSON(&asset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := repo.UpdateByID(uint(id), &asset); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Asset updated successfully", "asset": asset})
	}
}
