package entity

import "errors"

var (
	ErrUnauthorized         = errors.New("root authorization required")
	ErrTenantAlreadyExists  = errors.New("tenant api label already exists")
	ErrPrimaryAdminConflict = errors.New("primary admin username or email already exists")
	ErrProvisioningFailed   = errors.New("tenant provisioning failed")
)
