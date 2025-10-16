# SSL/TLS Certificate Expiry Tracking Feature

## Overview
Implemented SSL/TLS certificate expiry tracking for HTTPS-enabled apps with daily monitoring and visual dashboard display.

## Features

### ✅ What Was Implemented

1. **Daily SSL Certificate Checker**
   - Automatically checks all HTTPS apps every 24 hours
   - Extracts certificate expiry date, days until expiry, and issuer
   - Stores SSL data in database for quick access

2. **Database Schema**
   - Added 4 new columns to `apps` table:
     - `ssl_expiry_date` - Timestamp of when cert expires
     - `ssl_days_until_expiry` - Number of days until expiry
     - `ssl_issuer` - Certificate issuer organization
     - `ssl_last_checked` - Last time SSL was checked

3. **Plan-Based Access Control**
   - **Free Tier**: Shows locked SSL section with upgrade prompt
   - **Pro & Business**: Full SSL tracking with visual indicators

4. **Enhanced Dashboard UI**
   - Visual status badges with color coding:
     - 🔴 **CRITICAL** (≤7 days or expired) - Red with pulsing animation
     - ⚡ **SOON** (≤30 days) - Yellow with pulse
     - ✓ **VALID** (>30 days) - Green
     - 🔒 **LOCKED** (Free users) - Upgrade prompt
   
   - Displays:
     - Days until expiry
     - Certificate issuer
     - Status badge
     - Helpful tooltips

## File Changes

### Backend Files
1. **`/backend/worker/ssl_checker.go`** (NEW)
   - SSL certificate checking worker
   - Runs daily at startup and every 24 hours
   - Connects via TLS to extract certificate data

2. **`/db/db.go`**
   - Added `UpdateSSLInfo()` function
   - Modified `GetUserAppsWithStatus()` to include SSL data (Pro/Business only)
   - Added `time` import

3. **`/db/migrations/add_ssl_tracking.sql`** (NEW)
   - Migration to add SSL columns to apps table

4. **`/main.go`**
   - Initialize and start SSL checker worker

### Frontend Files
1. **`/frontend/src/Dashboard.jsx`**
   - Added `Shield` and `ShieldAlert` icons
   - Enhanced `getSSLStatus()` function with tooltips and badges
   - Plan-based rendering (locked for free users)

2. **`/frontend/src/Dashboard.css`**
   - Extensive SSL styling with animations
   - Color-coded status indicators
   - Locked state styling for free users
   - Pulsing animations for critical/warning states

## How It Works

### Worker Flow
```
1. App starts → SSL checker initializes
2. Immediately checks all HTTPS apps
3. For each app:
   - Parse URL and connect via TLS
   - Extract certificate from connection
   - Calculate days until expiry
   - Save to database
4. Sleep for 24 hours
5. Repeat from step 2
```

### Database Query
- Free users: SSL fields return NULL
- Pro/Business: SSL fields populated with actual data

### UI Display Logic
```
HTTPS App?
  ├─ No → Don't show SSL section
  └─ Yes → Check user plan
      ├─ Free → Show locked upgrade prompt
      └─ Pro/Business → Show SSL status
          ├─ No data → "Checking..."
          ├─ Expired → Red badge "EXPIRED"
          ├─ ≤7 days → Red badge "URGENT"
          ├─ ≤30 days → Yellow badge "SOON"
          └─ >30 days → Green badge "VALID"
```

## Migration Instructions

### If Database Already Exists:
```bash
docker compose up -d db
docker compose exec db psql -U postgres -d statusframe -c "
  ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_expiry_date TIMESTAMPTZ;
  ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_days_until_expiry INTEGER;
  ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_issuer TEXT;
  ALTER TABLE apps ADD COLUMN IF NOT EXISTS ssl_last_checked TIMESTAMPTZ;
"
```

### For Fresh Installs:
The columns will be created automatically from the migration file on first run.

## Testing

1. **Add HTTPS App**: Create an app with an HTTPS URL
2. **Wait for Check**: SSL checker runs on startup (immediate)
3. **View Dashboard**: 
   - Pro/Business users see SSL status with days remaining
   - Free users see upgrade prompt
4. **Check Logs**: `docker compose logs app | grep SSL`

## Future Enhancements (Optional)

- [ ] Email alerts when SSL expires in X days
- [ ] SSL history/timeline view
- [ ] Manual SSL check trigger
- [ ] Support for custom SSL check intervals
- [ ] SSL strength/grade analysis

## Effort

**Size: S (Small)** ✅
- Estimated: 2-3 hours
- Actual: ~2 hours
- Clean implementation with minimal dependencies
