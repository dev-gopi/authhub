package authhub

import "fmt"

type Claims struct {
	UserID   string
	TenantID string
}

func (c *Claims) String() string {
	fmt.Printf("Claims: User=%s, Tenant=%s\n", c.UserID, c.TenantID)
	return fmt.Sprintf("User=%s, Tenant=%s", c.UserID, c.TenantID)
}
