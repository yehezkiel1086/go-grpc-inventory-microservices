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
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/rabbitmq"
	rabbitmqRepo "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/rabbitmq/repository"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/core/service"
	"google.golang.org/grpc"
)

func failOnError(err error, msg string) {
  if err != nil {
    log.Panicf("%s: %s", msg, err)
  }
}

func main() {
	// load .env configs
	conf, err := config.New()
	failOnError(err, "failed to load .env configs")
	fmt.Println(".env configs loaded successfully")

	// init context
	ctx := context.Background()

	// init db connection
	db, err := postgres.New(ctx, conf.DB)
	failOnError(err, "failed to init db connection")
	fmt.Println("DB connection established successfully")

	// migrate dbs
	err = db.Migrate(&domain.User{}, &domain.Product{}, &domain.Order{})
	failOnError(err, "failed to migrate dbs")
	fmt.Println("DB migrated successfully")

	// init grpc connection
	clientUri := fmt.Sprintf("%s:%s", conf.GRPC.Host, conf.GRPC.Port)
	conn, err := grpc.Dial(clientUri, grpc.WithInsecure())
	failOnError(err, "failed to init grpc connection")
	fmt.Println("GRPC connection established successfully")

	// init grpc clients
	inventoryClient := inventory.NewInventoryServiceClient(conn)

	// init rabbitmq
	mq, err := rabbitmq.New(conf.Rabbitmq)
	failOnError(err, "failed to init rabbitmq")
	fmt.Println("rabbitmq initialized successfully")

	defer mq.CloseConn()
	defer mq.CloseChan()

	// declare queue
	q, err := mq.DeclareQueue("notification_queue")
	failOnError(err, "failed to declare queue")
	fmt.Println("queue declared successfully")

	// dependency injections
	notifRepo, err := rabbitmqRepo.NewNotificationRepository(mq)
	failOnError(err, "failed to init notification repository")

	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo, inventoryClient, notifRepo)
	productHandler := handler.NewProductHandler(productSvc)

	orderRepo := repository.NewOrderRepository(db)
	orderSvc := service.NewOrderService(orderRepo, inventoryClient, mq, q)
	orderHandler := handler.NewOrderHandler(orderSvc)

	// init router
	r := handler.NewRouter(
		conf.HTTP,
		productHandler,
		orderHandler,
	)

	// start server
	if err := r.Serve(conf.HTTP); err != nil {
		log.Fatal(err)
	}
}
