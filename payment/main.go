package main

import (
	"context"
	"payment/cmd/payment/handler"
	"payment/cmd/payment/repository"
	"payment/cmd/payment/resource"
	"payment/cmd/payment/service"
	"payment/cmd/payment/usecase"
	"payment/config"
	"payment/infrastructure/constant"
	"payment/infrastructure/log"
	"payment/kafka"
	"payment/models"
	"payment/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := resource.InitDB(&cfg)
	kafkaWriter := kafka.NewWriter(cfg.Kafka.Broker, cfg.Kafka.KafkaTopics[constant.KafkaTopicPaymentSuccess])

	log.SetupLogger()

	databaseRepository := repository.NewPaymentDatabase(db)
	publisherRepository := repository.NewKafkaPublisher(kafkaWriter)
	paymentService := service.NewPaymentService(databaseRepository, publisherRepository)
	paymentUsecase := usecase.NewPaymentUsecase(paymentService)
	paymentHandler := handler.NewPaymentHandler(paymentUsecase, cfg.Xendit.XenditWebhookToken)

	xenditRepository := repository.NewXenditClient(cfg.Xendit.XenditAPIKey)
	xenditService := service.NewXenditService(databaseRepository, xenditRepository)
	xenditUsecase := usecase.NewXenditUsecase(xenditService)

	kafka.StartOrderConsumer(cfg.Kafka.Broker, cfg.Kafka.KafkaTopics[constant.KafkaTopicOrderCreated],
		func(event models.OrderCreatedEvent) {
			if err := xenditUsecase.CreateInvoice(context.Background(), event); err != nil {
				log.Logger.Println("failed handling order created event: ", err.Error())
			}
		})

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, paymentHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
