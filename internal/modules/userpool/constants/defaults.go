package constants

const (
	DefaultPoolAPILabel = "default"

	DefaultPoolDisplayName = "Default User Pool"

	DefaultPoolStatus = "active"

	DefaultLocale = "en"
)

const (
	DefaultMFAPolicy = "optional"

	DefaultPostRecoverySessionAction = "revoke_all"
)

/*
These numeric values are implementation defaults.

Your Task 5 requirements do not prescribe exact numeric
values for session durations/password lengths.

Keep them centralized here for now instead of scattering
magic numbers through provisioning code.
*/
const (
	DefaultSessionIdleSeconds = 30 * 60

	DefaultSessionAbsoluteSeconds = 12 * 60 * 60

	DefaultRecentAuthSeconds = 5 * 60

	DefaultPasswordMinLength = 14

	DefaultPasswordMaxLength = 128

	DefaultPasswordHistoryCount = 5

	DefaultPasswordResetTokenTTLSeconds = 15 * 60
)
