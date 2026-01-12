package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/yehezkiel1086/go-grpc-inventory-microservices/services/order-service/internal/adapter/config"
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

func (r *Rabbitmq) DeclareQueue(name string) (*amqp.Queue, error) {
	q, err := r.ch.QueueDeclare(
		name, // name
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

func (r *Rabbitmq) Publish(ctx context.Context, q *amqp.Queue, body string) error {
	return r.ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
	})
}

func (r *Rabbitmq) CloseChan() error {
	return r.ch.Close()
}

func (r *Rabbitmq) CloseConn() error {
	return r.conn.Close()
}