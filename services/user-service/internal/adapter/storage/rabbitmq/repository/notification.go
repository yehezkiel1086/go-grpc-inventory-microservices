package repository

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/user-service/internal/adapter/storage/rabbitmq"
)

type NotificationRepository struct {
	mq *rabbitmq.Rabbitmq
	nq *amqp.Queue
}

func NewNotificationRepository(mq *rabbitmq.Rabbitmq) (*NotificationRepository, error) {
	nq, err := mq.DeclareQueue("notification_queue")
	if err != nil {
		return nil, err
	}

	return &NotificationRepository{
		mq: mq,
		nq: nq,
	}, nil
}

func (nr *NotificationRepository) SendEmailNotification(ctx context.Context, email string) error {
	return nr.mq.Publish(ctx, nr.nq, fmt.Sprintf("%s: user registered successfully", email))
}
