package vault

import (
	"context"
	"encoding/base64"
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

func (c *Client) EncryptTransit(
	ctx context.Context,
	mount string,
	keyName string,
	plaintext []byte,
) (string, error) {
	mount = strings.Trim(strings.TrimSpace(mount), "/")
	keyName = strings.TrimSpace(keyName)
	if mount == "" || keyName == "" {
		return "", fmt.Errorf("vault transit mount and key name are required")
	}
	if len(plaintext) == 0 {
		return "", fmt.Errorf("vault transit plaintext is required")
	}

	path := fmt.Sprintf("%s/encrypt/%s", mount, keyName)
	payload := map[string]interface{}{
		"plaintext": base64.StdEncoding.EncodeToString(plaintext),
	}
	secret, err := c.Client.Logical().WriteWithContext(ctx, path, payload)
	if err != nil {
		return "", fmt.Errorf("encrypt with vault transit key %q: %w", keyName, err)
	}
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("vault transit key %q returned empty encryption response", keyName)
	}
	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok || ciphertext == "" {
		return "", fmt.Errorf("vault transit key %q returned no ciphertext", keyName)
	}
	return ciphertext, nil
}
