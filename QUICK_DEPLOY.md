# UpLitycs - Quick EC2 Deployment

## 🚀 Fast Setup (For First-Time EC2 Setup)

### 1. Connect to EC2
```bash
ssh -i your-key.pem ubuntu@your-ec2-public-ip
```

### 2. Install All Dependencies (One Command)
```bash
curl -fsSL https://raw.githubusercontent.com/Ricardo-Ceia/UpLitycs/main/scripts/install_dependencies.sh | bash
```

Or manually:
```bash
git clone https://github.com/Ricardo-Ceia/UpLitycs.git
cd UpLitycs
chmod +x scripts/install_dependencies.sh
./scripts/install_dependencies.sh
```

### 3. Apply Docker Group Changes
```bash
newgrp docker
# Or logout and login again
```

### 4. Configure Environment
```bash
cd ~/UpLitycs
cp .env.example .env
nano .env  # Edit with your actual values
```

### 5. Configure EC2 Security Group
Allow inbound traffic on:
- **Port 3333** (Application) - Custom TCP
- **Port 22** (SSH) - SSH
- **Port 80** (Optional, for Nginx) - HTTP
- **Port 443** (Optional, for SSL) - HTTPS

### 6. Run the Application
```bash
cd ~/UpLitycs/scripts
chmod +x setup_ec2.sh
./setup_ec2.sh
```

---

## 🔄 Updating Your Application

```bash
cd ~/UpLitycs
git pull
cd scripts
./setup_ec2.sh
```

---

## ⚙️ Environment Variables You Need

Create `.env` file with:

```bash
# Required for Database
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your-secure-password
POSTGRES_DB=uplytics

# Required for OAuth
GOOGLE_CLIENT_ID=from-google-console
GOOGLE_CLIENT_SECRET=from-google-console
GOOGLE_REDIRECT_URL=http://YOUR-EC2-IP:3333/auth/google/callback

GITHUB_CLIENT_ID=from-github-settings
GITHUB_CLIENT_SECRET=from-github-settings  
GITHUB_REDIRECT_URL=http://YOUR-EC2-IP:3333/auth/github/callback

# Required for Sessions
SESSION_KEY=generate-random-32-char-key

# Required for Stripe
STRIPE_SECRET_KEY=from-stripe-dashboard
STRIPE_PUBLISHABLE_KEY=from-stripe-dashboard
STRIPE_WEBHOOK_SECRET=from-stripe-dashboard
STRIPE_PRICE_ID_STARTER=from-stripe-products
STRIPE_PRICE_ID_PRO=from-stripe-products

# App URL
APP_URL=http://YOUR-EC2-IP:3333
```

---

## 🔍 Troubleshooting

### Check if services are running:
```bash
docker ps                    # Database
curl http://localhost:3333   # Backend
```

### View logs:
```bash
docker logs uplitycs-db-1    # Database logs
# Backend logs appear in terminal
```

### Restart everything:
```bash
cd ~/UpLitycs
docker compose down
cd scripts
./setup_ec2.sh
```

### Port already in use:
```bash
sudo lsof -i :3333
sudo kill -9 <PID>
```

---

## 📝 Common Issues

### "npm: command not found" or "go: command not found"
- Run the install_dependencies.sh script
- Make sure to run `source ~/.bashrc` after Go installation
- Or logout and login again

### Database connection errors
- Check if database container is running: `docker ps`
- Check database logs: `docker logs uplitycs-db-1`
- Verify .env has correct database credentials

### OAuth not working
- Update redirect URLs in Google/GitHub OAuth settings
- Make sure URLs use your actual EC2 IP or domain
- Check if port 3333 is open in Security Group

### Stripe webhooks not working
- For webhooks to work, you need a public URL
- Consider using Stripe CLI for local testing
- Or set up ngrok/domain for production

---

## 🎯 Production Recommendations

### 1. Use a Domain Name
- Point a domain to your EC2 IP
- Update .env with your domain

### 2. Set Up SSL (HTTPS)
```bash
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d yourdomain.com
```

### 3. Use a Process Manager
Install PM2 to keep app running:
```bash
sudo npm install -g pm2
cd ~/UpLitycs
pm2 start scripts/setup_ec2.sh --name uplytics --interpreter bash
pm2 save
pm2 startup
```

### 4. Set Up Nginx Reverse Proxy
See EC2_SETUP.md for detailed instructions

### 5. Enable Automatic Backups
Back up your database regularly:
```bash
docker exec uplitycs-db-1 pg_dump -U postgres uplytics > backup.sql
```

---

## 📚 Additional Resources

- **Full Setup Guide**: See `EC2_SETUP.md`
- **Environment Template**: See `.env.example`
- **Stripe Docs**: https://stripe.com/docs
- **OAuth Setup**: Check your provider's documentation

---

## 🆘 Need Help?

1. Check the full EC2_SETUP.md guide
2. Review application logs
3. Check Docker container status
4. Verify all environment variables are set correctly
