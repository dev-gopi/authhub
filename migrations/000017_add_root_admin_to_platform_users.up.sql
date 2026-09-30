ALTER TABLE platform_users
ADD COLUMN is_root_admin boolean NOT NULL DEFAULT false;

CREATE INDEX idx_platform_users_root_admin
ON platform_users (is_root_admin)
WHERE is_root_admin = true
  AND is_active = true
  AND is_deleted = false;