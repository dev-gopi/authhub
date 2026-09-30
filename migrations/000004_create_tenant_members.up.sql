CREATE TABLE tenant_members (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    tenant_id uuid NOT NULL
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    platform_user_id uuid NOT NULL
        REFERENCES platform_users(id)
        ON DELETE CASCADE,

    member_class text NOT NULL,

    status text NOT NULL,

    protected_admin boolean NOT NULL DEFAULT false,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT tenant_members_unique
        UNIQUE (tenant_id, platform_user_id),

    CONSTRAINT tenant_members_class_check
        CHECK (
            member_class IN (
                'primary_admin',
                'delegated_admin',
                'standard_user'
            )
        )
);