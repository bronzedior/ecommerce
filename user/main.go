package main

import (
	"user/cmd/user/handler"
	"user/cmd/user/repository"
	"user/cmd/user/resource"
	"user/cmd/user/service"
	"user/cmd/user/usecase"
	"user/config"
	"user/infrastructure/log"
	"user/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	log.SetupLogger()

	userRepository := repository.NewUserRepository(db, redis)
	userService := service.NewUserService(*userRepository)
	userUsecase := usecase.NewUserUsecase(*userService, cfg.Secret.JWTSecret)
	userHandler := handler.NewUserHandler(*userUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *userHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
