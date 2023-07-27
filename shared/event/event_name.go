package event

import "fmt"

type EventName string

func RoutingKey(eventName EventName) string {
	return fmt.Sprintf("calclator.event.%s", eventName)
}
