package main

import (
	"user/cmd/user/handler"
	"user/config"
	"user/infrastructure/log"
	"user/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	userHandler := handler.NewUserHandler()

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *userHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
