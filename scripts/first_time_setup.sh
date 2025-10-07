#!/bin/bash

# UpLitycs - Complete EC2 First-Time Setup
# This script guides you through the entire setup process

set -e

clear

cat << "EOF"
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║              UpLitycs EC2 First-Time Setup                   ║
║                                                              ║
║  This script will guide you through setting up UpLitycs     ║
║  on your EC2 instance for the first time.                   ║
║                                                              ║
╚══════════════════════════════════════════════════════════════╝
EOF

echo ""
echo "Press Enter to continue..."
read

# Step 1: Install Dependencies
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 1: Installing Dependencies"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "This will install:"
echo "  - Docker & Docker Compose"
echo "  - Node.js & npm"
echo "  - Go"
echo ""
read -p "Continue with installation? (y/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Setup cancelled."
    exit 1
fi

# Run dependency installation
if [ -f "install_dependencies.sh" ]; then
    ./install_dependencies.sh
else
    echo "Error: install_dependencies.sh not found!"
    exit 1
fi

# Step 2: Apply Docker group changes
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 2: Applying Docker Group Changes"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "To use Docker without sudo, we need to apply group changes."
echo "This requires starting a new shell session."
echo ""
newgrp docker << EOFNEWGRP

# Step 3: Environment Configuration
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 3: Environment Configuration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

cd ..

if [ -f ".env" ]; then
    echo "⚠️  .env file already exists."
    read -p "Do you want to edit it? (y/N) " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        nano .env
    fi
else
    echo "Creating .env file from template..."
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo ""
        echo "✅ .env file created!"
        echo ""
        echo "Please edit the .env file with your actual configuration:"
        echo "  - Database credentials"
        echo "  - OAuth client IDs and secrets (Google, GitHub)"
        echo "  - Stripe API keys"
        echo "  - Session secret key"
        echo "  - Your EC2 public IP or domain"
        echo ""
        read -p "Press Enter to open the editor..."
        nano .env
    else
        echo "❌ Error: .env.example not found!"
        exit 1
    fi
fi

# Step 4: Security Group Check
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 4: EC2 Security Group Configuration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⚠️  IMPORTANT: Make sure your EC2 Security Group allows:"
echo ""
echo "  Inbound Rules:"
echo "    - Port 22 (SSH)        - To connect to your instance"
echo "    - Port 3333 (Custom)   - For the application"
echo "    - Port 80 (HTTP)       - Optional, for Nginx"
echo "    - Port 443 (HTTPS)     - Optional, for SSL"
echo ""
echo "To configure:"
echo "  1. Go to AWS Console → EC2 → Security Groups"
echo "  2. Find your instance's security group"
echo "  3. Edit Inbound Rules"
echo "  4. Add the rules above"
echo ""
read -p "Have you configured the Security Group? (y/N) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo ""
    echo "Please configure your Security Group before continuing."
    echo "Press Enter when done..."
    read
fi

# Step 5: Run the application
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 5: Starting UpLitycs"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Ready to start UpLitycs!"
echo ""
read -p "Start the application now? (y/N) " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    cd scripts
    ./setup_ec2.sh
else
    echo ""
    echo "Setup complete! To start the application later, run:"
    echo "  cd ~/UpLitycs/scripts"
    echo "  ./setup_ec2.sh"
fi

EOFNEWGRP
