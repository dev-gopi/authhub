CREATE TABLE user_pool_auth_policies (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    user_pool_id uuid NOT NULL
        REFERENCES user_pools(id)
        ON DELETE CASCADE,

    username_login_enabled boolean NOT NULL DEFAULT true,

    email_login_enabled boolean NOT NULL DEFAULT false,

    password_login_enabled boolean NOT NULL DEFAULT true,

    passkey_login_enabled boolean NOT NULL DEFAULT true,

    registration_enabled boolean NOT NULL DEFAULT true,

    email_verification_required boolean NOT NULL DEFAULT true,

    mfa_policy text NOT NULL DEFAULT 'optional',

    totp_enabled boolean NOT NULL DEFAULT true,

    session_idle_seconds integer NOT NULL,

    session_absolute_seconds integer NOT NULL,

    max_active_sessions integer,

    recent_auth_seconds integer NOT NULL,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT user_pool_auth_policy_unique
        UNIQUE (user_pool_id)
);