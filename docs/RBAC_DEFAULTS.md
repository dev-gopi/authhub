# AuthHub Default RBAC Catalog

This document describes the default tenant RBAC configuration used when a new tenant is provisioned.

## Design rules

- Permissions are global, immutable API labels stored in `permissions`.
- Tenant roles are materialized per tenant in `roles` and linked through `role_permissions`.
- Root Admin authorization is separate and is **not** represented by tenant roles.
- `primary_admin` is the protected owner-administrator role for one tenant and is not generally assignable.
- Role templates are versioned in `default_role_snapshots`.
- New tenant provisioning selects a template version using `DEFAULT_ROLE_TEMPLATE_VERSION` (default `v1`).
- Existing tenant roles are not silently rewritten when a new template version is introduced. Future upgrades should use an explicit migration/synchronization workflow.
- Every permission and role has JSONB metadata so future UI/policy engines can group, filter, classify, and evolve the catalog without changing API labels.

## Permission metadata

Every seeded permission includes metadata such as catalog version, system ownership, immutable API label, category, resource, action, scope, risk level, and tenant-role assignability.

## Permission groups

### Tenant
`tenant.read`, `tenant.update`

### Identity and membership
`user.read`, `user.create`, `user.update`, `user.disable`, `user.block`, `user.delete`, `user.invite`, `user_pool.read`, `user_pool.create`, `user_pool.update`, `user_pool.delete`, `profile_field.read`, `profile_field.create`, `profile_field.update`, `profile_field.delete`

### RBAC
`role.read`, `role.create`, `role.update`, `role.delete`, `role.assign`, `permission.read`

### Policies
`auth_policy.read`, `auth_policy.update`, `password_policy.read`, `password_policy.update`, `recovery_policy.read`, `recovery_policy.update`, `login_policy.read`, `login_policy.update`

### Sessions, devices, and MFA
`session.read`, `session.revoke`, `device.read`, `device.revoke`, `mfa.read`, `mfa.reset`, `totp.read`, `totp.reset`, `webauthn.read`, `webauthn.revoke`

### Applications and federation
`application.read`, `application.create`, `application.update`, `application.disable`, `application.delete`, `application.secret.rotate`, `oauth.read`, `oauth.revoke`, `oidc.read`, `oidc.update`, `saml.read`, `saml.update`

### Communications and integrations
`email.read`, `email.update`, `sms.read`, `sms.update`, `notification.read`, `notification.update`, `webhook.read`, `webhook.update`, `webhook.replay`, `post_auth_hook.read`, `post_auth_hook.create`, `post_auth_hook.update`, `post_auth_hook.delete`, `realtime.read`, `realtime.update`

### Key management, audit, and support
`key.read`, `key.rotate`, `audit.read`, `support.read`, `support.send`

### End-user self service
`self.profile.read`, `self.profile.update`, `self.session.read`, `self.session.revoke`, `self.device.read`, `self.device.revoke`, `self.mfa.manage`

## Default tenant roles

### `primary_admin`
Protected tenant owner-administrator. Receives every tenant-administrative permission. It is not generally assignable and remains distinct from Root Admin.

### `tenant_admin`
Broad tenant administration. It receives most tenant administrative permissions but intentionally excludes ownership-sensitive key rotation and direct application-secret rotation by default.

### `user_manager`
Identity operations: user lifecycle, session/device revocation, MFA/TOTP/passkey recovery operations, profile-field visibility, and bounded support visibility.

### `application_manager`
Application and integration operations: application clients, OAuth, OIDC, SAML, webhooks, post-auth hooks, realtime configuration, and application-secret rotation.

### `security_viewer`
Read-only security visibility across policies, sessions, devices, MFA state, federation, key metadata, and audit events.

### `support_agent`
Bounded user-support operations including user lookup, session/device revocation, MFA recovery operations, and support actions.

### `viewer`
General read-only tenant administration. It intentionally does not include cryptographic key metadata.

### `standard_user`
Only self-service permissions for the authenticated user's own profile, sessions, devices, and MFA.

## Future role-template changes

Do not mutate `v1` snapshots after production tenants rely on them. Instead:

1. Add new `default_role_snapshots` rows for `v2`.
2. Validate every referenced permission exists in the global catalog.
3. Set `DEFAULT_ROLE_TEMPLATE_VERSION=v2` for newly provisioned tenants.
4. Migrate existing tenants only through an explicit upgrade process with audit logging.

This keeps tenant provisioning deterministic and prevents a deployment from silently changing existing customer authorization.
