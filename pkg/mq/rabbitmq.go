package mq

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
)

var conn *amqp091.Connection

func Init(addr, username, password string) {
	var err error
	url := fmt.Sprintf("amqp://%s:%s@%s/", username, password, addr)
	conn, err = amqp091.Dial(url)
	if err != nil {
		panic(err)
	}
}

func Publish(exchange, routeKey string, data []byte) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	msg := amqp091.Publishing{
		ContentType: "text/plain",
		Body:        data,
	}
	return ch.PublishWithContext(context.Background(), exchange, routeKey, true, false, msg)
}

func Consume(ctx context.Context, queue string, fn func(ctx context.Context, msgs <-chan amqp091.Delivery)) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	consume, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	go fn(ctx, consume)
	return nil
}
