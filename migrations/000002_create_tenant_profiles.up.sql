CREATE TABLE tenant_profiles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    tenant_id uuid NOT NULL
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    legal_name text,
    logo_asset_id uuid,

    website_url text,

    support_email text,
    support_phone text,

    security_contact_email text,

    locale text,
    timezone text,

    country_code text,
    region text,

    privacy_policy_url text,
    terms_url text,

    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT tenant_profiles_tenant_unique
        UNIQUE (tenant_id)
);