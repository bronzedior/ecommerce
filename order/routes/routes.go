package routes

import (
	"order/cmd/order/handler"
	"order/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, orderHandler handler.OrderHandler) {
	router.Use(middleware.RequestLogger())
	router.GET("/ping", orderHandler.Ping)
}
