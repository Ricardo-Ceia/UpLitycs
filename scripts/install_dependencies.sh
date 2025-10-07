#!/bin/bash

# UpLitycs EC2 Dependency Installation Script
# Run this script on a fresh Ubuntu EC2 instance to install all required dependencies

set -e

echo "🔧 Installing UpLitycs Dependencies on EC2"
echo "==========================================="
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_info() {
    echo "ℹ️  $1"
}

# Check if running on Ubuntu
if [ ! -f /etc/os-release ]; then
    echo -e "${RED}❌ Could not detect OS. This script is designed for Ubuntu.${NC}"
    exit 1
fi

source /etc/os-release
if [[ "$ID" != "ubuntu" ]]; then
    echo -e "${RED}❌ This script is designed for Ubuntu. Detected: $ID${NC}"
    exit 1
fi

print_success "Running on Ubuntu $VERSION"
echo ""

# Update package index
print_info "Updating package index..."
sudo apt-get update -qq

# Install essential tools
print_info "Installing essential tools..."
sudo apt-get install -y -qq ca-certificates curl gnupg git wget unzip

print_success "Essential tools installed"
echo ""

# Install Docker
print_info "Installing Docker..."
if command -v docker >/dev/null 2>&1; then
    print_success "Docker already installed: $(docker --version)"
else
    # Add Docker's official GPG key
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc

    # Add Docker repository
    echo \
      "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
      $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
      sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

    # Install Docker
    sudo apt-get update -qq
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

    # Add current user to docker group
    sudo usermod -aG docker $USER

    print_success "Docker installed: $(docker --version)"
    print_success "Docker Compose installed: $(docker compose version)"
    echo ""
    echo -e "${YELLOW}⚠️  You may need to log out and log back in for Docker group changes to take effect${NC}"
    echo -e "${YELLOW}   Or run: newgrp docker${NC}"
    echo ""
fi

# Install Node.js
print_info "Installing Node.js..."
if command -v node >/dev/null 2>&1; then
    print_success "Node.js already installed: $(node --version)"
else
    # Install Node.js 20.x LTS
    curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
    sudo apt-get install -y nodejs

    print_success "Node.js installed: $(node --version)"
    print_success "npm installed: v$(npm --version)"
fi
echo ""

# Install Go
print_info "Installing Go..."
if command -v go >/dev/null 2>&1; then
    print_success "Go already installed: $(go version)"
else
    GO_VERSION="1.24.1"
    GO_ARCH="linux-amd64"
    
    # Download Go
    wget -q https://go.dev/dl/go${GO_VERSION}.${GO_ARCH}.tar.gz
    
    # Remove old installation and install new
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf go${GO_VERSION}.${GO_ARCH}.tar.gz
    rm go${GO_VERSION}.${GO_ARCH}.tar.gz
    
    # Add Go to PATH in .bashrc if not already there
    if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        echo 'export GOPATH=$HOME/go' >> ~/.bashrc
        echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
    fi
    
    # Add to current session
    export PATH=$PATH:/usr/local/go/bin
    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    print_success "Go installed: $(go version)"
    echo -e "${YELLOW}⚠️  Run 'source ~/.bashrc' to update your PATH in the current session${NC}"
fi
echo ""

# Summary
echo "==========================================="
echo -e "${GREEN}🎉 All dependencies installed successfully!${NC}"
echo "==========================================="
echo ""
echo "Installed versions:"
echo "  - Docker: $(docker --version)"
echo "  - Docker Compose: $(docker compose version)"
echo "  - Node.js: $(node --version)"
echo "  - npm: v$(npm --version)"
echo "  - Go: $(go version)"
echo ""
echo "Next steps:"
echo "  1. Clone your repository:"
echo "     git clone https://github.com/Ricardo-Ceia/UpLitycs.git"
echo "     cd UpLitycs"
echo ""
echo "  2. Create a .env file with your configuration"
echo "     (See EC2_SETUP.md for details)"
echo ""
echo "  3. Make setup script executable and run it:"
echo "     chmod +x scripts/setup_ec2.sh"
echo "     cd scripts"
echo "     ./setup_ec2.sh"
echo ""
echo "==========================================="
