package entity

import "github.com/google/uuid"

type RootAuthContext struct {
	PlatformUserID uuid.UUID

	SessionID uuid.UUID

	CredentialVersion int64

	IsRootAdmin bool
}
