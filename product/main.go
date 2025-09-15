package main

import (
	"product/cmd/product/handler"
	"product/cmd/product/repository"
	"product/cmd/product/resource"
	"product/cmd/product/service"
	"product/cmd/product/usecase"
	"product/config"
	"product/infrastructure/log"
	"product/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)

	log.SetupLogger()

	productRepository := repository.NewProductRepository(redis)
	productService := service.NewProductService(*productRepository)
	productUsecase := usecase.NewProductUsecase(*productService)
	productHandler := handler.NewProductHandler(*productUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *productHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
