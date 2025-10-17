-- Add Slack integration table
CREATE TABLE IF NOT EXISTS slack_integrations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    app_id INTEGER REFERENCES apps(id) ON DELETE CASCADE,
    slack_team_id VARCHAR(255) NOT NULL,
    slack_team_name VARCHAR(255),
    slack_bot_token VARCHAR(500) NOT NULL,
    slack_channel_id VARCHAR(255) NOT NULL,
    slack_channel_name VARCHAR(255),
    is_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add incident notifications table (optional, for tracking sent alerts)
CREATE TABLE IF NOT EXISTS incident_notifications (
    id SERIAL PRIMARY KEY,
    app_id INTEGER NOT NULL REFERENCES apps(id) ON DELETE CASCADE,
    notification_type VARCHAR(50) NOT NULL, -- 'slack', 'email', etc
    status VARCHAR(50), -- 'sent', 'failed', 'pending'
    message_timestamp VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_slack_integrations_user_id ON slack_integrations(user_id);
CREATE INDEX IF NOT EXISTS idx_slack_integrations_app_id ON slack_integrations(app_id);
CREATE INDEX IF NOT EXISTS idx_incident_notifications_app_id ON incident_notifications(app_id);
