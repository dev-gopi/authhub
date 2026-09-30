CREATE TABLE platform_user_passwords (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    platform_user_id uuid NOT NULL
        REFERENCES platform_users(id)
        ON DELETE CASCADE,

    password_hash text NOT NULL,

    password_algorithm text NOT NULL
        DEFAULT 'argon2id',

    password_params jsonb NOT NULL
        DEFAULT '{}'::jsonb,

    password_version integer NOT NULL DEFAULT 1,

    changed_at timestamptz NOT NULL DEFAULT now(),

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT platform_user_passwords_user_unique
        UNIQUE (platform_user_id),

    CONSTRAINT platform_user_password_algorithm_check
        CHECK (password_algorithm IN ('argon2id')),

    CONSTRAINT platform_user_password_version_check
        CHECK (password_version > 0)
);

CREATE INDEX idx_platform_user_passwords_platform_user_id
ON platform_user_passwords(platform_user_id);