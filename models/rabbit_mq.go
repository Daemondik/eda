package models

import (
	"github.com/streadway/amqp"
	"sync"
)

var (
	conn     *amqp.Connection
	connErr  error
	connOnce sync.Once
)

const ChatQueue = "chat"

func initRabbitMQ() {
	conn, connErr = amqp.Dial("amqp://admin:admin@rabbitmq/")
	if connErr != nil {
		panic(connErr)
	}
}

func RabbitMQQueue(name string) (*amqp.Channel, error) {
	connOnce.Do(initRabbitMQ)
	if connErr != nil {
		return nil, connErr
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
