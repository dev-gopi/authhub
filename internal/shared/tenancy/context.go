package tenancy

import "fmt"

type Context struct {
	TenantID string
}

func NewContext(tenantID string) *Context {
	fmt.Printf("Creating tenancy context for tenant: %s\n", tenantID)
	return &Context{TenantID: tenantID}
}
