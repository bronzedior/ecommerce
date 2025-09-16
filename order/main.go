package main

import (
	"order/cmd/order/handler"
	"order/cmd/order/repository"
	"order/cmd/order/resource"
	"order/cmd/order/service"
	"order/cmd/order/usecase"
	"order/config"
	"order/infrastructure/log"
	"order/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	log.SetupLogger()

	orderRepository := repository.NewOrderRepository(db, redis)
	orderService := service.NewOrderService(*orderRepository)
	orderUsecase := usecase.NewOrderUsecase(*orderService)
	orderHandler := handler.NewOrderHandler(*orderUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *orderHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
