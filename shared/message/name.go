package message

import "fmt"

type ServiceName string

type EventName string

type ServiceCommandName string

func EeventRoutingKey(service ServiceName, name EventName) string {
	return fmt.Sprintf("%s.event.%s", service, name)
}

func ServiceCommandRoutingKey(service ServiceName, name ServiceCommandName) string {
	return fmt.Sprintf("%s.command.%s", service, name)
}
