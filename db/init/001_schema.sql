-- Schema initialization
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  homepage TEXT NOT NULL,
  alerts TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);
