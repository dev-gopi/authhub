package entity

type PlatformUserStatus string

const (
	PlatformUserStatusActive   PlatformUserStatus = "active"
	PlatformUserStatusDisabled PlatformUserStatus = "disabled"
	PlatformUserStatusBlocked  PlatformUserStatus = "blocked"
)

func (s PlatformUserStatus) String() string {
	return string(s)
}

func (s PlatformUserStatus) IsValid() bool {
	switch s {
	case PlatformUserStatusActive,
		PlatformUserStatusDisabled,
		PlatformUserStatusBlocked:
		return true
	default:
		return false
	}
}
