CREATE TABLE roles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    tenant_id uuid
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    api_label text NOT NULL,

    display_name text NOT NULL,

    description text,

    is_default boolean NOT NULL DEFAULT false,

    protected_from_delete boolean NOT NULL DEFAULT false,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT roles_tenant_api_label_unique
        UNIQUE (tenant_id, api_label)
);