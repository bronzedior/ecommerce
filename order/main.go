package main

import (
	"order/cmd/order/handler"
	"order/config"
	"order/infrastructure/log"
	"order/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	orderHandler := handler.NewOrderHandler()

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *orderHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
