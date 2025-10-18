# UpLitycs - Status Page & Uptime Monitoring Platform

<div align="center">

[![Go](https://img.shields.io/badge/Go-1.24.1-00ADD8?logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-19.1-61DAFB?logo=react)](https://react.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Latest-336791?logo=postgresql)](https://www.postgresql.org)

**A modern, real-time uptime monitoring and status page platform with integrations for Slack, Discord, and Stripe payments**

[Features](#features) • [Quick Start](#quick-start) • [Installation](#installation) • [API Documentation](#api-documentation) • [Contributing](#contributing)

</div>

---

## Overview

UpLitycs is a comprehensive status page and uptime monitoring solution that helps businesses track the health of their applications and services in real-time. Built with a modern Go backend and React frontend, it provides beautiful, customizable status pages with deep integrations into popular communication platforms.

Monitor multiple applications, track SSL certificate expiration, send automated alerts to Slack and Discord, and manage your infrastructure with a single, intuitive dashboard.

---

## Features

### 🚀 Core Monitoring
- **Real-time Health Checks** - Continuous application health monitoring with configurable intervals
- **Multi-App Dashboard** - Monitor multiple applications from a single interface
- **Response Time Tracking** - Track application response times and performance metrics
- **24-Hour Uptime Statistics** - Historical uptime data and trend analysis
- **SSL Certificate Monitoring** - Automatic SSL expiry date tracking and renewal alerts
- **Status Code Tracking** - Detailed HTTP status code logging and analysis

### 🎨 Status Pages
- **Public Status Pages** - Share your application status with customers
- **Theme Customization** - Multiple theme options including cyberpunk retro style
- **Custom Domains** - Each status page gets a unique slug-based URL
- **Responsive Design** - Mobile-friendly status page display
- **Public API** - Access status data via public API endpoints
- **Uptime Badges** - Embeddable uptime badges for websites

### 💬 Integrations
- **Slack Integration** - Real-time incident notifications to Slack channels
  - Automatic downtime alerts
  - Recovery notifications
  - Status updates to configured channels
  
- **Discord Integration** - Discord webhook support for status updates
  - Server and channel configuration
  - Customizable notifications
  - Automatic incident tracking

### 💳 Subscription Management
- **Stripe Integration** - Secure payment processing
- **Multi-Tier Pricing** - Free, Pro, and Business plans with different features
- **Plan Limits** - Free tier: 1 app, Pro tier: 5 apps, Business tier: unlimited
- **Customer Portal** - Self-service subscription management
- **Usage Tracking** - Monitor plan usage and limits

### 📊 Pro/Business Features
- **Logo Upload** - Custom branding for status pages
- **Advanced Integrations** - Slack and Discord webhook support
- **Priority Support** - Dedicated support channels

### 🔐 Security & Authentication
- **OAuth 2.0 Integration** - Google OAuth authentication
- **Session Management** - Secure session handling with cookies
- **Admin Panel** - Administrative dashboard for site-wide monitoring
- **User Roles** - Admin and user permission levels

---

## Tech Stack

### Backend
- **Language**: Go 1.24.1
- **Framework**: Chi Router (chi/v5) - lightweight HTTP router
- **Database**: PostgreSQL
- **API Integrations**:
  - Stripe (payments)
  - AWS SES (email)
  - AWS S3 (logo storage)
  - Slack API
  - Discord Webhooks
  - Google OAuth 2.0

### Frontend
- **Framework**: React 19.1
- **Build Tool**: Vite 7.1
- **Routing**: React Router DOM 7.9
- **Styling**: Tailwind CSS 4.1 with custom animations
- **UI Components**: Lucide React icons
- **Animation**: Tailwind CSS Animate

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **Web Server**: Caddy (reverse proxy)
- **Database**: PostgreSQL
- **Email**: AWS Simple Email Service (SES)
- **Storage**: AWS S3

---

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.24.1 (for local development)
- Node.js 18+ (for frontend development)
- PostgreSQL (included in Docker)

### Using Docker (Recommended)

1. **Clone the repository**
```bash
git clone https://github.com/Ricardo-Ceia/UpLitycs.git
cd UpLitycs
```

2. **Configure environment variables**
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. **Start the application**
```bash
docker-compose up --build
```

4. **Access the application**
```
Frontend: http://localhost:3000 (development)
Backend API: http://localhost:8080
```

### Local Development Setup

#### Backend Setup
```bash
# Install Go dependencies
go mod download

# Set up PostgreSQL database
cd db
psql -U postgres -d statusframe < init/init.sql

# Run migrations
psql -U postgres -d statusframe < migrations/*.sql

# Start the backend server
go run main.go
```

#### Frontend Setup
```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

---

## Installation

### Prerequisites
- Docker & Docker Compose (recommended)
- PostgreSQL 13+
- AWS Account (for SES, S3)
- Stripe Account (for payments)
- OAuth Applications (Google, Discord, Slack)

### Environment Configuration

Create a `.env` file in the root directory:

```env
# Database
POSTGRES_DB=statusframe
POSTGRES_USER=postgres
POSTGRES_PASSWORD=your_secure_password

# Authentication
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
SESSION_SECRET=your_random_session_key
JWT_SECRET=your_jwt_secret

# Stripe
STRIPE_SECRET_KEY=your_stripe_secret_key
STRIPE_PUBLISHABLE_KEY=your_stripe_publishable_key
STRIPE_WEBHOOK_SECRET=your_stripe_webhook_secret
STRIPE_PRO_MONTHLY_PRICE_ID=price_xxx
STRIPE_PRO_YEARLY_PRICE_ID=price_xxx
STRIPE_BUSINESS_MONTHLY_PRICE=price_xxx
STRIPE_BUSINESS_YEARLY_PRICE=price_xxx

# AWS
AWS_REGION=eu-north-1
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_S3_BUCKET_NAME=your_s3_bucket

# Slack Integration
SLACK_CLIENT_ID=your_slack_client_id
SLACK_CLIENT_SECRET=your_slack_client_secret
SLACK_REDIRECT_URI=http://localhost:8080/api/slack/callback

# Discord Integration
DISCORD_CLIENT_ID=your_discord_client_id
DISCORD_CLIENT_SECRET=your_discord_client_secret
DISCORD_REDIRECT_URI=http://localhost:8080/api/discord/callback

# Application
APP_URL=http://localhost:8080
DOMAIN=yourdomain.com
ENVIRONMENT=development
```

### Running with Docker

```bash
# Build and run
docker-compose up --build

# Run in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for production deployment instructions.

---

## API Documentation

### Authentication
All protected endpoints require authentication via OAuth 2.0. Use the session cookie set by the auth endpoints.

### Public Endpoints

#### Get Public Status
```http
GET /api/public/status/{slug}
```
Returns the current status of a monitored application.

**Response:**
```json
{
  "app_name": "My Service",
  "status": "operational",
  "status_code": 200,
  "uptime_24h": 99.9,
  "last_checked": "2024-10-18T10:30:00Z"
}
```

#### Get Current Response Time
```http
GET /api/public/ping/{slug}
```
Returns the current response time for an application.

#### Get Uptime Badge
```http
GET /api/badge/{slug}
```
Returns an SVG uptime badge that can be embedded in websites.

---

### Protected Endpoints

#### Get User Status
```http
GET /api/user-status
Authorization: Bearer {token}
```
Returns the authenticated user's status and primary app slug.

#### Get User Apps
```http
GET /api/user-apps
Authorization: Bearer {token}
```
Returns all applications associated with the user.

**Response:**
```json
{
  "apps": [
    {
      "id": 1,
      "app_name": "My Service",
      "slug": "my-service",
      "health_url": "https://api.example.com/health",
      "theme": "cyberpunk",
      "status": "operational",
      "uptime_24h": 99.9,
      "ssl_days_until_expiry": 45
    }
  ]
}
```

#### Create New App (Onboarding)
```http
POST /api/go-to-dashboard
Authorization: Bearer {token}
Content-Type: application/json or multipart/form-data
```

**Request Body:**
```json
{
  "app_name": "My Service",
  "name": "Service Name",
  "homepage": "https://example.com",
  "slug": "my-service",
  "health_url": "https://api.example.com/health",
  "alerts": "y",
  "theme": "cyberpunk",
  "logo": "file" // optional, Pro/Business only
}
```

#### Delete App
```http
DELETE /api/apps/{appId}
Authorization: Bearer {token}
```

#### Get Plan Features
```http
GET /api/plan-features
Authorization: Bearer {token}
```

---

### Stripe Integration

#### Create Checkout Session
```http
POST /api/create-checkout-session
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "price_id": "price_xxx",
  "interval": "month"
}
```

#### Handle Webhook
```http
POST /api/stripe-webhook
Content-Type: application/json
```
Stripe sends webhook events here. Automatically handles:
- `checkout.session.completed`
- `customer.subscription.updated`
- `customer.subscription.deleted`

---

### Slack Integration

#### Start Slack Authentication
```http
GET /api/slack/start-auth
Authorization: Bearer {token}
```

#### Slack Callback (OAuth)
```http
GET /api/slack/callback?code={code}&state={state}
```

#### Save Slack Integration
```http
POST /api/slack/save-integration
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "channel_id": "C123456",
  "channel_name": "#alerts"
}
```

#### Get Slack Integration Status
```http
GET /api/slack/integration
Authorization: Bearer {token}
```

#### Disable Slack Integration
```http
POST /api/slack/disable
Authorization: Bearer {token}
```

---

### Discord Integration

#### Start Discord Authentication
```http
GET /api/discord/start-auth
Authorization: Bearer {token}
```

#### Discord Callback (OAuth)
```http
GET /api/discord/callback?code={code}&state={state}
```

#### Update Discord Webhook
```http
POST /api/discord/webhook
Authorization: Bearer {token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "webhook_url": "https://discordapp.com/api/webhooks/xxx/yyy"
}
```

#### Get Discord Integration Status
```http
GET /api/discord/integration
Authorization: Bearer {token}
```

#### Disable Discord Integration
```http
POST /api/discord/disable
Authorization: Bearer {token}
```

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  avatar_url TEXT,
  plan TEXT DEFAULT 'free',
  stripe_customer_id TEXT,
  stripe_subscription_id TEXT,
  created_at TIMESTAMPTZ DEFAULT now()
);
```

### Apps Table
```sql
CREATE TABLE apps (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id),
  app_name TEXT NOT NULL,
  slug TEXT UNIQUE NOT NULL,
  health_url TEXT NOT NULL,
  theme TEXT DEFAULT 'cyberpunk',
  alerts TEXT DEFAULT 'n',
  ssl_expiry_date TIMESTAMPTZ,
  ssl_days_until_expiry INTEGER,
  created_at TIMESTAMPTZ DEFAULT now()
);
```

### Status Tracking Table
```sql
CREATE TABLE user_status (
  id SERIAL PRIMARY KEY,
  app_id INTEGER REFERENCES apps(id),
  status_code INTEGER NOT NULL,
  checked_at TIMESTAMPTZ DEFAULT now()
);
```

### Integrations
- `slack_integrations` - Slack workspace configurations
- `discord_integrations` - Discord webhook configurations

Full schema available in `db/init/init.sql`

---

## Project Structure

```
UpLitycs/
├── backend/
│   ├── auth/              # OAuth and authentication
│   ├── handlers/          # HTTP request handlers
│   ├── email/             # Email service (AWS SES)
│   ├── stripe_config/     # Stripe integration
│   ├── utils/             # Utility functions
│   └── worker/            # Background workers
│       ├── health_checker.go
│       └── ssl_checker.go
├── frontend/
│   ├── src/
│   │   ├── components/    # React components
│   │   ├── pages/         # Page components
│   │   ├── styles/        # CSS files
│   │   └── App.jsx
│   ├── public/            # Static files
│   └── vite.config.js
├── db/
│   ├── init/              # Database initialization
│   └── migrations/        # Database migrations
├── main.go                # Application entry point
├── docker-compose.yaml
├── Dockerfile
└── README.md
```

---

## Configuration

### Health Check Configuration
Health checks run every 30 seconds. Configure the check interval in `main.go`:

```go
healthChecker := worker.NewHealthChecker(conn, 30*time.Second)
```

### SSL Certificate Checking
SSL certificates are checked daily. Configure in `main.go`:

```go
sslChecker := worker.NewSSLChecker(conn)
go sslChecker.Start()
```

### CORS Configuration
Configure allowed origins in `main.go`:

```go
AllowedOrigins: []string{
  "http://localhost:3000",
  "http://localhost:5173",
  "https://yourdomain.com",
}
```

---

## Pricing Tiers

### Free Plan
- ✅ 1 monitored application
- ✅ Real-time health checks
- ✅ 24-hour uptime tracking
- ✅ Public status page
- ❌ Logo upload
- ❌ Slack/Discord integrations

### Pro Plan
- ✅ 5 monitored applications
- ✅ All Free features
- ✅ Custom logo upload
- ✅ Slack integration
- ✅ Discord integration
- ✅ SSL certificate monitoring

### Business Plan
- ✅ Unlimited applications
- ✅ All Pro features
- ✅ Priority support
- ✅ Advanced analytics
- ✅ Dedicated account manager (optional)

---

## Development

### Prerequisites
- Go 1.24.1+
- Node.js 18+
- PostgreSQL 13+

### Running Tests

```bash
# Backend tests
go test ./... -v

# Frontend tests
cd frontend
npm run test
```

### Code Quality

```bash
# Backend linting
golangci-lint run

# Frontend linting
cd frontend
npm run lint
```

### Building

```bash
# Backend
go build -o bin/uplitycs

# Frontend
cd frontend
npm run build
```

---

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Contribution Guidelines

- Follow existing code style and conventions
- Write clear commit messages
- Include tests for new features
- Update documentation as needed
- Ensure all tests pass before submitting PR

### Reporting Issues

Found a bug? Please open an issue with:
- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Your environment details


## Performance & Scalability

### Optimization Techniques
- **Database Indexing** - Strategic indexes on frequently queried columns
- **Connection Pooling** - PostgreSQL connection pooling
- **Caching** - Response time caching for status pages
- **Batch Processing** - Batch health checks and SSL certificate checks

### Current Limits
- Health checks: 30-second intervals
- SSL checks: Daily checks
- App limit: Based on subscription tier

---

## Security Considerations

### Best Practices Implemented
- ✅ **OAuth 2.0** - Secure authentication with Google
- ✅ **CSRF Protection** - Secure session handling
- ✅ **SQL Injection Prevention** - Parameterized queries
- ✅ **CORS Configuration** - Restricted cross-origin requests
- ✅ **Environment Variables** - Sensitive data not in code
- ✅ **HTTPS Ready** - Full HTTPS support with Caddy
- ✅ **Admin Authentication** - Protected admin endpoints



### Database Connection Issues
```bash
# Check PostgreSQL is running
docker-compose ps

# Verify connection
psql -U postgres -h localhost -d statusframe -c "SELECT 1"
```

### Health Check Not Working
1. Verify health URL is accessible
2. Check network connectivity
3. Review backend logs: `docker-compose logs backend`
4. Check health checker configuration

### OAuth Issues
1. Verify client ID and secret in `.env`
2. Check redirect URIs match configuration
3. Ensure OAuth app is approved
4. Clear browser cookies and retry

### Stripe Payment Issues
1. Verify Stripe API keys in `.env`
2. Check webhook configuration
3. Review Stripe dashboard for errors
4. Ensure price IDs are correct

---

## Deployment

### Docker Deployment
See `docker-compose.yaml` for configuration.

### Cloud Platforms
- **AWS EC2** - Deploy with Docker
- **DigitalOcean** - Docker + App Platform
- **Heroku** - Buildpack deployment
- **Render** - Native support for Go + React

- Built with ❤️ by Ricardo Ceia
- Inspired by modern status page solutions
- Thanks to all contributors and community members


