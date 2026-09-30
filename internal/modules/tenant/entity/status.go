package entity

type Status string

const (
	StatusProvisioning    Status = "provisioning"
	StatusActive          Status = "active"
	StatusSuspended       Status = "suspended"
	StatusDeleteRequested Status = "delete_requested"
	StatusDeleted         Status = "deleted"
)

func (s Status) String() string {
	return string(s)
}

func (s Status) IsValid() bool {
	switch s {
	case StatusProvisioning,
		StatusActive,
		StatusSuspended,
		StatusDeleteRequested,
		StatusDeleted:
		return true
	default:
		return false
	}
}
