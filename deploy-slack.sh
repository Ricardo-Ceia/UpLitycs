#!/bin/bash
# Slack Integration Deployment Script
# This script helps you deploy the Slack integration to your database

set -e

echo "🚀 StatusFrame Slack Integration Deployment"
echo "==========================================="
echo ""

# Check if migrations file exists
if [ ! -f "db/migrations/add_slack_integration.sql" ]; then
    echo "❌ Migration file not found at db/migrations/add_slack_integration.sql"
    exit 1
fi

echo "📝 Migration file found"
echo ""

# Prompt for database connection details
read -p "Enter PostgreSQL username (default: postgres): " DB_USER
DB_USER=${DB_USER:-postgres}

read -p "Enter database name (default: uplitycs): " DB_NAME
DB_NAME=${DB_NAME:-uplitycs}

read -p "Enter PostgreSQL host (default: localhost): " DB_HOST
DB_HOST=${DB_HOST:-localhost}

read -p "Enter PostgreSQL port (default: 5432): " DB_PORT
DB_PORT=${DB_PORT:-5432}

echo ""
echo "🔗 Connecting to database: $DB_USER@$DB_HOST:$DB_PORT/$DB_NAME"
echo ""

# Run migration
echo "⏳ Running migration..."
PGPASSWORD="" psql -U "$DB_USER" -h "$DB_HOST" -p "$DB_PORT" -d "$DB_NAME" -f "db/migrations/add_slack_integration.sql"

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ Slack integration database tables created successfully!"
    echo ""
    echo "📋 Tables created:"
    echo "   - slack_integrations"
    echo "   - incident_notifications"
    echo ""
    echo "🎯 Next steps:"
    echo "   1. Update frontend - add SlackIntegration component to settings page"
    echo "   2. Configure environment variables in .env or deployment config"
    echo "   3. Restart backend service"
    echo "   4. Test OAuth flow with a Pro/Business plan user"
    echo ""
else
    echo "❌ Migration failed. Please check the error above."
    exit 1
fi

echo "📚 For more information, see:"
echo "   - SLACK_INTEGRATION.md (User guide)"
echo "   - SLACK_IMPLEMENTATION.md (Developer guide)"
echo ""
echo "Done! 🎉"
