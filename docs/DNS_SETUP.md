# DNS Setup Quick Reference

## Required DNS Records for StatusFrame

Replace `YOUR_SERVER_IP` with your actual server's public IP address.

### For Custom Domain (e.g., statusframe.com)

#### Basic Setup
```
Type: A
Name: @
Value: YOUR_SERVER_IP
TTL: 3600

Type: A
Name: www
Value: YOUR_SERVER_IP
TTL: 3600
```

#### With IPv6 (Optional)
```
Type: AAAA
Name: @
Value: YOUR_IPV6_ADDRESS
TTL: 3600

Type: AAAA
Name: www
Value: YOUR_IPV6_ADDRESS
TTL: 3600
```

## Popular DNS Providers

### Cloudflare
1. Log in to Cloudflare dashboard
2. Select your domain
3. Go to DNS tab
4. Click "Add record"
5. Add the A records above
6. **Important**: Set proxy status to "DNS only" (gray cloud) initially
7. After SSL is setup, you can enable proxy (orange cloud)

### Namecheap
1. Log in to Namecheap
2. Go to Domain List
3. Click "Manage" next to your domain
4. Go to "Advanced DNS" tab
5. Add A records as shown above

### GoDaddy
1. Log in to GoDaddy
2. Go to My Products
3. Click "DNS" next to your domain
4. Add A records as shown above

### AWS Route 53
1. Open Route 53 console
2. Select your hosted zone
3. Click "Create record"
4. Add A records as shown above

### Google Domains
1. Log in to Google Domains
2. Click your domain
3. Click "DNS" in left menu
4. Scroll to "Custom resource records"
5. Add A records as shown above

## Verification

### Check DNS Propagation
```bash
# Check A record
dig statusframe.com
nslookup statusframe.com

# Check from different locations
dig @8.8.8.8 statusframe.com        # Google DNS
dig @1.1.1.1 statusframe.com        # Cloudflare DNS

# Check www subdomain
dig www.statusframe.com
```

### Online Tools
- https://dnschecker.org/
- https://www.whatsmydns.net/
- https://mxtoolbox.com/SuperTool.aspx

## Propagation Time
- Usually: 5-30 minutes
- Maximum: 48 hours
- Cloudflare: Often under 5 minutes

## Subdomain Setup (Optional)

### API Subdomain
```
Type: A
Name: api
Value: YOUR_SERVER_IP
TTL: 3600
```

### Status Page Subdomain
```
Type: A
Name: status
Value: YOUR_SERVER_IP
TTL: 3600
```

### Development/Staging
```
Type: A
Name: staging
Value: YOUR_STAGING_SERVER_IP
TTL: 3600
```

## Email Configuration (Optional)

### MX Records for Email
```
Type: MX
Name: @
Priority: 10
Value: mail.your-email-provider.com
TTL: 3600
```

### SPF Record
```
Type: TXT
Name: @
Value: v=spf1 include:_spf.your-email-provider.com ~all
TTL: 3600
```

## Common Issues

### DNS Not Resolving
- Wait for propagation (up to 48 hours)
- Clear DNS cache: `sudo systemd-resolve --flush-caches`
- Check if nameservers are correct
- Verify no typos in records

### Certificate Errors After DNS Setup
- Wait for DNS to fully propagate before running certbot
- Verify DNS is pointing to server: `dig +short your-domain.com`
- Ensure firewall allows port 80 and 443

### Cloudflare Specific
- If using Cloudflare proxy (orange cloud):
  - SSL/TLS mode should be "Full" or "Full (strict)"
  - Don't use Cloudflare's Universal SSL initially
  - Setup your own SSL first, then enable proxy

## Example Complete DNS Configuration

```
# Main domain
Type: A    | Name: @   | Value: 203.0.113.10  | TTL: 3600

# WWW subdomain
Type: A    | Name: www | Value: 203.0.113.10  | TTL: 3600

# API subdomain (optional)
Type: A    | Name: api | Value: 203.0.113.10  | TTL: 3600

# IPv6 (optional)
Type: AAAA | Name: @   | Value: 2001:db8::1   | TTL: 3600

# Email (if using)
Type: MX   | Name: @   | Priority: 10 | Value: mail.provider.com | TTL: 3600

# SPF (if using email)
Type: TXT  | Name: @   | Value: v=spf1 include:_spf.provider.com ~all

# DMARC (if using email)
Type: TXT  | Name: _dmarc | Value: v=DMARC1; p=none; rua=mailto:admin@statusframe.com
```

## Testing After Setup

```bash
# 1. Verify DNS resolution
dig statusframe.com

# 2. Test HTTP (before SSL)
curl http://statusframe.com

# 3. After SSL setup, test HTTPS
curl https://statusframe.com

# 4. Test SSL certificate
openssl s_client -connect statusframe.com:443 -servername statusframe.com

# 5. Check HTTP to HTTPS redirect
curl -I http://statusframe.com
# Should return 301 redirect
```

## Next Steps

After DNS is configured and propagated:
1. ✅ Install Nginx
2. ✅ Configure Nginx for HTTP (temporary)
3. ✅ Get SSL certificate from Let's Encrypt
4. ✅ Update Nginx config for HTTPS
5. ✅ Test the application

See `docs/PRODUCTION_SETUP.md` for complete deployment guide.
