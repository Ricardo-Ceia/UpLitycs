# Project Structure

Complete overview of the UpLitycs project directory structure and file purposes.

## Directory Tree

```
UpLitycs/
├── .github/                          # GitHub configuration
│   └── workflows/                    # CI/CD workflows
│
├── app/                              # Application builds
│
├── backend/                          # Go backend application
│   ├── auth/                         # Authentication & OAuth
│   │   └── auth.go                   # OAuth handlers and middleware
│   │
│   ├── email/                        # Email service
│   │   └── ses.go                    # AWS SES integration
│   │
│   ├── handlers/                     # HTTP request handlers
│   │   ├── handlers.go               # Core CRUD handlers
│   │   ├── admin_handlers.go         # Admin panel handlers
│   │   ├── discord_handlers.go       # Discord integration
│   │   ├── slack_handlers.go         # Slack integration
│   │   └── stripe_handlers.go        # Stripe payment handlers
│   │
│   ├── stripe_config/                # Stripe configuration
│   │   └── config.go                 # Stripe client setup
│   │
│   ├── utils/                        # Utility functions
│   │   └── utils.go                  # Helper functions
│   │
│   └── worker/                       # Background workers
│       ├── health_checker.go         # App health monitoring
│       └── ssl_checker.go            # SSL certificate monitoring
│
├── db/                               # Database configuration
│   ├── init/                         # Database initialization
│   │   └── init.sql                  # SQL schema creation
│   │
│   ├── migrations/                   # Database migrations
│   │   ├── add_discord_integration.sql
│   │   ├── add_logo_url.sql
│   │   ├── add_slack_integration.sql
│   │   └── add_ssl_tracking.sql
│   │
│   └── db.go                         # Database connection & queries
│
├── frontend/                         # React frontend application
│   ├── public/                       # Static files (not processed)
│   │
│   ├── src/                          # Source code
│   │   ├── components/               # Reusable components
│   │   ├── pages/                    # Page components
│   │   │
│   │   ├── Admin.jsx                 # Admin panel page
│   │   ├── Admin.css                 # Admin styles
│   │   ├── App.jsx                   # Root application component
│   │   ├── Dashboard.jsx             # User dashboard
│   │   ├── Dashboard.css             # Dashboard styles
│   │   ├── DiscordIntegration.jsx    # Discord settings
│   │   ├── DiscordIntegration.css    # Discord styles
│   │   ├── Home.jsx                  # Landing page
│   │   ├── Home.css                  # Home page styles
│   │   ├── index.css                 # Global styles
│   │   ├── main.jsx                  # Application entry point
│   │   ├── Onboard.jsx               # Onboarding flow
│   │   ├── PlanFeatures.jsx          # Pricing features display
│   │   ├── PlanFeatures.css          # Plan features styles
│   │   ├── Pricing.jsx               # Pricing page
│   │   ├── Pricing.css               # Pricing styles
│   │   ├── ProtectedRoute.jsx        # Authentication wrapper
│   │   ├── RetroAuth.jsx             # Login page
│   │   ├── RetroAuth.css             # Login styles
│   │   ├── Settings.jsx              # User settings page
│   │   ├── Settings.css              # Settings styles
│   │   ├── SlackIntegration.jsx      # Slack settings
│   │   ├── SlackIntegration.css      # Slack styles
│   │   ├── StatusPage.jsx            # Public status page
│   │   ├── StatusPage.css            # Status page styles
│   │   ├── UpgradeModal.jsx          # Plan upgrade modal
│   │   ├── UpgradeModal.css          # Modal styles
│   │   ├── UptimeBarGraph.jsx        # Uptime chart
│   │   ├── UptimeBarGraph.css        # Chart styles
│   │   │
│   │   └── assets/                   # Images, logos, etc.
│   │
│   ├── package.json                  # Frontend dependencies
│   ├── package-lock.json             # Dependency lock file
│   ├── vite.config.js                # Vite build configuration
│   ├── tailwind.config.js            # Tailwind CSS config
│   ├── eslint.config.js              # ESLint configuration
│   ├── jsconfig.json                 # JavaScript configuration
│   ├── index.html                    # HTML entry point
│   └── README.md                     # Frontend-specific docs
│
├── tests/                            # Test files
│   └── health_checker_test.go        # Health checker tests
│
├── .env                              # Environment variables (local)
├── .env.example                      # Environment variables template
├── .gitignore                        # Git ignore rules
├── .dockerignore                     # Docker ignore rules
│
├── main.go                           # Backend entry point
├── go.mod                            # Go module definition
├── go.sum                            # Go dependency hash lock
│
├── Caddyfile                         # Caddy web server config (local)
├── Caddyfile.global                  # Caddy global config (production)
├── Dockerfile                        # Docker image definition
├── docker-compose.yaml               # Multi-container setup
│
├── start.sh                          # Start script
├── stop.sh                           # Stop script
├── deploy-slack.sh                   # Deployment notification script
│
├── ssh.pem                           # SSH key (local dev only)
│
├── README.md                         # Main project documentation
```
