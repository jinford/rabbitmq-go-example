package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventSubscriber struct {
	conn *amqp.Connection
}

func NewEventSubscriber(
	conn *amqp.Connection,
) *EventSubscriber {
	return &EventSubscriber{
		conn: conn,
	}
}

func (sub *EventSubscriber) Subscribe(ctx context.Context, eventName string, handler func(ctx context.Context, msgBody []byte) error) error {
	ch, err := sub.conn.Channel()
	if err != nil {
		return fmt.Errorf("conn.Channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		eventName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
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
