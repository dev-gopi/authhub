CREATE TABLE role_permissions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    role_id uuid NOT NULL
        REFERENCES roles(id)
        ON DELETE CASCADE,

    permission_id uuid NOT NULL
        REFERENCES permissions(id)
        ON DELETE CASCADE,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT role_permissions_unique
        UNIQUE (role_id, permission_id)
);