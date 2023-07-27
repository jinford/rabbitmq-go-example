package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jinford/rabbitmq-go-example/calculator/pkg/command"
	"github.com/jinford/rabbitmq-go-example/calculator/pkg/event"
	"github.com/jinford/rabbitmq-go-example/calculator/pkg/service"
	"github.com/jinford/rabbitmq-go-example/shared/message"
)

func AddCommandHandler(p *EventPublisher) func(ctx context.Context, msgBpdy []byte) error {
	return func(ctx context.Context, msgBpdy []byte) error {
		reqBody := new(command.RequestBody)
		if err := json.Unmarshal(msgBpdy, reqBody); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}

		log.Println("[CONSUME MESSAGE] caluculator.add:", reqBody.ID)

		result := CalcAdd(reqBody.A, reqBody.B)
		log.Printf("[CALCULATE] A + B = %d\n", result)

		eventBody := &event.CalculatedEventBody{
			ID:     reqBody.ID,
			Result: result,
		}

		msgBody, err := json.Marshal(eventBody)
		if err != nil {
			return fmt.Errorf("json.Marshal: %w", err)
		}

		log.Println("[PUBLISH EVNET] calculated:", reqBody.ID)

		routingKey := message.EeventRoutingKey(service.Calculator, event.Calculated)
		if err := p.Publish(ctx, routingKey, msgBody); err != nil {
			return fmt.Errorf("p.Publish: %w", err)
		}

		return nil
	}

}

func CalcAdd(a int, b int) int {
	return a + b
}
