CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  homepage TEXT NOT NULL,
  alerts TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);

INSERT INTO users (username, homepage, alerts)
VALUES
 ('openai','https://openai.com','yes'),
 ('github','https://github.com','yes'),
 ('wikipedia','https://www.wikipedia.org','no')
ON CONFLICT (username) DO NOTHING;