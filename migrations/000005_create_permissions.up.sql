CREATE TABLE permissions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),

    api_label text NOT NULL UNIQUE,

    display_name text NOT NULL,

    description text,

    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz,

    is_active boolean NOT NULL DEFAULT true,
    is_deleted boolean NOT NULL DEFAULT false
);