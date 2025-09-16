package main

import (
	"order/config"
	"order/infrastructure/log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	log.SetupLogger()

	port := cfg.App.Port
	router := gin.Default()
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
