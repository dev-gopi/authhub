package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	permissionmodel "github.com/dev-gopi/authhub/internal/modules/permission/model"
	roleconstants "github.com/dev-gopi/authhub/internal/modules/role/constants"
	rolemodel "github.com/dev-gopi/authhub/internal/modules/role/model"
	sharedmodel "github.com/dev-gopi/authhub/internal/shared/model"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (s *Service) provisionDefaultRoles(
	ctx context.Context,
	tx *gorm.DB,
	tenantID uuid.UUID,
	actorID uuid.UUID,
	now time.Time,
) (*roleProvisionResult, error) {
	definitions, err := s.loadDefaultRoleDefinitions(ctx, tx)
	if err != nil {
		return nil, err
	}

	allLabels := collectDefaultPermissionLabels(definitions)
	permissions, err := s.permissions.FindByAPILabels(ctx, tx, allLabels)
	if err != nil {
		return nil, fmt.Errorf("load global permissions: %w", err)
	}

	permissionMap := make(map[string]permissionmodel.Permission, len(permissions))
	for _, permission := range permissions {
		permissionMap[permission.APILabel] = permission
	}

	// Fail closed. A tenant must never receive a partially configured default role.
	for _, label := range allLabels {
		if _, ok := permissionMap[label]; !ok {
			return nil, fmt.Errorf("required permission %q is not seeded", label)
		}
	}

	result := &roleProvisionResult{}

	for _, definition := range definitions {
		description := definition.Description
		metadata, err := json.Marshal(definition.Metadata)
		if err != nil {
			return nil, fmt.Errorf("marshal metadata for default role %q: %w", definition.APILabel, err)
		}

		role := &rolemodel.Role{
			BaseModel:           sharedmodel.NewBaseModelAt(now, &actorID),
			TenantID:            &tenantID,
			APILabel:            definition.APILabel,
			DisplayName:         definition.DisplayName,
			Description:         &description,
			Metadata:            datatypes.JSON(metadata),
			IsDefault:           true,
			ProtectedFromDelete: definition.ProtectedFromDelete,
		}

		storedRole, err := s.roles.EnsureRole(ctx, tx, role)
		if err != nil {
			return nil, fmt.Errorf("ensure default role %q: %w", definition.APILabel, err)
		}

		if definition.APILabel == "primary_admin" {
			result.PrimaryAdminRoleID = storedRole.ID
		}

		for _, permissionLabel := range definition.PermissionLabels {
			permission := permissionMap[permissionLabel]
			rp := &rolemodel.RolePermission{
				BaseModel:    sharedmodel.NewBaseModelAt(now, &actorID),
				RoleID:       storedRole.ID,
				PermissionID: permission.ID,
			}

			if err := s.roles.EnsureRolePermission(ctx, tx, rp); err != nil {
				return nil, fmt.Errorf("assign permission %q to role %q: %w", permissionLabel, definition.APILabel, err)
			}
		}
	}

	if result.PrimaryAdminRoleID == uuid.Nil {
		return nil, fmt.Errorf("primary_admin role was not provisioned")
	}

	return result, nil
}

func (s *Service) loadDefaultRoleDefinitions(
	ctx context.Context,
	tx *gorm.DB,
) ([]roleconstants.DefaultRole, error) {
	version := s.defaultRoleTemplateVersion
	if version == "" {
		version = roleconstants.DefaultRoleTemplateVersion
	}

	snapshots, err := s.roles.ListDefaultRoleSnapshots(ctx, tx, version)
	if err != nil {
		return nil, fmt.Errorf("load default role snapshots for version %q: %w", version, err)
	}
	if len(snapshots) == 0 {
		return nil, fmt.Errorf("no default role snapshots found for version %q", version)
	}

	definitions := make([]roleconstants.DefaultRole, 0, len(snapshots))
	seen := make(map[string]struct{}, len(snapshots))

	for _, snapshot := range snapshots {
		if err := verifyDefaultRoleSnapshot(snapshot); err != nil {
			return nil, err
		}

		var definition roleconstants.DefaultRole
		if err := json.Unmarshal(snapshot.Definition, &definition); err != nil {
			return nil, fmt.Errorf("decode default role snapshot %q: %w", snapshot.RoleAPILabel, err)
		}
		if definition.APILabel == "" || definition.APILabel != snapshot.RoleAPILabel {
			return nil, fmt.Errorf("invalid default role snapshot %q: api_label mismatch", snapshot.RoleAPILabel)
		}
		if _, exists := seen[definition.APILabel]; exists {
			return nil, fmt.Errorf("duplicate default role snapshot %q for version %q", definition.APILabel, version)
		}
		seen[definition.APILabel] = struct{}{}
		definitions = append(definitions, definition)
	}

	return definitions, nil
}

func collectDefaultPermissionLabels(definitions []roleconstants.DefaultRole) []string {
	seen := make(map[string]struct{})
	labels := make([]string, 0)
	for _, role := range definitions {
		for _, label := range role.PermissionLabels {
			if _, exists := seen[label]; exists {
				continue
			}
			seen[label] = struct{}{}
			labels = append(labels, label)
		}
	}
	return labels
}

func verifyDefaultRoleSnapshot(snapshot rolemodel.DefaultRoleSnapshot) error {
	var raw any
	if err := json.Unmarshal(snapshot.Definition, &raw); err != nil {
		return fmt.Errorf("decode default role snapshot %q for integrity check: %w", snapshot.RoleAPILabel, err)
	}

	canonical, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("canonicalize default role snapshot %q: %w", snapshot.RoleAPILabel, err)
	}

	sum := sha256.Sum256(canonical)
	actual := hex.EncodeToString(sum[:])
	if actual != snapshot.DefinitionSHA256 {
		return fmt.Errorf("default role snapshot %q failed integrity verification", snapshot.RoleAPILabel)
	}

	return nil
}
