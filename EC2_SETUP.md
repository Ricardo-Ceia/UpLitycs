# AWS EC2 Deployment Guide for UpLitycs

## Prerequisites on EC2 Instance

Your EC2 instance needs the following installed:
- Docker & Docker Compose
- Node.js (v18 or higher) and npm
- Go (v1.24.1 or compatible)
- Git

## Step-by-Step Setup Instructions

### 1. Connect to your EC2 instance
```bash
ssh -i your-key.pem ubuntu@your-ec2-public-ip
```

### 2. Install Docker and Docker Compose
```bash
# Update package index
sudo apt-get update

# Install Docker
sudo apt-get install -y ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Add your user to docker group (to run docker without sudo)
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
docker compose version
```

### 3. Install Node.js and npm
```bash
# Install Node.js 20.x (LTS)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify installation
node --version
npm --version
```

### 4. Install Go
```bash
# Download and install Go 1.24.1 (or latest compatible version)
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz
rm go1.24.1.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

### 5. Clone Your Repository
```bash
cd ~
git clone https://github.com/Ricardo-Ceia/UpLitycs.git
cd UpLitycs
```

### 6. Create Environment File
Create a `.env` file in the root of your project:
```bash
nano .env
```

Add your environment variables (example):
```env
# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=example
POSTGRES_DB=uplytics
DB_HOST=localhost
DB_PORT=5432

# OAuth Configuration
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://your-ec2-public-ip:3333/auth/google/callback

GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://your-ec2-public-ip:3333/auth/github/callback

# Session Secret
SESSION_KEY=your-random-session-secret-key-here

# Stripe Configuration
STRIPE_SECRET_KEY=your-stripe-secret-key
STRIPE_PUBLISHABLE_KEY=your-stripe-publishable-key
STRIPE_WEBHOOK_SECRET=your-stripe-webhook-secret
STRIPE_PRICE_ID_STARTER=your-stripe-price-id-starter
STRIPE_PRICE_ID_PRO=your-stripe-price-id-pro

# Application URL
APP_URL=http://your-ec2-public-ip:3333
```

Save with `Ctrl+X`, then `Y`, then `Enter`.

### 7. Configure Security Group
Make sure your EC2 Security Group allows inbound traffic on:
- Port 3333 (your application)
- Port 22 (SSH)
- Port 80/443 (if you plan to use a reverse proxy later)

### 8. Run the Application
```bash
cd ~/UpLitycs/scripts
chmod +x start.sh
./start.sh
```

## Production Considerations

### 1. Use a Process Manager (Recommended)
Install PM2 or use systemd to keep your app running:

```bash
# Install PM2
sudo npm install -g pm2

# Create a PM2 ecosystem file
cd ~/UpLitycs
```

Create `ecosystem.config.js`:
```javascript
module.exports = {
  apps: [{
    name: 'uplytics',
    script: './scripts/start.sh',
    interpreter: '/bin/bash',
    env: {
      NODE_ENV: 'production'
    }
  }]
};
```

Start with PM2:
```bash
pm2 start ecosystem.config.js
pm2 save
pm2 startup
```

### 2. Set Up Nginx Reverse Proxy (Recommended)
```bash
sudo apt-get install -y nginx

# Configure Nginx
sudo nano /etc/nginx/sites-available/uplytics
```

Add:
```nginx
server {
    listen 80;
    server_name your-domain.com;  # or your EC2 public IP

    location / {
        proxy_pass http://localhost:3333;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/uplytics /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 3. Set Up SSL with Let's Encrypt (Recommended for Production)
```bash
sudo apt-get install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### 4. Auto-start Docker on Boot
```bash
sudo systemctl enable docker
```

## Updating Your Application

To update your application after making changes:

```bash
cd ~/UpLitycs
git pull
cd scripts
./start.sh
```

## Troubleshooting

### Check Docker containers
```bash
docker ps
docker logs uplitycs-db-1
```

### Check application logs
```bash
# If running with PM2
pm2 logs uplytics

# If running directly
# Check terminal output where you ran start.sh
```

### Database connection issues
```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check database logs
docker logs uplitycs-db-1

# Test database connection
docker exec -it uplitycs-db-1 psql -U postgres -d uplytics
```

### Port already in use
```bash
# Find what's using port 3333
sudo lsof -i :3333
# Kill the process if needed
sudo kill -9 <PID>
```

## Monitoring

### Check system resources
```bash
# Memory usage
free -h

# Disk usage
df -h

# CPU usage
top
```

### Monitor Docker containers
```bash
docker stats
```
