package vault

import (
	"context"
	"fmt"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

type TransitKeySpec struct {
	Name string
	Type string
}

// EnsureTransitKey creates a Transit key when it does not already exist.
//
// The method is intentionally idempotent:
//
// existing key
//
//	-> success
//
// missing key
//
//	-> create
//	-> success
func (c *Client) EnsureTransitKey(
	ctx context.Context,
	mount string,
	spec TransitKeySpec,
) error {
	mount = strings.Trim(
		strings.TrimSpace(mount),
		"/",
	)

	if mount == "" {
		return fmt.Errorf(
			"vault transit mount is required",
		)
	}

	if strings.TrimSpace(spec.Name) == "" {
		return fmt.Errorf(
			"vault transit key name is required",
		)
	}

	if strings.TrimSpace(spec.Type) == "" {
		return fmt.Errorf(
			"vault transit key type is required",
		)
	}

	path := fmt.Sprintf(
		"%s/keys/%s",
		mount,
		spec.Name,
	)

	existing, err := c.Client.Logical().
		ReadWithContext(
			ctx,
			path,
		)

	if err != nil {
		return fmt.Errorf(
			"read vault transit key %q: %w",
			spec.Name,
			err,
		)
	}

	if existing != nil {
		return nil
	}

	payload := map[string]interface{}{
		"type": spec.Type,

		"exportable": false,

		"allow_plaintext_backup": false,
	}

	_, err = c.Client.Logical().
		WriteWithContext(
			ctx,
			path,
			payload,
		)

	if err != nil {
		return fmt.Errorf(
			"create vault transit key %q: %w",
			spec.Name,
			err,
		)
	}

	return nil
}

// ReadTransitKey returns metadata for an existing Transit key.
//
// This returns Vault metadata only.
// It does not export private key material.
func (c *Client) ReadTransitKey(
	ctx context.Context,
	mount string,
	keyName string,
) (*vaultapi.Secret, error) {
	mount = strings.Trim(
		strings.TrimSpace(mount),
		"/",
	)

	keyName = strings.TrimSpace(
		keyName,
	)

	if mount == "" {
		return nil, fmt.Errorf(
			"vault transit mount is required",
		)
	}

	if keyName == "" {
		return nil, fmt.Errorf(
			"vault transit key name is required",
		)
	}

	path := fmt.Sprintf(
		"%s/keys/%s",
		mount,
		keyName,
	)

	secret, err := c.Client.Logical().
		ReadWithContext(
			ctx,
			path,
		)

	if err != nil {
		return nil, fmt.Errorf(
			"read vault transit key %q: %w",
			keyName,
			err,
		)
	}

	return secret, nil
}
