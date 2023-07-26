package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageConsumer struct {
	conn *amqp.Connection
}

func NewMessageConsumer(conn *amqp.Connection) *MessageConsumer {
	return &MessageConsumer{
		conn: conn,
	}
}

func (c *MessageConsumer) Consume(ctx context.Context, routingKey string, handler func(ctx context.Context, msgBody []byte) error) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("c.conn.Channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		routingKey, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("ch.QueueDeclare: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("ch.Consume: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			if err := handler(ctx, msg.Body); err != nil {
				return fmt.Errorf("handler: %w", err)
			}
		}
	}
}
