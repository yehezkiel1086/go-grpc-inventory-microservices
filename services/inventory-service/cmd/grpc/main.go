package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/adapter/handler"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/inventory-service/internal/core/service"
)

func main() {
	// load .env configs
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(".env configs loaded successfully")

	// init postgres db conn
	ctx := context.Background()
	
	db, err := postgres.New(ctx, conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("database connected successfully")

	// migrate dbs
	if err := db.Migrate(&domain.Inventory{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("database migrated successfully")

	// dependency injections
	invRepo := repository.NewInventoryRepository(db)
	invSvc := service.NewInventoryService(invRepo)
	invHandler := handler.NewInventoryHandler(invSvc)

	// init server
	s := handler.NewServer(
		invHandler,
	)

	// run server
	if err := s.Run(conf.GRPC); err != nil {
		log.Fatal(err)
	}
}
