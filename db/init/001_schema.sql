-- Schema initialization
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  avatar_url TEXT,
  health_url TEXT,
  theme TEXT DEFAULT 'cyberpunk',
  alerts TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_status (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    page TEXT NOT NULL,
    status TEXT NOT NULL,
    status_code INT NOT NULL,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT now()
);