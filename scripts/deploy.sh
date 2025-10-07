#!/bin/bash

# StatusFrame Production Deployment Script
# This script builds and deploys StatusFrame to production

set -e  # Exit on any error

echo "🚀 StatusFrame Production Deployment"
echo "===================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_DIR="/home/ubuntu/statusframe"
BACKEND_BINARY="statusframe"
SYSTEMD_SERVICE="statusframe"

# Check if running as root or with sudo
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}Please run with sudo${NC}"
    exit 1
fi

echo -e "${YELLOW}Step 1: Stopping existing service...${NC}"
systemctl stop $SYSTEMD_SERVICE || true

echo -e "${GREEN}✓ Service stopped${NC}"
echo ""

echo -e "${YELLOW}Step 2: Building frontend...${NC}"
cd $APP_DIR/frontend
npm install
npm run build

if [ ! -d "dist" ]; then
    echo -e "${RED}✗ Frontend build failed - dist folder not found${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Frontend built successfully${NC}"
echo ""

echo -e "${YELLOW}Step 3: Building Go backend...${NC}"
cd $APP_DIR
go build -o $BACKEND_BINARY main.go

if [ ! -f "$BACKEND_BINARY" ]; then
    echo -e "${RED}✗ Backend build failed - binary not found${NC}"
    exit 1
fi

# Make binary executable
chmod +x $BACKEND_BINARY

echo -e "${GREEN}✓ Backend built successfully${NC}"
echo ""

echo -e "${YELLOW}Step 4: Updating environment...${NC}"
if [ ! -f ".env" ]; then
    echo -e "${RED}✗ .env file not found!${NC}"
    echo "Please create .env file with required configuration"
    exit 1
fi

echo -e "${GREEN}✓ Environment configured${NC}"
echo ""

echo -e "${YELLOW}Step 5: Database migrations...${NC}"
# Ensure PostgreSQL is running
systemctl status postgresql || systemctl start postgresql

echo -e "${GREEN}✓ Database ready${NC}"
echo ""

echo -e "${YELLOW}Step 6: Starting service...${NC}"
systemctl daemon-reload
systemctl start $SYSTEMD_SERVICE
systemctl enable $SYSTEMD_SERVICE

# Wait a moment for service to start
sleep 2

# Check if service is running
if systemctl is-active --quiet $SYSTEMD_SERVICE; then
    echo -e "${GREEN}✓ Service started successfully${NC}"
else
    echo -e "${RED}✗ Service failed to start${NC}"
    echo "Check logs with: journalctl -u $SYSTEMD_SERVICE -n 50"
    exit 1
fi

echo ""
echo -e "${YELLOW}Step 7: Reloading Nginx...${NC}"
nginx -t && systemctl reload nginx

echo -e "${GREEN}✓ Nginx reloaded${NC}"
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}🎉 Deployment completed successfully!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Service status:"
systemctl status $SYSTEMD_SERVICE --no-pager -l
echo ""
echo "Useful commands:"
echo "  View logs: journalctl -u $SYSTEMD_SERVICE -f"
echo "  Restart: sudo systemctl restart $SYSTEMD_SERVICE"
echo "  Stop: sudo systemctl stop $SYSTEMD_SERVICE"
echo ""
