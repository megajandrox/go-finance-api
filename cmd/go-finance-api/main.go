package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/internal/middleweare"
	"github.com/megajandrox/go-finance-api/pkg/repository"
	"github.com/megajandrox/go-finance-api/pkg/routerapi"
)

func main() {

	router := gin.Default()
	db := middleweare.InitializeDatabase()
	positionRepo := repository.NewPositionRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	v1 := router.Group("/api/v1")
	{
		routerapi.QuoteRoutes(v1)
		routerapi.IndexRoutes(v1)
		routerapi.AssetRoutes(v1, assetRepo, positionRepo)
	}

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
