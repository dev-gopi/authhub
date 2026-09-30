package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	auditmodel "github.com/dev-gopi/authhub/internal/modules/audit/model"
	outboxmodel "github.com/dev-gopi/authhub/internal/modules/outbox/model"
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (s *Service) createProvisioningEvents(
	ctx context.Context,
	tx *gorm.DB,
	tenantID uuid.UUID,
	actorID uuid.UUID,
	primaryAdminID uuid.UUID,
	primaryAdminEmail string,
	encryptedTemporaryPassword string,
	now time.Time,
) error {
	payloadBytes, err := json.Marshal(map[string]any{
		"tenant_id":                     tenantID.String(),
		"primary_admin_id":              primaryAdminID.String(),
		"primary_admin_email":           primaryAdminEmail,
		"temporary_password_ciphertext": encryptedTemporaryPassword,
		"credential_delivery":           "vault_transit_ciphertext",
	})
	if err != nil {
		return fmt.Errorf("marshal tenant-created outbox payload: %w", err)
	}

	event := &outboxmodel.Event{
		BaseModel: sharedmodel.NewBaseModelAt(now, &actorID),
		TenantID:  &tenantID, EventType: "tenant.created", AggregateType: "tenant", AggregateID: &tenantID,
		Payload: datatypes.JSON(payloadBytes), Status: "pending", Attempts: 0, AvailableAt: now,
	}
	if err := s.outbox.Create(ctx, tx, event); err != nil {
		return fmt.Errorf("create tenant-created outbox event: %w", err)
	}

	auditMetadata, err := json.Marshal(map[string]any{"primary_admin_id": primaryAdminID.String()})
	if err != nil {
		return fmt.Errorf("marshal tenant audit metadata: %w", err)
	}
	auditEvent := &auditmodel.Event{
		BaseModel: sharedmodel.NewBaseModelAt(now, &actorID),
		TenantID:  &tenantID, ActorPlatformUserID: &actorID, Action: "tenant.create", ResourceType: "tenant", ResourceID: &tenantID,
		Success: true, Metadata: datatypes.JSON(auditMetadata),
	}
	if err := s.audit.Create(ctx, tx, auditEvent); err != nil {
		return fmt.Errorf("create tenant audit event: %w", err)
	}
	return nil
}

func (s *Service) recordProvisioningFailureAudit(ctx context.Context, tenantID, actorID uuid.UUID, errValue error) {
	if s.audit == nil {
		return
	}
	now := time.Now().UTC()
	metadata, err := json.Marshal(map[string]any{"error": "tenant provisioning failed"})
	if err != nil {
		return
	}
	event := &auditmodel.Event{
		BaseModel: sharedmodel.NewBaseModelAt(now, &actorID),
		TenantID:  &tenantID, ActorPlatformUserID: &actorID, Action: "tenant.create", ResourceType: "tenant", ResourceID: &tenantID,
		Success: false, Metadata: datatypes.JSON(metadata),
	}
	if err := s.audit.Create(ctx, s.db, event); err != nil && s.logger != nil {
		s.logger.Warn("failed to persist tenant provisioning failure audit", zap.Error(err), zap.String("tenant_id", tenantID.String()))
	}
	_ = errValue
}
