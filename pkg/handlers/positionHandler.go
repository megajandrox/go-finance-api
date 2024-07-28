package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/dto"
	"github.com/megajandrox/go-finance-api/pkg/models"
	"github.com/megajandrox/go-finance-api/pkg/repository"
)

func BuyPosition(repo repository.PositionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var addPosition dto.BuyPosition
		if err := c.ShouldBindJSON(&addPosition); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//TODO falta validar que exista la accion y que la cantidad sea positiva
		position, errNew := models.NewPosition(addPosition.Symbol, addPosition.Price, addPosition.Quantity, addPosition.MarketType)
		if errNew != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": errNew.Error()})
			return
		}
		if err := repo.Create(position); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Position created successfully", "position": position})
	}
}

func SellPosition(repo repository.PositionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var sellPosition dto.SellPosition
		if err := c.ShouldBindJSON(&sellPosition); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		idStr := c.Param("id")
		idInt, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		// Convertir el entero a uint
		id := uint(idInt)
		position, errGetByID := repo.GetByID(id)
		if errGetByID != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": errGetByID.Error()})
			return
		}
		if sellPosition.Quantity >= position.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect quantity, should be greather than previous one."})
			return
		}
		position.PositionType = models.Sold
		position.ExitTime = time.Now()
		position.Quantity = position.Quantity - sellPosition.Quantity
		position.Balance = (sellPosition.Price - position.EntryPrice) * float64(sellPosition.Quantity)
		if err := repo.Update(position); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Position created successfully", "position": position})
	}
}
