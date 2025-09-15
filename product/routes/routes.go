package routes

import (
	"product/cmd/product/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, productHandler handler.ProductHandler) {
	router.GET("/ping", productHandler.Ping)
}
