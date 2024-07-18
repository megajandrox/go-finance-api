package routerapi

import (
	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/services"
)

func QuoteRoutes(v1 *gin.RouterGroup) {
	quoteGroup := v1.Group("/quote")
	{

		quoteGroup.GET("/:symbol", services.GetQuote)

	}
}

func IndexRoutes(v1 *gin.RouterGroup) {
	quoteGroup := v1.Group("/index")
	{

		quoteGroup.GET("/:symbol", services.GetIndex)

	}
}
