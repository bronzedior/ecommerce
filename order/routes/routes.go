package routes

import (
	"order/cmd/order/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, orderHandler handler.OrderHandler) {
	router.GET("/ping", orderHandler.Ping)
}
