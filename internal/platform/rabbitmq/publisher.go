package rabbitmq

import "fmt"

func PublishMessage(exchange string, routingKey string, body []byte) {
	fmt.Printf("Publishing message to %s/%s\n", exchange, routingKey)
}
