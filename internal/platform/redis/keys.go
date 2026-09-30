package redis

const (
	KeyPrefix = "authhub:"
)

func GetTenantKey(tenantID string) string {
	return KeyPrefix + "tenant:" + tenantID
}
