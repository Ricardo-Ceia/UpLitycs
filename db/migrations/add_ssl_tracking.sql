-- Add SSL certificate tracking columns to apps table
ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_expiry_date TIMESTAMPTZ;
ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_days_until_expiry INTEGER;
ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_issuer TEXT;
ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_last_checked TIMESTAMPTZ;
