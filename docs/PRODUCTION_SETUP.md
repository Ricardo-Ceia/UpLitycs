# StatusFrame Production Setup Guide

This guide will help you set up StatusFrame on a production server with Nginx, SSL, and proper DNS configuration.

## Prerequisites

- Ubuntu 22.04 LTS server
- Domain name pointing to your server's IP
- Root or sudo access
- At least 2GB RAM

## 1. Initial Server Setup

### Update system
```bash
sudo apt update && sudo apt upgrade -y
```

### Install required packages
```bash
sudo apt install -y nginx postgresql postgresql-contrib certbot python3-certbot-nginx git curl
```

### Install Go (1.21+)
```bash
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### Install Node.js (v20+)
```bash
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs
node --version
npm --version
```

## 2. DNS Configuration

Point your domain to your server's IP address:

### A Records (IPv4)
```
Type: A
Name: @ (or your domain)
Value: YOUR_SERVER_IP
TTL: 3600
```

```
Type: A
Name: www
Value: YOUR_SERVER_IP
TTL: 3600
```

### AAAA Records (IPv6 - optional)
```
Type: AAAA
Name: @
Value: YOUR_IPV6_ADDRESS
TTL: 3600
```

### Verify DNS propagation
```bash
# Check if DNS is pointing to your server
nslookup statusframe.com
dig statusframe.com
```

Wait for DNS propagation (can take up to 48 hours, usually minutes)

## 3. Clone and Configure Application

### Clone repository
```bash
cd /home/ubuntu
git clone https://github.com/YOUR_USERNAME/statusframe.git
cd statusframe
```

### Create environment file
```bash
cp .env.example .env
nano .env
```

Add your configuration:
```env
# Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=YOUR_SECURE_PASSWORD
POSTGRES_DB=statusframe
DATABASE_URL=postgresql://postgres:YOUR_SECURE_PASSWORD@localhost:5432/statusframe

# OAuth (Google)
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=https://statusframe.com/auth/google/callback

# Stripe
STRIPE_SECRET_KEY=your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=your_stripe_publishable_key
STRIPE_WEBHOOK_SECRET=your_stripe_webhook_secret

# Application
APP_URL=https://statusframe.com
SESSION_SECRET=generate_a_random_secure_string_here
ENVIRONMENT=production
```

### Set proper permissions
```bash
chmod 600 .env
```

## 4. Database Setup

### Configure PostgreSQL
```bash
sudo -u postgres psql
```

```sql
-- Create database and user
CREATE DATABASE statusframe;
CREATE USER statusframe_user WITH ENCRYPTED PASSWORD 'YOUR_SECURE_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE statusframe TO statusframe_user;

-- Connect to database
\c statusframe;

-- Grant schema permissions
GRANT ALL ON SCHEMA public TO statusframe_user;

\q
```

### Run migrations
```bash
# Migrations will run automatically when the app starts
# Or run manually with:
psql -U postgres -d statusframe -f db/init/001_schema.sql
psql -U postgres -d statusframe -f db/init/002_simplify_health_checks.sql
# ... etc
```

## 5. Build Application

### Build frontend
```bash
cd frontend
npm install
npm run build
cd ..
```

### Build backend
```bash
go mod download
go build -o statusframe main.go
chmod +x statusframe
```

### Test locally
```bash
./statusframe
# Should start on :8080
# Press Ctrl+C to stop
```

## 6. Setup Systemd Service

### Copy service file
```bash
sudo cp systemd/statusframe.service /etc/systemd/system/
sudo systemctl daemon-reload
```

### Start and enable service
```bash
sudo systemctl start statusframe
sudo systemctl enable statusframe
sudo systemctl status statusframe
```

### View logs
```bash
# Real-time logs
sudo journalctl -u statusframe -f

