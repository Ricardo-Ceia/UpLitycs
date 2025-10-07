# UpLitycs 📊

A modern uptime monitoring SaaS platform with real-time status tracking, beautiful dashboards, and multi-tier subscription plans.

## Features

- 🔄 **Real-time Uptime Monitoring** - Track your websites and APIs 24/7
- 📊 **Beautiful Dashboards** - Visualize uptime stats with interactive charts
- 🎨 **Customizable Status Pages** - Public status pages for your services
- 👥 **Multi-App Support** - Monitor multiple applications from one account
- 💳 **Stripe Integration** - Flexible pricing tiers (Free, Starter, Pro)
- 🔐 **OAuth Authentication** - Sign in with Google or GitHub
- ⚡ **Real-time Health Checks** - Automatic health monitoring every 30 seconds
- 🎯 **Response Time Tracking** - Monitor API performance

## Tech Stack

**Backend:**
- Go 1.24.1
- Chi Router
- PostgreSQL 17
- Stripe API

**Frontend:**
- React 19
- Vite
- TailwindCSS
- React Router

**Infrastructure:**
- Docker & Docker Compose
- AWS EC2 Ready

## Quick Start (Local Development)

### Prerequisites
- Go 1.24.1+
- Node.js 18+
- Docker & Docker Compose

### 1. Clone the Repository
```bash
git clone https://github.com/Ricardo-Ceia/UpLitycs.git
cd UpLitycs
```

### 2. Set Up Environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

### 3. Run the Application
```bash
cd scripts
chmod +x start.sh
./start.sh
```

The application will be available at `http://localhost:3333`

## 🚀 AWS EC2 Deployment

### Quick Deploy
See **[QUICK_DEPLOY.md](./QUICK_DEPLOY.md)** for a fast deployment guide.

### Detailed Setup
See **[EC2_SETUP.md](./EC2_SETUP.md)** for comprehensive EC2 setup instructions.

### First-Time EC2 Setup
```bash
# On your EC2 instance
curl -fsSL https://raw.githubusercontent.com/Ricardo-Ceia/UpLitycs/main/scripts/install_dependencies.sh | bash
```

Or:
```bash
git clone https://github.com/Ricardo-Ceia/UpLitycs.git
cd UpLitycs
chmod +x scripts/install_dependencies.sh
./scripts/install_dependencies.sh
```

Then:
```bash
cp .env.example .env
nano .env  # Configure your environment
cd scripts
chmod +x setup_ec2.sh
./setup_ec2.sh
```

## 📝 Environment Configuration

Required environment variables (see `.env.example` for full list):

- **Database**: PostgreSQL credentials
- **OAuth**: Google and GitHub client IDs and secrets
- **Stripe**: API keys and price IDs
- **Session**: Secret key for session management
- **App URL**: Your application's public URL

## 🗂️ Project Structure

```
UpLitycs/
├── backend/              # Go backend code
│   ├── auth/            # Authentication logic
│   ├── handlers/        # HTTP handlers
│   ├── stripe_config/   # Stripe configuration
│   ├── worker/          # Background health checker
│   └── utils/           # Utility functions
├── db/                  # Database migrations and setup
├── frontend/            # React frontend
│   └── src/            
│       ├── Dashboard.jsx    # Main dashboard
│       ├── Pricing.jsx      # Pricing page
│       ├── StatusPage.jsx   # Public status pages
│       └── ...
├── scripts/             # Deployment and setup scripts
│   ├── start.sh                # Local dev startup
│   ├── setup_ec2.sh           # EC2 setup with checks
│   └── install_dependencies.sh # Install all dependencies
├── docker-compose.yaml  # Docker services
├── main.go             # Application entry point
└── .env.example        # Environment template
```

## 🔧 Development Scripts

### Local Development
```bash
cd scripts
./start.sh
```

### EC2 Deployment
```bash
cd scripts
./setup_ec2.sh
```

### Install Dependencies (EC2)
```bash
./scripts/install_dependencies.sh
```

## 🧪 Testing

```bash
# Run Go tests
go test ./tests/...

# Run specific test
go test -run TestHTTP ./tests/
```

## 📊 Database

PostgreSQL database with automatic migrations on startup:
- User management
- App tracking
- Health check history
- Subscription plans

Migrations are in `db/init/` and run automatically when the database container starts.

## 🔐 Authentication

Supports OAuth 2.0 with:
- Google
- GitHub

Session-based authentication with secure cookies.

## 💳 Pricing Tiers

- **Free**: 1 app, basic monitoring
- **Starter**: 5 apps, $9.99/month
- **Pro**: Unlimited apps, $29.99/month

Managed through Stripe subscriptions.

## 🚀 Deployment Checklist

- [ ] Install dependencies (Node.js, Go, Docker)
- [ ] Configure `.env` file
- [ ] Set up OAuth apps (Google, GitHub)
- [ ] Configure Stripe products and prices
- [ ] Update EC2 Security Group (allow port 3333)
- [ ] Run setup script
- [ ] (Optional) Set up Nginx reverse proxy
- [ ] (Optional) Configure SSL with Let's Encrypt
- [ ] (Optional) Set up PM2 for process management

## 🐛 Troubleshooting

### Common Issues

**"npm: command not found" or "go: command not found"**
- Run `install_dependencies.sh`
- Source your bashrc: `source ~/.bashrc`

**Database connection errors**
- Check if container is running: `docker ps`
- View logs: `docker logs uplitycs-db-1`

**OAuth redirect issues**
- Verify redirect URLs match your domain/IP
- Check Security Group allows port 3333

See [QUICK_DEPLOY.md](./QUICK_DEPLOY.md) for more troubleshooting tips.

## 📚 Documentation

- **[QUICK_DEPLOY.md](./QUICK_DEPLOY.md)** - Fast deployment reference
- **[EC2_SETUP.md](./EC2_SETUP.md)** - Detailed EC2 setup guide
- **[.env.example](./.env.example)** - Environment configuration template

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## 📄 License

[Add your license here]

## 👤 Author

Ricardo Ceia

## 🔗 Links

- **GitHub**: https://github.com/Ricardo-Ceia/UpLitycs
- **Live Demo**: [Add your deployed URL]

---

Made with ❤️ using Go and React
