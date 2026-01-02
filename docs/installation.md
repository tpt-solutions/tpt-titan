# TPT Titan Installation Guide

This guide covers both development setup and production deployment of TPT Titan.

## Prerequisites

- **For Development:**
  - Go 1.19+ (for backend)
  - Node.js 18+ (for frontend/desktop)
  - SQLite (built-in) or PostgreSQL (optional)
  - Git

- **For Production:**
  - Docker and Docker Compose (recommended)
  - Or: Linux/macOS server with Go, Node.js, PostgreSQL, Redis, Nginx

## Development Setup

### 1. Clone the Repository
```bash
git clone https://github.com/yourorg/tpt-titan.git
cd tpt-titan
```

### 2. Backend Setup
```bash
cd backend
go mod tidy

# Copy environment configuration
cp ../.env.example .env
# Edit .env with your settings

# Initialize database (SQLite default)
go run main.go --init-db

# Run backend server
go run main.go
```

### 3. Frontend Setup
```bash
cd frontend
npm install

# Copy environment configuration
cp ../.env.example .env.local
# Edit .env.local with your API endpoint

# Run development server
npm run dev
```

### 4. Desktop Setup (Optional)
```bash
cd desktop
npm install

# Run desktop app in development
npm run tauri dev
```

### 5. Database Setup (Optional - PostgreSQL)
If using PostgreSQL instead of SQLite:

```bash
# Install PostgreSQL and create database
sudo -u postgres createdb tpt_titan
sudo -u postgres createuser tpt_titan
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE tpt_titan TO tpt_titan;"

# Update backend/.env with PostgreSQL connection details
```

## Production Deployment

### Option 1: Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/yourorg/tpt-titan.git
cd tpt-titan

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Option 2: Manual Installation

For manual installation on Linux/macOS servers, use the provided installation script:

```bash
# Download and run the installation script
curl -fsSL https://raw.githubusercontent.com/yourorg/tpt-titan/main/scripts/install.sh | bash

# Or clone and run locally
git clone https://github.com/yourorg/tpt-titan.git
cd tpt-titan
chmod +x scripts/install.sh
sudo ./scripts/install.sh
```

The script will:
- Install system dependencies (PostgreSQL, Redis, Nginx)
- Set up database and users
- Configure systemd services
- Set up SSL certificates
- Start all services

### Option 3: Binary Installation

Download the latest release binaries for your platform:

```bash
# Download the appropriate binary for your system
# Linux: tpt-titan-linux-amd64.tar.gz
# macOS: tpt-titan-darwin-amd64.tar.gz
# Windows: tpt-titan-windows-amd64.zip

# Extract and run
tar -xzf tpt-titan-linux-amd64.tar.gz
cd tpt-titan
./tpt-titan
```

## Configuration

### Environment Variables

Create a `.env` file in the backend directory:

```env
# Database
DATABASE_URL=sqlite://./data/tpt-titan.db
# Or for PostgreSQL:
# DATABASE_URL=postgresql://user:password@localhost:5432/tpt_titan

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-super-secret-jwt-key-here
ENCRYPTION_KEY=your-encryption-key-here

# Server
PORT=8080
HOST=0.0.0.0

# AI Configuration
OLLAMA_URL=http://localhost:11434
DEFAULT_AI_MODEL=llama2

# Email (optional)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

### Nginx Configuration (Production)

Example Nginx configuration for reverse proxy:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL configuration
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN";
    add_header X-Content-Type-Options "nosniff";
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains";

    # Main application
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }

    # Static files
    location /static/ {
        alias /opt/tpt-titan/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check DATABASE_URL in .env file
   - Ensure PostgreSQL is running: `sudo systemctl status postgresql`
   - Verify user permissions

2. **Frontend Can't Connect to Backend**
   - Check API_BASE_URL in frontend/.env.local
   - Ensure backend is running on correct port
   - Check CORS settings

3. **Permission Errors**
   - Run with appropriate user permissions
   - Check file ownership: `chown -R tpt-titan:tpt-titan /opt/tpt-titan`

4. **SSL Certificate Issues**
   - Run: `certbot --nginx -d your-domain.com`
   - Ensure DNS points to your server

### Logs

```bash
# Docker logs
docker-compose logs backend
docker-compose logs frontend

# Systemd logs
journalctl -u tpt-titan
journalctl -u nginx

# Application logs
tail -f /var/log/tpt-titan/tpt-titan.log
```

## Updating

### Docker Deployment
```bash
cd tpt-titan
git pull
docker-compose down
docker-compose pull
docker-compose up -d
```

### Manual Installation
```bash
cd /opt/tpt-titan
git pull
sudo systemctl restart tpt-titan
```

## Security Considerations

- Change default passwords and secrets
- Use strong JWT and encryption keys
- Keep system and dependencies updated
- Configure firewall rules
- Use HTTPS in production
- Regularly backup data
- Monitor logs for suspicious activity

## Support

- [Documentation](https://docs.tpt-titan.org)
- [GitHub Issues](https://github.com/yourorg/tpt-titan/issues)
- [Community Forum](https://community.tpt-titan.org)
