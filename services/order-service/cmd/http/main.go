package main

import (
	"context"
	"fmt"
	"log"

	inventory "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/common/genproto/inventory/protobuf"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/handler"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/service"
	"google.golang.org/grpc"
)

func main() {
	// load .env configs
	conf, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(".env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init db connection
	db, err := postgres.New(ctx, conf.DB)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB connection established successfully")

	// migrate dbs
	if err := db.Migrate(&domain.User{}, &domain.Product{}, &domain.Order{}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("DB migrated successfully")

	// init grpc connection
	clientUri := fmt.Sprintf("%s:%s", conf.GRPC.Host, conf.GRPC.Port)
	conn, err := grpc.Dial(clientUri, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("GRPC connection established successfully")

	// init grpc clients
	inventoryClient := inventory.NewInventoryServiceClient(conn)

	// dependency injections
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	authSvc := service.NewAuthService(conf.JWT, userRepo)
	authHandler := handler.NewAuthHandler(conf.JWT, authSvc)

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo, inventoryClient)
	productHandler := handler.NewProductHandler(productSvc)

	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, inventoryClient)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// init router
	r := handler.NewRouter(
		conf.HTTP,
		userHandler,
		authHandler,
		productHandler,
		orderHandler,
	)

	// start server
	if err := r.Serve(conf.HTTP); err != nil {
		log.Fatal(err)
	}
}
