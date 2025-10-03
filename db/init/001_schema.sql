-- Schema initialization
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  avatar_url TEXT,
  health_url TEXT,
  theme TEXT DEFAULT 'cyberpunk',
  alerts TEXT,
  slug TEXT UNIQUE,
  app_name TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

-- Health check results (status derived from status_code)
CREATE TABLE IF NOT EXISTS user_status (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    status_code INT NOT NULL,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_user_status_user_id ON user_status(user_id);
CREATE INDEX IF NOT EXISTS idx_user_status_checked_at ON user_status(checked_at);

-- Alert history for rate limiting (only stores when alerts were sent)
CREATE TABLE IF NOT EXISTS alerts (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_alerts_user_id ON alerts(user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_sent_at ON alerts(sent_at);