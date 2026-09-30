CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    tenant_id uuid NOT NULL
        REFERENCES tenants(id)
        ON DELETE CASCADE,

    user_pool_id uuid NOT NULL
        REFERENCES user_pools(id)
        ON DELETE CASCADE,

    username text,
    normalized_username text,

    email text,
    normalized_email text,

    email_verified_at timestamptz,

    display_name text,

    nickname text,

    given_name text,

    family_name text,

    locale text,

    zoneinfo text,

    status text NOT NULL,

    credential_version bigint NOT NULL DEFAULT 1,

    must_change_password boolean NOT NULL DEFAULT false,

    last_login_at timestamptz,

    disabled_at timestamptz,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX users_pool_username_unique
ON users(user_pool_id, normalized_username)
WHERE normalized_username IS NOT NULL
  AND is_deleted = false;

CREATE UNIQUE INDEX users_pool_email_unique
ON users(user_pool_id, normalized_email)
WHERE normalized_email IS NOT NULL
  AND is_deleted = false;