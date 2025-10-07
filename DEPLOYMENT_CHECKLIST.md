# Quick Start Checklist - StatusFrame Production Deployment

## Pre-Deployment Checklist

### Domain & DNS
- [ ] Domain purchased and ready
- [ ] DNS A record pointing to server IP
- [ ] DNS propagated (check with `dig yourdomain.com`)

### Server Setup
- [ ] Ubuntu 22.04 LTS server provisioned
- [ ] SSH access configured
- [ ] Firewall rules planned
- [ ] At least 2GB RAM available

### Credentials & Keys
- [ ] Google OAuth credentials (production URLs)
- [ ] Stripe API keys (production)
- [ ] Strong database password generated
- [ ] Session secret generated
- [ ] Email address for SSL certificates

---

## Deployment Steps

### 1. Initial Server Setup (15 minutes)
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install basic dependencies
sudo apt install -y nginx postgresql postgresql-contrib \
  certbot python3-certbot-nginx git curl ufw

# Install Go
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Verify installations
go version
node --version
npm --version
```

### 2. Clone & Configure (10 minutes)
```bash
# Clone repository
cd ~
git clone https://github.com/Ricardo-Ceia/statusframe.git
cd statusframe

# Create environment file
nano .env
```

**Required .env variables:**
```env
# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=YOUR_STRONG_PASSWORD_HERE
POSTGRES_DB=statusframe
DATABASE_URL=postgresql://postgres:YOUR_STRONG_PASSWORD_HERE@localhost:5432/statusframe

# OAuth
GOOGLE_CLIENT_ID=your_google_client_id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=https://yourdomain.com/auth/google/callback

# Stripe
STRIPE_SECRET_KEY=sk_live_your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=pk_live_your_stripe_publishable_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# Application
APP_URL=https://yourdomain.com
SESSION_SECRET=generate_random_64_char_string_here
ENVIRONMENT=production
PORT=8080
```

```bash
# Secure the file
chmod 600 .env
```

### 3. Database Setup (5 minutes)
```bash
# Access PostgreSQL
sudo -u postgres psql

# In PostgreSQL prompt:
CREATE DATABASE statusframe;
CREATE USER statusframe_user WITH ENCRYPTED PASSWORD 'YOUR_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE statusframe TO statusframe_user;
\c statusframe
GRANT ALL ON SCHEMA public TO statusframe_user;
\q
```

### 4. Build Application (10 minutes)
```bash
# Build frontend
cd ~/statusframe/frontend
npm install
npm run build

# Verify build
ls -la dist/

# Build backend
cd ~/statusframe
go mod download
go build -o statusframe main.go
chmod +x statusframe

# Quick test (optional)
./statusframe &
sleep 3
curl http://localhost:8080
pkill statusframe
```

### 5. Setup Systemd Service (5 minutes)
```bash
# Copy and install service
sudo cp ~/statusframe/systemd/statusframe.service /etc/systemd/system/

# Edit if needed (paths, user)
sudo nano /etc/systemd/system/statusframe.service

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable statusframe
sudo systemctl start statusframe

# Check status
sudo systemctl status statusframe

# View logs
sudo journalctl -u statusframe -n 50
```

### 6. Configure Nginx (10 minutes)
```bash
# Copy config
sudo cp ~/statusframe/nginx/statusframe.conf /etc/nginx/sites-available/statusframe

# Update domain name in config
sudo nano /etc/nginx/sites-available/statusframe
# Replace all instances of "statusframe.com" with your domain

# Initially comment out SSL lines (lines 23-28)
# We'll uncomment after getting certificate

# Enable site
sudo ln -s /etc/nginx/sites-available/statusframe /etc/nginx/sites-enabled/
sudo rm /etc/nginx/sites-enabled/default

# Test and start
sudo nginx -t
sudo systemctl enable nginx
sudo systemctl restart nginx
```

### 7. Get SSL Certificate (10 minutes)
```bash
# Create webroot
sudo mkdir -p /var/www/certbot

# Get certificate
sudo certbot certonly --webroot \
  -w /var/www/certbot \
  -d yourdomain.com \
  -d www.yourdomain.com \
  --email your@email.com \
  --agree-tos \
  --no-eff-email

# Uncomment SSL lines in Nginx config
sudo nano /etc/nginx/sites-available/statusframe
# Uncomment lines 23-28

