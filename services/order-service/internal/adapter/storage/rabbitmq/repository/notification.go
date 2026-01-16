package repository

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/storage/rabbitmq"
)

type NotificationRepository struct {
	mq *rabbitmq.Rabbitmq
	q *amqp.Queue
}

func NewNotificationRepository(mq *rabbitmq.Rabbitmq) (*NotificationRepository, error) {
	q, err := mq.DeclareQueue("notification_queue")
	if err != nil {
		return nil, err
	}

	return &NotificationRepository{
		mq: mq,
		q: q,
	}, nil
}

func (nr *NotificationRepository) SendNotification(ctx context.Context, msg string) error {
	return nr.mq.Publish(ctx, nr.q, msg)
}
