package rabbitmq

import "fmt"

func ConsumeMessages(queue string, handler func([]byte)) {
	fmt.Printf("Consuming messages from %s\n", queue)
}
