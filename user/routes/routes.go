package routes

import (
	"user/cmd/user/handler"
	"user/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler handler.UserHandler) {
	// Public API
	router.Use(middleware.RequestLogger())
	router.GET("/ping", userHandler.Ping)
	router.POST("/v1/register", userHandler.Register)

	// Private API
}
