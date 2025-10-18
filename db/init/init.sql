-- Final canonical schema (users = user-level data; apps = monitored sites)
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  avatar_url TEXT,
  plan TEXT DEFAULT 'free',
  plan_started_at TIMESTAMPTZ DEFAULT now(),
  stripe_customer_id TEXT,
  stripe_subscription_id TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_plan ON users(plan);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

CREATE TABLE IF NOT EXISTS apps (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  logo_url TEXT,
  app_name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  health_url TEXT NOT NULL,
  theme TEXT DEFAULT 'cyberpunk',
  alerts TEXT DEFAULT 'n',
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
  next_check_at TIMESTAMPTZ DEFAULT now(),
  ssl_expiry_date TIMESTAMPTZ,
  ssl_days_until_expiry INTEGER,
  ssl_issuer TEXT,
  ssl_last_checked TIMESTAMPTZ,
  UNIQUE(user_id, app_name)
);

CREATE INDEX IF NOT EXISTS idx_apps_user_id ON apps(user_id);
CREATE INDEX IF NOT EXISTS idx_apps_slug ON apps(slug);

CREATE TABLE IF NOT EXISTS user_status (
  id SERIAL PRIMARY KEY,
  app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
  status_code INTEGER NOT NULL,
  checked_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_user_status_app_id ON user_status(app_id);
CREATE INDEX IF NOT EXISTS idx_user_status_checked_at ON user_status(checked_at);

CREATE TABLE IF NOT EXISTS alerts (
  id SERIAL PRIMARY KEY,
  app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
  sent_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_alerts_app_id ON alerts(app_id);
CREATE INDEX IF NOT EXISTS idx_alerts_sent_at ON alerts(sent_at);

-- Slack integration table
CREATE TABLE IF NOT EXISTS slack_integrations (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  slack_team_id VARCHAR(255) NOT NULL,
  slack_team_name VARCHAR(255),
  slack_bot_token VARCHAR(1000) NOT NULL,
  slack_channel_id VARCHAR(255),
  slack_channel_name VARCHAR(255),
  is_enabled BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_slack_integrations_user_id ON slack_integrations(user_id);

-- Discord integration table
CREATE TABLE IF NOT EXISTS discord_integrations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    discord_user_id VARCHAR(255) NOT NULL,
    discord_username VARCHAR(255),
    webhook_url VARCHAR(1000) NOT NULL,
    server_id VARCHAR(255),
    server_name VARCHAR(255),
    channel_id VARCHAR(255),
    channel_name VARCHAR(255),
    is_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_discord_integrations_user_id ON discord_integrations(user_id);
CREATE INDEX IF NOT EXISTS idx_discord_integrations_discord_user_id ON discord_integrations(discord_user_id);

-- Incident notifications tracking
CREATE TABLE IF NOT EXISTS incident_notifications (
  id SERIAL PRIMARY KEY,
  app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
  notification_type VARCHAR(50) NOT NULL,
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_incident_notifications_app_id ON incident_notifications(app_id);