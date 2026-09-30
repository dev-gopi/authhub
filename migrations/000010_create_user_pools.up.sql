CREATE TABLE user_pools (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    tenant_id uuid NOT NULL
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    api_label text NOT NULL,

    display_name text NOT NULL,

    status text NOT NULL,

    issuer_uri text NOT NULL,

    default_locale text NOT NULL DEFAULT 'en',

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT user_pools_tenant_api_label_unique
        UNIQUE (tenant_id, api_label)
);