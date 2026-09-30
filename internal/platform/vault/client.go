package vault

import (
	"context"
	"errors"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

type Client struct {
	Client *vaultapi.Client
}

func NewClient(
	address string,
	token string,
) (*Client, error) {

	config := vaultapi.DefaultConfig()

	config.Address = address

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	}

	return &Client{
		Client: client,
	}, nil
}

func (c *Client) Health(ctx context.Context) error {
	health, err := c.Client.Sys().HealthWithContext(ctx)
	if err != nil {
		return fmt.Errorf("vault health check: %w", err)
	}

	if !health.Initialized {
		return errors.New("vault is not initialized")
	}

	if health.Sealed {
		return errors.New("vault is sealed")
	}

	return nil
}
