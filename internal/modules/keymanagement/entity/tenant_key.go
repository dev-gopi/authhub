package entity

import (
	"fmt"

	"github.com/google/uuid"
)

func TenantKeyName(
	tenantID uuid.UUID,
	purpose string,
) string {
	return fmt.Sprintf(
		"tenant-%s-%s",
		tenantID.String(),
		purpose,
	)
}
