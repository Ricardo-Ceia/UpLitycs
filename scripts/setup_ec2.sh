#!/bin/bash

# UpLitycs EC2 Setup Script
# This script checks for required dependencies and installs them if missing

set -e

echo "🚀 UpLitycs EC2 Setup Script"
echo "=============================="
echo ""

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print colored output
print_error() {
    echo -e "${RED}❌ ERROR: $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  WARNING: $1${NC}"
}

print_info() {
    echo "ℹ️  $1"
}

# Check for Docker
echo "Checking for Docker..."
if ! command_exists docker; then
    print_error "Docker is not installed!"
    echo ""
    echo "Please install Docker first:"
    echo "  curl -fsSL https://get.docker.com -o get-docker.sh"
    echo "  sudo sh get-docker.sh"
    echo "  sudo usermod -aG docker \$USER"
    echo "  newgrp docker"
    echo ""
    echo "See EC2_SETUP.md for detailed instructions."
    exit 1
else
    DOCKER_VERSION=$(docker --version)
    print_success "Docker found: $DOCKER_VERSION"
fi

# Check for Docker Compose
echo "Checking for Docker Compose..."
if ! docker compose version >/dev/null 2>&1; then
    print_error "Docker Compose is not installed!"
    echo ""
    echo "Please install Docker Compose plugin:"
    echo "  sudo apt-get update"
    echo "  sudo apt-get install docker-compose-plugin"
    echo ""
    echo "See EC2_SETUP.md for detailed instructions."
    exit 1
else
    COMPOSE_VERSION=$(docker compose version)
    print_success "Docker Compose found: $COMPOSE_VERSION"
fi

# Check for Node.js
echo "Checking for Node.js..."
if ! command_exists node; then
    print_error "Node.js is not installed!"
    echo ""
    echo "Please install Node.js first:"
    echo "  curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -"
    echo "  sudo apt-get install -y nodejs"
    echo ""
    echo "See EC2_SETUP.md for detailed instructions."
    exit 1
else
    NODE_VERSION=$(node --version)
    print_success "Node.js found: $NODE_VERSION"
fi

# Check for npm
echo "Checking for npm..."
if ! command_exists npm; then
    print_error "npm is not installed!"
    echo ""
    echo "npm should come with Node.js. Please reinstall Node.js."
    echo "See EC2_SETUP.md for detailed instructions."
    exit 1
else
    NPM_VERSION=$(npm --version)
    print_success "npm found: v$NPM_VERSION"
fi

# Check for Go
echo "Checking for Go..."
if ! command_exists go; then
    print_error "Go is not installed!"
    echo ""
    echo "Please install Go first:"
    echo "  wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz"
    echo "  sudo rm -rf /usr/local/go"
    echo "  sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz"
    echo "  echo 'export PATH=\$PATH:/usr/local/go/bin' >> ~/.bashrc"
    echo "  source ~/.bashrc"
    echo ""
    echo "See EC2_SETUP.md for detailed instructions."
    exit 1
else
    GO_VERSION=$(go version)
    print_success "Go found: $GO_VERSION"
fi

# Check for .env file
echo "Checking for .env file..."
if [ ! -f "../.env" ]; then
    print_warning ".env file not found!"
    echo ""
    echo "Please create a .env file in the project root directory."
    echo "See EC2_SETUP.md for the required environment variables."
    echo ""
    read -p "Do you want to continue anyway? (y/N) " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
else
    print_success ".env file found"
fi

echo ""
echo "=============================="
echo "All dependencies satisfied! 🎉"
echo "=============================="
echo ""

# Navigate to project root
cd ..

# Start database container
echo "📦 Starting database container..."
docker compose -f docker-compose.yaml up -d db

# Wait for database to be ready
echo "⏳ Waiting for database to be ready..."
sleep 5

# Check if database is running
if docker ps | grep -q uplitycs-db; then
    print_success "Database container is running"
else
    print_error "Database container failed to start"
    echo "Check logs with: docker logs uplitycs-db-1"
    exit 1
fi

# Build frontend
echo ""
echo "🎨 Building frontend..."
cd frontend

# Install dependencies
echo "📥 Installing npm dependencies..."
npm install

# Build the frontend
echo "🔨 Building React app..."
npm run build

if [ $? -eq 0 ]; then
    print_success "Frontend build completed"
else
    print_error "Frontend build failed"
    exit 1
fi

cd ..

# Download Go dependencies
echo ""
echo "📥 Downloading Go dependencies..."
go mod download

if [ $? -eq 0 ]; then
    print_success "Go dependencies downloaded"
else
    print_error "Failed to download Go dependencies"
    exit 1
fi

# Start backend
echo ""
echo "🚀 Starting backend server..."
echo "=============================="
echo ""
print_info "Server will start on http://localhost:3333"
print_info "Press Ctrl+C to stop the server"
echo ""

go run .
