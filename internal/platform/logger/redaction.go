package logger

const redactedValue = "[REDACTED]"

// RedactSensitiveData must be used when a value is known to contain a secret.
// It intentionally never returns the input value.
func RedactSensitiveData(_ string) string {
	return redactedValue
}
