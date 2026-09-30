CREATE TABLE default_role_snapshots (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    role_api_label text NOT NULL,

    platform_version text NOT NULL,

    definition jsonb NOT NULL,

    definition_sha256 text NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT default_role_snapshots_unique
        UNIQUE (
            role_api_label,
            platform_version
        )
);