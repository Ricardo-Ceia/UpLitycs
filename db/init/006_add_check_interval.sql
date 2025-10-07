-- Add check_interval column to apps table
-- This allows each app to have custom check intervals based on plan limits
ALTER TABLE apps ADD COLUMN IF NOT EXISTS check_interval INT DEFAULT 300;

-- Update existing apps to have default intervals based on current user plan
UPDATE apps a
SET check_interval = CASE 
    WHEN u.plan = 'business' THEN 30
    WHEN u.plan = 'pro' THEN 60
    ELSE 300
END
FROM users u
WHERE a.user_id = u.id;

-- Add comment for documentation
COMMENT ON COLUMN apps.check_interval IS 'Health check interval in seconds. Min values: free=300, pro=60, business=30';

-- Create index for interval-based queries (useful for scheduling)
CREATE INDEX IF NOT EXISTS idx_apps_check_interval ON apps(check_interval);
