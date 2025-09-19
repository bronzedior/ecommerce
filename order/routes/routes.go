package routes

import (
	"order/cmd/order/handler"
	"order/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, orderHandler handler.OrderHandler, jwtSecret string) {
	router.Use(middleware.RequestLogger())
	router.GET("/ping", orderHandler.Ping)
	authMiddleware := middleware.AuthMiddleware(jwtSecret)
	router.Use(authMiddleware)
	router.GET("/v1/order_history", orderHandler.GetOrderHistory)
	router.POST("/v1/checkout", orderHandler.Checkout)
}
