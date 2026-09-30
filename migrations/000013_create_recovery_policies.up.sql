CREATE TABLE recovery_policies (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    user_pool_id uuid NOT NULL
        REFERENCES user_pools(id)
        ON DELETE CASCADE,

    email_link_enabled boolean NOT NULL DEFAULT true,

    email_otp_enabled boolean NOT NULL DEFAULT true,

    sms_otp_enabled boolean NOT NULL DEFAULT false,

    saved_recovery_code_enabled boolean NOT NULL DEFAULT true,

    existing_totp_enabled boolean NOT NULL DEFAULT true,

    existing_passkey_enabled boolean NOT NULL DEFAULT true,

    admin_assisted_enabled boolean NOT NULL DEFAULT true,

    security_question_enabled boolean NOT NULL DEFAULT false,

    minimum_proofs integer NOT NULL DEFAULT 1,

    require_independent_channels boolean NOT NULL DEFAULT false,

    post_recovery_session_action text NOT NULL DEFAULT 'revoke_all',

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT recovery_policy_user_pool_unique
        UNIQUE (user_pool_id)
);