# Test and reload
sudo nginx -t
sudo systemctl reload nginx
```

### 8. Configure Firewall (2 minutes)
```bash
# Allow necessary ports
sudo ufw allow 22/tcp   # SSH
sudo ufw allow 80/tcp   # HTTP
sudo ufw allow 443/tcp  # HTTPS

# Enable firewall
sudo ufw --force enable

# Check status
sudo ufw status
```

### 9. Verification (5 minutes)
```bash
# Test all services
sudo systemctl status statusframe
sudo systemctl status nginx
sudo systemctl status postgresql

# Test endpoints
curl https://yourdomain.com
curl https://yourdomain.com/api/check-session

# Check SSL
curl -I https://yourdomain.com

# Monitor logs
sudo journalctl -u statusframe -f
```

---

## Post-Deployment Tasks

### Immediate
- [ ] Test user registration flow
- [ ] Test app creation
- [ ] Test status page viewing
- [ ] Verify email alerts work
- [ ] Test Stripe integration (sandbox first)

### Within 24 Hours
- [ ] Set up monitoring (Uptime Robot, Pingdom, etc.)
- [ ] Configure database backups
- [ ] Set up log rotation
- [ ] Configure Google Analytics (optional)
- [ ] Test SSL certificate auto-renewal: `sudo certbot renew --dry-run`

### Within Week
- [ ] Set up automated backups (cron job)
- [ ] Configure error tracking (Sentry, etc.)
- [ ] Set up metrics/monitoring (Prometheus + Grafana)
- [ ] Performance testing
- [ ] Security audit
- [ ] Document recovery procedures

---

## Common Issues & Solutions

### Service Won't Start
```bash
# Check logs
sudo journalctl -u statusframe -n 100

# Check if port in use
sudo lsof -i :8080

# Verify .env file
cat .env

# Check file permissions
ls -la statusframe
```

### Database Connection Errors
```bash
# Verify PostgreSQL is running
sudo systemctl status postgresql

# Test connection
psql -U postgres -d statusframe -c "SELECT 1;"

# Check connection string
grep DATABASE_URL .env
```

### Nginx Won't Start
```bash
# Test config
sudo nginx -t

# Check logs
sudo tail -f /var/log/nginx/error.log

# Verify syntax
sudo nginx -c /etc/nginx/nginx.conf -t
```

### SSL Certificate Issues
```bash
# Check certificate
sudo certbot certificates

# Manual renewal
sudo certbot renew

# Check certificate files
sudo ls -la /etc/letsencrypt/live/yourdomain.com/
```

---

## Maintenance Commands

### View Logs
```bash
# Application
sudo journalctl -u statusframe -f

# Nginx access
sudo tail -f /var/log/nginx/statusframe-access.log

# Nginx errors
sudo tail -f /var/log/nginx/statusframe-error.log
```

### Restart Services
```bash
# Application
sudo systemctl restart statusframe

# Nginx
sudo systemctl reload nginx

# PostgreSQL
sudo systemctl restart postgresql
```

### Update Application
```bash
cd ~/statusframe
git pull
sudo ./scripts/deploy.sh
```

### Backup Database
```bash
# Manual backup
pg_dump statusframe > ~/backups/statusframe_$(date +%Y%m%d).sql

# Restore backup
psql statusframe < ~/backups/statusframe_20250107.sql
```

---

## Automated Deployment

For future updates, use the deployment script:

```bash
cd ~/statusframe
git pull
sudo ./scripts/deploy.sh
```

This script:
1. Stops the service
2. Builds frontend
3. Builds backend
4. Starts service
5. Reloads Nginx
6. Verifies everything is running

---

## Support & Documentation

- **Full Setup Guide**: `docs/PRODUCTION_SETUP.md`
- **DNS Configuration**: `docs/DNS_SETUP.md`
- **Rebranding Summary**: `REBRANDING_SUMMARY.md`
- **Plan Limits**: `PLAN_LIMIT_IMPLEMENTATION.md`

---

**Total Deployment Time**: ~1-2 hours (including DNS propagation)

**Skills Needed**: Basic Linux, command line comfort

**Cost Estimate** (monthly):
- Server (2GB RAM): $10-20
- Domain: $10-15/year
- SSL Certificate: Free (Let's Encrypt)

---

*Ready to deploy? Start with Step 1!* 🚀
