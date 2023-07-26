package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageProducer struct {
	conn *amqp.Connection
}

func NewMessageProducer(conn *amqp.Connection) *MessageProducer {
	return &MessageProducer{
		conn: conn,
	}
}

func (p *MessageProducer) Produce(ctx context.Context, routingKey string, msgBody []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("p.conn.Channel: %w", err)
	}
	defer ch.Close()

	if err := ch.PublishWithContext(ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBody,
		},
	); err != nil {
		return fmt.Errorf("ch.PublishWithContext: %w", err)
	}

	return nil
}
