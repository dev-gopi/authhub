CREATE UNIQUE INDEX IF NOT EXISTS tenant_one_primary_admin_idx
ON tenant_members(tenant_id)
WHERE member_class = 'primary_admin'
  AND status = 'active'
  AND is_active = true
  AND is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_roles_tenant_active
ON roles(tenant_id, api_label)
WHERE is_active = true AND is_deleted = false;

CREATE INDEX IF NOT EXISTS idx_role_permissions_active
ON role_permissions(role_id, permission_id)
WHERE is_active = true AND is_deleted = false;
