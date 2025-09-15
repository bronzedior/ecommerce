package main

import (
	"product/cmd/product/handler"
	"product/config"
	"product/infrastructure/log"
	"product/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	productHandler := handler.NewProductHandler()

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *productHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
