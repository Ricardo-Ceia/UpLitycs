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
  app_name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  health_url TEXT NOT NULL,
  theme TEXT DEFAULT 'cyberpunk',
  alerts TEXT DEFAULT 'n',
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now(),
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