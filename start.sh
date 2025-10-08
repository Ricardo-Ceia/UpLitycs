#!/bin/bash

echo "🚀 Starting StatusFrame..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ .env file not found!"
    echo "📝 Copy .env.example to .env and fill in your values:"
    echo "   cp .env.example .env"
    echo "   nano .env"
    exit 1
fi

# Load environment
export $(grep -v '^#' .env | xargs)

# Pull latest code (optional)
# git pull

# Build and start
docker compose down
docker compose build --no-cache
docker compose up -d

# Wait for services
echo "⏳ Waiting for services to start..."
sleep 15

# Check status
echo ""
echo "📊 Service Status:"
docker compose ps

echo ""
echo "✅ StatusFrame is running!"
echo ""
echo "🌐 Access your app at: https://${DOMAIN}"
echo ""
echo "📝 View logs: docker compose logs -f"
echo "🛑 Stop: docker compose down"