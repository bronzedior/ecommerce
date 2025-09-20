package routes

import (
	"payment/cmd/payment/handler"
	"payment/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, paymentHandler handler.PaymentHandler) {
	router.Use(middleware.RequestLogger())
	router.GET("/ping", paymentHandler.Ping)
}
