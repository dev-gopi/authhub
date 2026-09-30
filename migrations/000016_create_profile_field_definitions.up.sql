CREATE TABLE profile_field_definitions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    tenant_id uuid NOT NULL
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    user_pool_id uuid NOT NULL
        REFERENCES user_pools(id)
        ON DELETE CASCADE,

    api_label text NOT NULL,

    display_label text NOT NULL,

    field_type text NOT NULL,

    user_readable boolean NOT NULL DEFAULT true,

    user_editable boolean NOT NULL DEFAULT true,

    admin_editable boolean NOT NULL DEFAULT true,

    required boolean NOT NULL DEFAULT false,

    validation jsonb NOT NULL DEFAULT '{}'::jsonb,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT profile_field_definition_unique
        UNIQUE (
            user_pool_id,
            api_label
        )
);