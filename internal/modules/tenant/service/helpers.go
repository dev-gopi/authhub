package service

import (
	"fmt"
	"strings"
)

func optionalString(
	value string,
) *string {
	value = strings.TrimSpace(value)

	if value == "" {
		return nil
	}

	return &value
}

func buildIssuerURI(
	baseURL string,
	tenantAPILabel string,
	poolAPILabel string,
) string {
	baseURL = strings.TrimRight(
		baseURL,
		"/",
	)

	return fmt.Sprintf(
		"%s/%s/%s",
		baseURL,
		tenantAPILabel,
		poolAPILabel,
	)
}
