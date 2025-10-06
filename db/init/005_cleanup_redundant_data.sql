-- Phase 1: Remove redundant columns from users table
-- These columns are now stored in the apps table and are no longer needed

-- Drop redundant app-related columns from users table
ALTER TABLE users DROP COLUMN IF EXISTS health_url;
ALTER TABLE users DROP COLUMN IF EXISTS theme;
ALTER TABLE users DROP COLUMN IF EXISTS alerts;
ALTER TABLE users DROP COLUMN IF EXISTS slug;
ALTER TABLE users DROP COLUMN IF EXISTS app_name;

-- Phase 2: Remove redundant user_id foreign keys
-- We can always get user via app_id -> apps.user_id, so user_id is redundant

-- Drop user_id from user_status (keep app_id only)
ALTER TABLE user_status DROP COLUMN IF EXISTS user_id;

-- Drop user_id from alerts (keep app_id only)
ALTER TABLE alerts DROP COLUMN IF EXISTS user_id;

-- Drop old indexes on removed columns
DROP INDEX IF EXISTS idx_user_status_user_id;
DROP INDEX IF EXISTS idx_alerts_user_id;

-- The following indexes are already created and will remain:
-- idx_user_status_app_id (for user_status.app_id)
-- idx_alerts_app_id (for alerts.app_id)
-- These are sufficient for efficient queries
