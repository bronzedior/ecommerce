package routes

import (
	"product/cmd/product/handler"
	"product/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, productHandler handler.ProductHandler) {
	router.Use(middleware.RequestLogger())
	router.GET("/ping", productHandler.Ping)
}
