package mq

import (
	"context"
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"sync"
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

func Consume(wg *sync.WaitGroup, queue string, fn func(wg *sync.WaitGroup, msgs <-chan amqp091.Delivery)) {
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	consume, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	go fn(wg, consume)
}
