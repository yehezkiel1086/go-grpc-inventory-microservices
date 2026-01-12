package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/notification-service/internal/adapter/config"
)

type Rabbitmq struct {
	conn *amqp.Connection
	ch *amqp.Channel
}

func New(conf *config.Rabbitmq) (*Rabbitmq, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", conf.User, conf.Password, conf.Host, conf.Port)

	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Rabbitmq{
		conn: conn,
		ch: ch,
	}, nil
}

func (r *Rabbitmq) DeclareQueue() (*amqp.Queue, error) {
	q, err := r.ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (r *Rabbitmq) Consume(q *amqp.Queue) (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (r *Rabbitmq) CloseChan() error {
	return r.ch.Close()
}

func (r *Rabbitmq) CloseConn() error {
	return r.conn.Close()
}