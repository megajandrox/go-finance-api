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

func AssetRoutes(v1 *gin.RouterGroup, repo repository.AssetRepository, repo2 repository.PositionRepository) {
	assetGroup := v1.Group("/assets")
	{
		assetGroup.GET("/", handlers.GetAllAssets(repo))
		assetGroup.POST("/", handlers.CreateAsset(repo))
		assetGroup.PUT("/:id", handlers.UpdateAsset(repo))
		positionGroup := assetGroup.Group("/:id/positions")
		{
			positionGroup.POST("/", handlers.BuyPosition(repo2))
			positionGroup.PUT("/:idPosition", handlers.SellPosition(repo2))
		}
	}
}
