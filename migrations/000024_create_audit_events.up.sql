CREATE TABLE audit_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid REFERENCES tenants(id) ON DELETE SET NULL,
    actor_platform_user_id uuid REFERENCES platform_users(id) ON DELETE SET NULL,
    action text NOT NULL,
    resource_type text NOT NULL,
    resource_id uuid,
    success boolean NOT NULL,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,
    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false
);
CREATE INDEX idx_audit_events_tenant_created_at ON audit_events(tenant_id, created_at DESC);
CREATE INDEX idx_audit_events_actor_created_at ON audit_events(actor_platform_user_id, created_at DESC);
