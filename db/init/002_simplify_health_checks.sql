-- Migration: Simplify health check schema
-- Remove redundant columns and derive status from status_code

-- Simplify user_status table
ALTER TABLE user_status DROP COLUMN IF EXISTS page;
ALTER TABLE user_status DROP COLUMN IF EXISTS status;
ALTER TABLE user_status DROP COLUMN IF EXISTS response_time_ms;

-- Simplify alerts table
ALTER TABLE alerts DROP COLUMN IF EXISTS status;
ALTER TABLE alerts DROP COLUMN IF EXISTS status_code;

-- Add comment to document the simplified schema
COMMENT ON TABLE user_status IS 'Stores health check results. Status is derived from status_code: 200-299=up, 300-399=degraded, 400-499=client_error, 500+=down, 0=error';
COMMENT ON TABLE alerts IS 'Stores when downtime alerts were sent to users for rate limiting';
