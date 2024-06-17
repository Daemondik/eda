package models

import (
	"github.com/streadway/amqp"
)

const ChatQueue = "chat"

func NewRabbitMQQueue(name string) (*amqp.Channel, error) {
	conn, err := amqp.Dial("amqp://admin:admin@rabbitmq/")
	if err != nil {
		return nil, err
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
