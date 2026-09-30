CREATE TABLE user_profiles (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id uuid NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    job_title text,

    department text,

    profile_picture_asset_id uuid,

    organization_details jsonb NOT NULL DEFAULT '{}'::jsonb,

    custom_profile jsonb NOT NULL DEFAULT '{}'::jsonb,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by uuid,
    updated_by uuid,
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false,

    CONSTRAINT user_profiles_user_unique
        UNIQUE (user_id)
);