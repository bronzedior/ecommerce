package main

import (
	"order/cmd/order/handler"
	"order/cmd/order/repository"
	"order/cmd/order/resource"
	"order/cmd/order/service"
	"order/cmd/order/usecase"
	"order/config"
	"order/infrastructure/log"
	"order/kafka"
	"order/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	log.SetupLogger()

	kafkaProducer := kafka.NewKafkaProducer([]string{"localhost:9093"}, "order.created")
	defer kafkaProducer.Close()

	orderRepository := repository.NewOrderRepository(db, redis, cfg.Product.Host)
	orderService := service.NewOrderService(*orderRepository)
	orderUsecase := usecase.NewOrderUsecase(*orderService, *kafkaProducer)
	orderHandler := handler.NewOrderHandler(*orderUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *orderHandler, cfg.Secret.JWTSecret)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
