package rabbitmq

import (
	"context"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	Conn *amqp.Connection
}

func NewConnection(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect rabbitmq: %w", err)
	}

	return &Connection{
		Conn: conn,
	}, nil
}

func (c *Connection) Health(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if c.Conn == nil {
		return errors.New("rabbitmq connection is nil")
	}

	if c.Conn.IsClosed() {
		return errors.New("rabbitmq connection is closed")
	}

	channel, err := c.Conn.Channel()
	if err != nil {
		return fmt.Errorf("open rabbitmq health channel: %w", err)
	}

	return channel.Close()
}

func (c *Connection) Close() error {
	if c.Conn == nil || c.Conn.IsClosed() {
		return nil
	}

	return c.Conn.Close()
}
