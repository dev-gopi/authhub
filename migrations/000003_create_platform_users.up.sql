CREATE TABLE platform_users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    username text NOT NULL UNIQUE,
    email text NOT NULL UNIQUE,

    email_verified_at timestamptz,

    display_name text,

    status text NOT NULL,

    credential_version bigint NOT NULL DEFAULT 1,

    last_login_at timestamptz,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT platform_users_status_check
        CHECK (
            status IN (
                'active',
                'disabled',
                'blocked'
            )
        )
);