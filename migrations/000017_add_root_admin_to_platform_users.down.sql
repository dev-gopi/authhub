DROP INDEX IF EXISTS idx_platform_users_root_admin;

ALTER TABLE platform_users
DROP COLUMN IF EXISTS is_root_admin;