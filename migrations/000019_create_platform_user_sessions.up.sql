CREATE TABLE platform_user_sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    platform_user_id uuid NOT NULL
        REFERENCES platform_users(id)
        ON DELETE CASCADE,

    session_token_hash char(64) NOT NULL,

    credential_version bigint NOT NULL,

    ip_address inet,

    user_agent text,

    expires_at timestamptz NOT NULL,

    idle_expires_at timestamptz NOT NULL,

    last_seen_at timestamptz NOT NULL DEFAULT now(),

    revoked_at timestamptz,

    revoked_reason text,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT platform_user_sessions_token_hash_unique
        UNIQUE (session_token_hash)
);

CREATE INDEX idx_platform_user_sessions_platform_user
ON platform_user_sessions(platform_user_id);

CREATE INDEX idx_platform_user_sessions_active
ON platform_user_sessions(platform_user_id, expires_at)
WHERE revoked_at IS NULL
  AND is_active = true
  AND is_deleted = false;