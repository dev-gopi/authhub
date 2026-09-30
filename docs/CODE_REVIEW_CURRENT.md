# Current Code Review Notes

This review covers the current runtime and Task 2-5 implementation after the RBAC refactor.

## Checked and aligned

- Fresh migration chain `000001` through `000024` remains ordered correctly.
- Every application table created by the migration chain has the shared base fields: `id`, `created_at`, `updated_at`, `created_by`, `updated_by`, `deleted_at`, `is_active`, `is_deleted`.
- `roles.metadata` and `permissions.metadata` are JSONB and indexed for future catalog/configuration queries.
- 80 default permissions are seeded at the permission-table foundation migration.
- Eight default role templates are seeded as versioned snapshots (`v1`).
- Tenant provisioning reads the configured snapshot version instead of depending only on compiled role mappings.
- Default role snapshot integrity hashes are validated before tenant roles are materialized.
- Missing required permissions fail tenant provisioning closed instead of creating partial roles.
- Root Admin authorization remains separate from tenant roles.
- Generic unfinished authentication/permission/tenant middleware now fails closed rather than silently allowing requests.
- Secret redaction helper no longer returns the secret input.
- Security headers are enabled on both Auth API and Control API.

## Intentionally unfinished/future modules

The following areas still contain scaffold files and are not treated as complete in the current milestone:

- worker/outbox relay runtime
- realtime gateway
- generic tenant-user authentication middleware
- tenant permission middleware/evaluator
- role CRUD HTTP APIs
- permission catalog HTTP APIs
- user-pool CRUD HTTP APIs
- platform-user CRUD HTTP APIs
- tenant-member CRUD HTTP APIs
- RabbitMQ topology setup beyond the current connection layer
- later OAuth/OIDC/SAML/WebAuthn/TOTP/session/device/recovery modules

These scaffolds should not be considered production implementations until their corresponding milestones are completed.

## Local verification

The project declares Go `1.27.1`. The review environment has Go `1.23.x`, so full `go test ./...` / `go vet ./...` verification cannot be claimed here without changing the project's required toolchain. Source formatting, migration consistency, RBAC catalog integrity, snapshot hashes, and pure RBAC constant tests were checked separately.
