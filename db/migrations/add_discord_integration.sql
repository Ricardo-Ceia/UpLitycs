-- Add Discord integration table
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

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_discord_integrations_user_id ON discord_integrations(user_id);
CREATE INDEX IF NOT EXISTS idx_discord_integrations_discord_user_id ON discord_integrations(discord_user_id);
