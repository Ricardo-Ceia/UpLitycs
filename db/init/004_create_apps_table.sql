-- Create apps table for multi-app support
CREATE TABLE IF NOT EXISTS apps (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    app_name TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    health_url TEXT NOT NULL,
    theme TEXT DEFAULT 'cyberpunk',
    alerts TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_apps_user_id ON apps(user_id);
CREATE INDEX IF NOT EXISTS idx_apps_slug ON apps(slug);

-- Migrate existing user data to apps table
INSERT INTO apps (user_id, app_name, slug, health_url, theme, alerts, created_at)
SELECT id, app_name, slug, health_url, theme, alerts, created_at 
FROM users 
WHERE slug IS NOT NULL AND slug != '' AND health_url IS NOT NULL AND health_url != '';

-- Update user_status to reference apps instead of users directly
ALTER TABLE user_status ADD COLUMN IF NOT EXISTS app_id INT REFERENCES apps(id) ON DELETE CASCADE;

-- Migrate existing status checks to use app_id
UPDATE user_status us
SET app_id = a.id
FROM apps a
WHERE us.user_id = a.user_id;

-- Create index for app_id lookups
CREATE INDEX IF NOT EXISTS idx_user_status_app_id ON user_status(app_id);

-- Update alerts table to reference apps
ALTER TABLE alerts ADD COLUMN IF NOT EXISTS app_id INT REFERENCES apps(id) ON DELETE CASCADE;

-- Migrate existing alerts to use app_id  
UPDATE alerts al
SET app_id = a.id
FROM apps a
WHERE al.user_id = a.user_id;

-- Create index for app_id in alerts
CREATE INDEX IF NOT EXISTS idx_alerts_app_id ON alerts(app_id);

-- We'll keep user_id in both tables for backwards compatibility during migration
-- Can drop old columns later: ALTER TABLE users DROP COLUMN health_url, slug, app_name, theme, alerts;
