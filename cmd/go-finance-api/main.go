package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/megajandrox/go-finance-api/pkg/routerapi"
)

func main() {

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		routerapi.QuoteRoutes(v1)
		routerapi.IndexRoutes(v1)

	}

	// Start the HTTP server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
