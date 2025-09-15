package routes

import (
	"user/cmd/user/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userHandler handler.UserHandler) {
	// Public API
	router.GET("/ping", userHandler.Ping)
	router.POST("/v1/register", userHandler.Register)

	// Private API
}
