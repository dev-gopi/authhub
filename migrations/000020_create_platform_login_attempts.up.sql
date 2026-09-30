CREATE TABLE platform_login_attempts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    platform_user_id uuid
        REFERENCES platform_users(id)
        ON DELETE SET NULL,

    identifier_hash char(64) NOT NULL,

    ip_address inet,

    user_agent text,

    success boolean NOT NULL DEFAULT false,

    failure_reason text,

    attempted_at timestamptz NOT NULL DEFAULT now(),

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE INDEX idx_platform_login_attempts_user
ON platform_login_attempts(platform_user_id, attempted_at DESC);

CREATE INDEX idx_platform_login_attempts_identifier
ON platform_login_attempts(identifier_hash, attempted_at DESC);

CREATE INDEX idx_platform_login_attempts_ip
ON platform_login_attempts(ip_address, attempted_at DESC);