package routes

import (
	"payment/cmd/payment/handler"
	"payment/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, paymentHandler handler.PaymentHandler) {
	router.Use(middleware.RequestLogger())
	router.POST("/v1/payment/webhook", paymentHandler.HandleXenditWebhook)
	router.GET("/v1/invoice/:order_id/pdf", paymentHandler.HandleDownloadPDFInvoice)
}
