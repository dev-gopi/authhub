CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE tenants (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    api_label text NOT NULL UNIQUE,
    display_name text NOT NULL,

    status text NOT NULL,

    default_locale text NOT NULL DEFAULT 'en',

    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

    suspended_at timestamptz,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT tenants_status_check
        CHECK (
            status IN (
                'provisioning',
                'active',
                'suspended',
                'delete_requested',
                'deleted'
            )
        )
);