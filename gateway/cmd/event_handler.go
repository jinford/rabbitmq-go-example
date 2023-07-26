package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jinford/rabbitmq-go-example/calculator/pkg/event"
)

func CalculatedEventHandler() func(context.Context, []byte) error {
	return func(ctx context.Context, msgBody []byte) error {
		eventBody := new(event.CalculatedEventBody)
		if err := json.Unmarshal(msgBody, eventBody); err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}

		log.Println("[SUBSCRIBE EVNET] calculated", eventBody.ID)

		addResultMem[eventBody.ID] = eventBody.Result

		return nil
	}
}
