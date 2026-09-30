CREATE TABLE password_policies (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    user_pool_id uuid NOT NULL
        REFERENCES user_pools(id)
        ON DELETE CASCADE,

    min_length integer NOT NULL,

    max_length integer NOT NULL,

    prevent_common_passwords boolean NOT NULL DEFAULT true,

    breached_password_check_enabled boolean NOT NULL DEFAULT false,

    password_history_count integer NOT NULL DEFAULT 0,

    reset_token_ttl_seconds integer NOT NULL,

    lockout_policy jsonb NOT NULL DEFAULT '{}'::jsonb,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT password_policy_user_pool_unique
        UNIQUE (user_pool_id)
);