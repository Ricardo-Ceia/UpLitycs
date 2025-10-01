# Quick Reference: Health Check Templates

This document provides a quick overview of all available health check templates in the UpLitycs onboarding process.

## Available Frameworks

### JavaScript/TypeScript
1. **Node.js + Express**
   - Key: `node-express`
   - Best for: Traditional Node.js backends
   - Port: 3000
   - Endpoint: `/health`

2. **Next.js API Routes**
   - Key: `node-nextjs`
   - Best for: Full-stack Next.js apps
   - Supports: Pages Router & App Router
   - Endpoint: `/api/health`

### Python
3. **Flask**
   - Key: `python-flask`
   - Best for: Simple Python APIs
   - Dependencies: `flask`, `psutil`
   - Port: 5000
   - Endpoint: `/health`

4. **FastAPI**
   - Key: `python-fastapi`
   - Best for: Modern async Python APIs
   - Dependencies: `fastapi`, `uvicorn`, `psutil`
   - Port: 8000
   - Endpoint: `/health`
   - Features: Auto-generated OpenAPI docs

### Go
5. **Gin Framework**
   - Key: `go-gin`
   - Best for: High-performance Go APIs
   - Package: `github.com/gin-gonic/gin`
   - Port: 8080
   - Endpoint: `/health`

6. **Go Standard Library**
   - Key: `go-standard`
   - Best for: Minimal dependencies
   - No external packages required
   - Port: 8080
   - Endpoint: `/health`

### Ruby
7. **Ruby on Rails**
   - Key: `ruby-rails`
   - Best for: Full-stack Rails apps
   - Controller: `HealthController`
   - Route: `get '/health', to: 'health#check'`
   - Endpoint: `/health`

8. **Sinatra**
   - Key: `ruby-sinatra`
   - Best for: Lightweight Ruby apps
   - Gem: `sinatra`
   - Port: 4567
   - Endpoint: `/health`

### PHP
9. **Laravel**
   - Key: `php-laravel`
   - Best for: Full-stack PHP apps
   - Controller: `HealthController`
   - Route: `routes/api.php`
   - Endpoint: `/api/health`

### .NET
10. **.NET / C#**
    - Key: `dotnet`
    - Best for: ASP.NET Core apps
    - Supports: Minimal API & Controllers
    - Port: 5001 (HTTPS)
    - Endpoint: `/health`

### Java
11. **Spring Boot**
    - Key: `java-spring`
    - Best for: Enterprise Java apps
    - Controller: `HealthController`
    - Port: 8080
    - Endpoint: `/health`
    - Alternative: Spring Boot Actuator

---

## Standard Response Format

All templates return a consistent JSON structure:

```json
{
  "status": "UP",
  "timestamp": "2025-10-01T12:00:00Z",
  "service": "my-service-name",
  "version": "1.0.0",
  "uptime": 3600,
  "memory": {
    "used_mb": 256
  },
  "environment": "production"
}
```

### Required Fields
- `status`: "UP" or "DOWN"

### Optional Fields (but included in templates)
- `timestamp`: ISO 8601 format
- `service`: Service identifier
- `version`: Semantic version
- `uptime`: Seconds since start
- `memory`: Memory usage stats
- `environment`: deployment environment

---

## Common Hosting Platforms by Framework

### Node.js
- Vercel (Next.js optimized)
- Heroku
- Railway
- Render
- Fly.io
- AWS Elastic Beanstalk

### Python
- Heroku
- Fly.io
- Railway
- Render
- Google Cloud Run
- AWS Elastic Beanstalk

### Go
- Fly.io (highly recommended)
- Railway
- Render
- Heroku
- Google Cloud Run
- Any VPS (easy deployment)

### Ruby
- Heroku (traditional choice)
- Fly.io
- Render
- Railway

### PHP
- Any shared hosting
- Heroku
- Laravel Forge
- Platform.sh
- DigitalOcean App Platform

### .NET
- Azure (native support)
- AWS Elastic Beanstalk
- Heroku
- Fly.io

### Java
- Heroku
- AWS Elastic Beanstalk
- Google Cloud Run
- Azure
- Any VPS (Tomcat/Spring Boot JAR)

---

## Testing Your Health Endpoint

### Using curl
```bash
curl https://your-domain.com/health
```

### Using browser
Simply navigate to:
```
https://your-domain.com/health
```

### Expected Response
You should see JSON output with at minimum:
```json
{"status": "UP"}
```

---

## Common Issues & Solutions

### Issue: 404 Not Found
**Solution:** 
- Verify the endpoint route is correct
- Check if the server is running
- Confirm you deployed the latest code

### Issue: 500 Internal Server Error
**Solution:**
- Check server logs for errors
- Verify all dependencies are installed
- Confirm environment variables are set

### Issue: CORS Error (in browser)
**Solution:**
- This is normal - our monitoring doesn't use browser
- If needed, add CORS headers to your endpoint

### Issue: Endpoint Behind Authentication
**Solution:**
- Health endpoints MUST be publicly accessible
- Remove authentication middleware from health route
- Use a separate endpoint for public monitoring

### Issue: Returns HTML Instead of JSON
**Solution:**
- Verify you're setting `Content-Type: application/json`
- Check route isn't being caught by a catch-all HTML handler
- Confirm the health route is registered correctly

---

## Best Practices

### ✅ DO
- Keep health checks lightweight (< 100ms response time)
- Return immediately without external dependencies
- Include useful metadata (uptime, version)
- Use standard HTTP status codes (200 for UP)
- Make endpoint publicly accessible

### ❌ DON'T
- Don't check database connections (too slow)
- Don't require authentication
- Don't perform heavy computations
- Don't make external API calls
- Don't include sensitive information

---

## Customization Tips

### Adding Custom Metrics
You can extend the response with your own fields:
```json
{
  "status": "UP",
  "custom_metric": "your_value",
  "active_connections": 42,
  "queue_size": 15
}
```

### Different Status Values
Besides "UP", you might use:
- "DOWN" - Service unavailable
- "DEGRADED" - Partial functionality
- "MAINTENANCE" - Planned downtime

### Status Codes
- 200: Service is UP
- 503: Service is DOWN
- 429: Rate limited (if applicable)

---

## Support

If you encounter issues with any template:
1. Check the setup instructions carefully
2. Verify dependencies are installed
3. Test the endpoint locally first
4. Check your hosting provider's logs
5. Ensure the endpoint is publicly accessible

For framework-specific help, refer to:
- Official framework documentation
- Hosting provider documentation
- Community forums and Stack Overflow
