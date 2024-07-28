package routerapi

import (
	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/handlers"
	"github.com/megajandrox/go-finance-api/pkg/repository"
)

func QuoteRoutes(v1 *gin.RouterGroup) {
	quoteGroup := v1.Group("/quote")
	{

		quoteGroup.GET("/:symbol", handlers.GetQuote)

	}
}

func IndexRoutes(v1 *gin.RouterGroup) {
	quoteGroup := v1.Group("/index")
	{

		quoteGroup.GET("/:symbol", handlers.GetIndex)

	}
}

func PositionRoutes(v1 *gin.RouterGroup, repo repository.PositionRepository) {
	quoteGroup := v1.Group("/position")
	{

		quoteGroup.POST("/", handlers.BuyPosition(repo))
		quoteGroup.PUT("/:id", handlers.SellPosition(repo))

	}
}
