package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/handler"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/storage/postgres"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/storage/postgres/repository"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/storage/rabbitmq"
	rabbitmqRepository "github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/storage/rabbitmq/repository"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/core/domain"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/core/service"
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
	failOnError(err, "unable to connect to postgres db")
	fmt.Println("DB connection established successfully")

	// migrate dbs
	err = db.Migrate(&domain.User{})
	failOnError(err, "failed to migrate dbs")
	fmt.Println("DB migrated successfully")

	// init rabbitmq
	mq, err := rabbitmq.New(conf.Rabbitmq)
	failOnError(err, "failed to establish connection with rabbitmq")
	fmt.Println("RabbitMQ connection established successfully")

	defer mq.CloseConn()
	defer mq.CloseChan()

	// dependency injections
	notifRepo, err := rabbitmqRepository.NewNotificationRepository(mq)
	failOnError(err, "failed to declare notification queue")

	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo, notifRepo)
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
