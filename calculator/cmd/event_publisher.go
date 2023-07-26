package main

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPublisher struct {
	conn *amqp.Connection
}

func NewEventPublisher(conn *amqp.Connection) *EventPublisher {
	return &EventPublisher{
		conn: conn,
	}
}

func (p *EventPublisher) Publish(ctx context.Context, eventName string, msgBody []byte) error {
	ch, err := p.conn.Channel()
	if err != nil {
		return fmt.Errorf("p.conn.Channel: %w", err)
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
		return fmt.Errorf("ch.ExchangeDeclare: %w", err)
	}

	if err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msgBody,
		},
	); err != nil {
		return fmt.Errorf("ch.PublishWithContext: %w", err)
	}

	return nil
}
