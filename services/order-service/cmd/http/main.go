package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/handler"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/service"
)

func main() {
	// load .env configs
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ .env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init db connection
	db, err := postgres.New(ctx, conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ DB connection established successfully")

	// migrate dbs
	if err := db.Migrate(&domain.User{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("✅ DB migrated successfully")

	// dependency injections
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	// init router
	r := handler.NewRouter(
		conf.HTTP,
		userHandler,
		authHandler,
	)

	// start server
	if err := r.Serve(conf.HTTP); err != nil {
		log.Fatal(err)
	}
}
