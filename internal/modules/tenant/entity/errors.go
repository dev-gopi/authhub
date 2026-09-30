package entity

import "fmt"

type TenantError struct {
	Message string
}

func (e *TenantError) Error() string {
	fmt.Printf("Tenant error: %s\n", e.Message)
	return e.Message
}
