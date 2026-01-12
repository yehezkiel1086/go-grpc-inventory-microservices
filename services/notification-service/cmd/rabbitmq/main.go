package main

import (
	"fmt"
	"log"

	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/notification-service/internal/adapter/config"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/notification-service/internal/adapter/storage/rabbitmq"
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

	// init rabbitmq
	r, err := rabbitmq.New(conf.Rabbitmq)
	failOnError(err, "failed to init rabbitmq")
	fmt.Println("rabbitmq initialized successfully")

	defer r.CloseConn()
	defer r.CloseChan()

	// declare queue
	q, err := r.DeclareQueue()
	failOnError(err, "failed to declare queue")
	fmt.Println("queue declared successfully")

	// consume messages
	msgs, err := r.Consume(q)
	failOnError(err, "Failed to register a consumer")
	fmt.Println("consumer registered successfully")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