# Last 100 lines
sudo journalctl -u statusframe -n 100
```

## 7. Configure Nginx

### Copy Nginx configuration
```bash
sudo cp nginx/statusframe.conf /etc/nginx/sites-available/statusframe
```

### Update the configuration
```bash
sudo nano /etc/nginx/sites-available/statusframe
```

Replace `statusframe.com` with your actual domain.

### Enable site
```bash
sudo ln -s /etc/nginx/sites-available/statusframe /etc/nginx/sites-enabled/
sudo nginx -t
```

### Remove default Nginx site
```bash
sudo rm /etc/nginx/sites-enabled/default
```

## 8. SSL Certificate (Let's Encrypt)

### Create webroot directory
```bash
sudo mkdir -p /var/www/certbot
```

### Temporarily start Nginx (HTTP only for cert generation)
```bash
# Comment out SSL lines in nginx config first
sudo nano /etc/nginx/sites-available/statusframe
# Comment out lines 23-28 (ssl_certificate lines)
sudo nginx -t && sudo systemctl restart nginx
```

### Get SSL certificate
```bash
sudo certbot certonly --webroot \
  -w /var/www/certbot \
  -d statusframe.com \
  -d www.statusframe.com \
  --email your@email.com \
  --agree-tos \
  --no-eff-email
```

### Uncomment SSL lines in Nginx config
```bash
sudo nano /etc/nginx/sites-available/statusframe
# Uncomment lines 23-28
sudo nginx -t && sudo systemctl restart nginx
```

### Auto-renewal
Certbot automatically sets up renewal. Test it:
```bash
sudo certbot renew --dry-run
```

## 9. Firewall Configuration

### Setup UFW
```bash
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable
sudo ufw status
```

## 10. Final Checks

### Test endpoints
```bash
# Health check
curl https://statusframe.com/health

# API
curl https://statusframe.com/api/check-session
```

### Check services
```bash
# Application
sudo systemctl status statusframe

# Database
sudo systemctl status postgresql

# Nginx
sudo systemctl status nginx
```

### Monitor logs
```bash
# Application logs
sudo journalctl -u statusframe -f

# Nginx access logs
sudo tail -f /var/log/nginx/statusframe-access.log

# Nginx error logs
sudo tail -f /var/log/nginx/statusframe-error.log
```

## 11. Deployment Script

For future deployments, use the provided script:

```bash
chmod +x scripts/deploy.sh
sudo ./scripts/deploy.sh
```

## Troubleshooting

### Application won't start
```bash
# Check logs
sudo journalctl -u statusframe -n 100

# Check if port is in use
sudo lsof -i :8080

# Check environment file
cat .env
```

### Database connection errors
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Test database connection
psql -U postgres -d statusframe -c "SELECT 1;"

# Check connection string in .env
```

### Nginx errors
```bash
# Test configuration
sudo nginx -t

# Check error logs
sudo tail -f /var/log/nginx/error.log

# Restart Nginx
sudo systemctl restart nginx
```

### SSL issues
```bash
# Renew certificates
sudo certbot renew

# Check certificate expiry
sudo certbot certificates

# Test SSL configuration
openssl s_client -connect statusframe.com:443 -servername statusframe.com
```

## Security Best Practices

1. **Keep system updated**
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

2. **Regular backups**
   ```bash
   # Database backup
   pg_dump statusframe > backup_$(date +%Y%m%d).sql
   
   # Files backup
   tar -czf statusframe_backup_$(date +%Y%m%d).tar.gz /home/ubuntu/statusframe
   ```

3. **Monitor logs regularly**
   - Application logs
   - Nginx access/error logs
   - System logs

4. **Use strong passwords**
   - Database passwords
   - Session secrets
   - Server SSH keys

5. **Enable fail2ban** (optional)
   ```bash
   sudo apt install fail2ban
   sudo systemctl enable fail2ban
   ```

## Performance Optimization

### Database optimization
```sql
-- Add indexes if needed
CREATE INDEX idx_apps_user_id ON apps(user_id);
CREATE INDEX idx_user_status_app_id ON user_status(app_id);
```

### Nginx caching (optional)
Add to Nginx config:
```nginx
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=statusframe_cache:10m max_size=100m inactive=60m;
```

### Monitoring
Consider setting up monitoring tools:
- **Prometheus + Grafana** for metrics
- **Uptime Robot** for external monitoring
- **Sentry** for error tracking

## Support

For issues, check:
- GitHub Issues: https://github.com/YOUR_USERNAME/statusframe/issues
- Documentation: https://statusframe.com/docs
- Logs: `sudo journalctl -u statusframe -f`
