package main

import (
	"payment/cmd/payment/handler"
	"payment/cmd/payment/repository"
	"payment/cmd/payment/resource"
	"payment/cmd/payment/service"
	"payment/cmd/payment/usecase"
	"payment/config"
	"payment/infrastructure/log"
	"payment/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)

	log.SetupLogger()

	paymentRepository := repository.NewPaymentRepository(redis)
	paymentService := service.NewPaymentService(*paymentRepository)
	paymentUsecase := usecase.NewPaymentUsecase(*paymentService)
	paymentHandler := handler.NewPaymentHandler(*paymentUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *paymentHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
