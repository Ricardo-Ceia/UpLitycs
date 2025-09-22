-- Initial seed data
INSERT INTO users (username, homepage, alerts)
VALUES
  ('openai','https://openai.com','yes'),
  ('github','https://github.com','yes'),
  ('wikipedia','https://www.wikipedia.org','no'),
  ('trpgenie','https://trpgenie.com','yes')
ON CONFLICT (username) DO NOTHING;
