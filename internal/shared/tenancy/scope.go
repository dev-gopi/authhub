package tenancy

import "fmt"

func ScopeQuery(query string, tenantID string) string {
	fmt.Printf("Scoping query for tenant %s: %s\n", tenantID, query)
	return query
}
