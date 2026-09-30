CREATE TABLE tenant_member_roles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    tenant_member_id uuid NOT NULL
        REFERENCES tenant_members(id)
        ON DELETE CASCADE,

    role_id uuid NOT NULL
        REFERENCES roles(id)
        ON DELETE CASCADE,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT tenant_member_roles_unique
        UNIQUE (tenant_member_id, role_id)
);