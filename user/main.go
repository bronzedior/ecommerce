package main

import (
	"net"
	"user/cmd/user/handler"
	"user/cmd/user/repository"
	"user/cmd/user/resource"
	"user/cmd/user/service"
	"user/cmd/user/usecase"
	"user/config"
	grpcUser "user/grpc"
	"user/infrastructure/log"
	"user/proto/userpb"
	"user/routes"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// cfg := config.LoadConfig()
	// redis := resource.InitRedis(&cfg)
	// db := resource.InitDB(&cfg)

	// log.SetupLogger()

	// userRepository := repository.NewUserRepository(db, redis)
	// userService := service.NewUserService(*userRepository)
	// userUsecase := usecase.NewUserUsecase(*userService, cfg.Secret.JWTSecret)
	// userHandler := handler.NewUserHandler(*userUsecase)

	// go func() {
	// 	port := cfg.App.Port
	// 	router := gin.Default()
	// 	routes.SetupRoutes(router, *userHandler, cfg.Secret.JWTSecret)
	// 	router.Run(":" + port)

	// 	log.Logger.Printf("HTTP Server running on port: %s", port)
	// }()

	// grpcServer := grpc.NewServer()
	// userpb.RegisterUserServiceServer(grpcServer, &grpcUser.GRPCServer{UserUsecase: *userUsecase})

	// lis, _ := net.Listen("tcp", ":50051")
	// grpcServer.Serve(lis)

	// log.Logger.Printf("GRPC Server running on port: %s", ":50051")
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	log.SetupLogger()

	userRepository := repository.NewUserRepository(db, redis)
	userService := service.NewUserService(*userRepository)
	userUsecase := usecase.NewUserUsecase(*userService, cfg.Secret.JWTSecret)
	userHandler := handler.NewUserHandler(*userUsecase)

	// Start HTTP server in a goroutine
	go func() {
		port := cfg.App.Port
		router := gin.Default()
		routes.SetupRoutes(router, *userHandler, cfg.Secret.JWTSecret)

		log.Logger.Printf("HTTP Server running on port: %s", port)
		if err := router.Run(":" + port); err != nil {
			log.Logger.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	// Start gRPC server
	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, &grpcUser.GRPCServer{UserUsecase: *userUsecase})
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Logger.Fatalf("failed to listen: %v", err)
	}

	log.Logger.Printf("gRPC Server running on port: %s", ":50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Logger.Fatalf("failed to serve gRPC: %v", err)
	}
}
