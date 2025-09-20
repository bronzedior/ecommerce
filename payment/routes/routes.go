package routes

import (
	"payment/cmd/payment/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, paymentHandler handler.PaymentHandler) {
	router.GET("/ping", paymentHandler.Ping)
}
