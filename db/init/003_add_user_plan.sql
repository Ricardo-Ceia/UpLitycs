-- Add plan column to users table for subscription management
ALTER TABLE users ADD COLUMN IF NOT EXISTS plan TEXT DEFAULT 'free';

-- Add plan_started_at to track when current plan began (for trials)
ALTER TABLE users ADD COLUMN IF NOT EXISTS plan_started_at TIMESTAMPTZ DEFAULT now();

-- Add stripe_customer_id for payment integration
ALTER TABLE users ADD COLUMN IF NOT EXISTS stripe_customer_id TEXT;

-- Add stripe_subscription_id to track active subscriptions
ALTER TABLE users ADD COLUMN IF NOT EXISTS stripe_subscription_id TEXT;

-- Create index for faster plan queries
CREATE INDEX IF NOT EXISTS idx_users_plan ON users(plan);
