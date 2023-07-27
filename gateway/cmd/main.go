package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jinford/rabbitmq-go-example/calculator/pkg/event"
	"github.com/jinford/rabbitmq-go-example/calculator/pkg/service"
	"github.com/jinford/rabbitmq-go-example/shared/message"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/sync/errgroup"
)

var (
	addResultMem = map[string]int{}
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	rabbitmqConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("amqp.Dial: %w", err)
	}

	producer := NewMessageProducer(rabbitmqConn)

	htppSrv := NewHttpServer(producer)

	subscriber := NewEventSubscriber(rabbitmqConn)

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return htppSrv.Run(ctx)
	})

	eg.Go(func() error {
		routingKey := message.EeventRoutingKey(service.Calculator, event.Calculated)
		return subscriber.Subscribe(ctx, routingKey, CalculatedEventHandler())
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
}
