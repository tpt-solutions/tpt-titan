#!/bin/bash

# TPT Titan Installation Script for Linux/macOS
# This script installs TPT Titan without Docker

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="/opt/tpt-titan"
DATA_DIR="/var/lib/tpt-titan"
CONFIG_DIR="/etc/tpt-titan"
LOG_DIR="/var/log/tpt-titan"
USER="tpt-titan"
GROUP="tpt-titan"

# Resolve the directory this script lives in (used to locate deploy/ templates)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Version and URLs (would be dynamically set)
VERSION="1.0.0"
DOWNLOAD_URL="https://github.com/tpt-titan/tpt-titan/releases/download/v${VERSION}"

# Detect OS and architecture
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        OS="linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        OS="darwin"
    else
        echo -e "${RED}Unsupported OS: $OSTYPE${NC}"
        exit 1
    fi
}

detect_arch() {
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="arm"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
}

print_banner() {
    echo -e "${BLUE}"
    cat << 'EOF'
████████╗██████╗ ████████╗    ████████╗██╗████████╗ █████╗ ███╗   ██╗
╚══██╔══╝██╔══██╗╚══██╔══╝    ╚══██╔══╝██║╚══██╔══╝██╔══██╗████╗  ██║
   ██║   ██████╔╝   ██║          ██║   ██║   ██║   ███████║██╔██╗ ██║
   ██║   ██╔═══╝    ██║          ██║   ██║   ██║   ██╔══██║██║╚██╗██║
   ██║   ██║        ██║          ██║   ██║   ██║   ██║  ██║██║ ╚████║
   ╚═╝   ╚═╝        ╚═╝          ╚═╝   ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═══╝
EOF
    echo -e "${NC}"
    echo -e "${GREEN}TPT Titan Installation Script${NC}"
    echo -e "${YELLOW}Version: ${VERSION}${NC}"
    echo ""
}

check_requirements() {
    echo -e "${BLUE}Checking system requirements...${NC}"

    # Check if running as root
    if [[ $EUID -eq 0 ]]; then
        echo -e "${YELLOW}Warning: Running as root. Consider using a non-privileged user.${NC}"
    fi

    # Check available disk space (need at least 2GB)
    AVAILABLE_SPACE=$(df /opt | tail -1 | awk '{print $4}')
    if [[ $AVAILABLE_SPACE -lt 2097152 ]]; then # 2GB in KB
        echo -e "${RED}Error: Insufficient disk space. Need at least 2GB available.${NC}"
        exit 1
    fi

    # Check for required commands
    REQUIRED_COMMANDS=("curl" "tar" "systemctl")
    for cmd in "${REQUIRED_COMMANDS[@]}"; do
        if ! command -v $cmd &> /dev/null; then
            echo -e "${RED}Error: Required command '$cmd' not found.${NC}"
            exit 1
        fi
    done

    echo -e "${GREEN}✓ System requirements met${NC}"
}

create_user() {
    echo -e "${BLUE}Creating TPT Titan system user...${NC}"

    # Create group if it doesn't exist
    if ! getent group $GROUP > /dev/null 2>&1; then
        groupadd -r $GROUP
    fi

    # Create user if it doesn't exist
    if ! id $USER > /dev/null 2>&1; then
        useradd -r -g $GROUP -d $DATA_DIR -s /bin/false $USER
    fi

    echo -e "${GREEN}✓ System user created${NC}"
}

create_directories() {
    echo -e "${BLUE}Creating directories...${NC}"

    # Create directories
    mkdir -p $INSTALL_DIR
    mkdir -p $DATA_DIR
    mkdir -p $CONFIG_DIR
    mkdir -p $LOG_DIR
    mkdir -p $DATA_DIR/uploads
    mkdir -p $DATA_DIR/backups
    mkdir -p $DATA_DIR/plugins

    # Set ownership
    chown -R $USER:$GROUP $INSTALL_DIR
    chown -R $USER:$GROUP $DATA_DIR
    chown -R $USER:$GROUP $LOG_DIR

    # Set permissions
    chmod 755 $INSTALL_DIR
    chmod 755 $DATA_DIR
    chmod 755 $CONFIG_DIR
    chmod 755 $LOG_DIR

    echo -e "${GREEN}✓ Directories created${NC}"
}

download_and_extract() {
    local filename="tpt-titan-${VERSION}-${OS}-${ARCH}.tar.gz"
    local download_path="/tmp/$filename"

    # Try to download the prebuilt release tarball. Releases are published on the
    # GitHub Releases page; until they exist for a given version/tag this will fail
    # and we transparently fall back to building from the local source tree below.
    echo -e "${BLUE}Downloading TPT Titan ${VERSION}...${NC}"
    if curl -fL -o "$download_path" "${DOWNLOAD_URL}/${filename}"; then
        echo -e "${BLUE}Extracting files...${NC}"
        if tar -xzf "$download_path" -C "$INSTALL_DIR"; then
            rm -f "$download_path"
            echo -e "${GREEN}✓ TPT Titan downloaded and extracted${NC}"
            return 0
        fi
        rm -f "$download_path"
        echo -e "${YELLOW}Downloaded archive failed to extract; falling back to local build.${NC}"
    else
        echo -e "${YELLOW}Prebuilt release not available for v${VERSION}/${OS}/${ARCH}.${NC}"
    fi

    build_from_source
}

install_dependencies() {
    echo -e "${BLUE}Installing system dependencies...${NC}"

    if command -v apt-get > /dev/null 2>&1; then
        # Debian/Ubuntu
        apt-get update
        apt-get install -y postgresql postgresql-contrib redis-server nginx certbot python3-certbot-nginx
    elif command -v yum > /dev/null 2>&1; then
        # CentOS/RHEL
        yum update -y
        yum install -y postgresql-server postgresql-contrib redis nginx certbot python3-certbot-nginx
        postgresql-setup initdb
    elif command -v dnf > /dev/null 2>&1; then
        # Fedora
        dnf update -y
        dnf install -y postgresql-server postgresql-contrib redis nginx certbot python3-certbot-nginx
        postgresql-setup --initdb
    elif command -v zypper > /dev/null 2>&1; then
        # openSUSE
        zypper refresh
        zypper install -y postgresql-server postgresql-contrib redis nginx certbot
    elif command -v pacman > /dev/null 2>&1; then
        # Arch Linux
        pacman -Syu --noconfirm
        pacman -S --noconfirm postgresql redis nginx certbot
    elif [[ "$OS" == "darwin" ]]; then
        # macOS
        if ! command -v brew > /dev/null 2>&1; then
            echo -e "${YELLOW}Homebrew not found. Please install Homebrew first: https://brew.sh/${NC}"
            echo -e "${BLUE}Installing PostgreSQL and Redis via Homebrew...${NC}"
            brew install postgresql redis nginx
        else
            brew install postgresql redis nginx
        fi
    else
        echo -e "${YELLOW}Warning: Unsupported package manager. Please install dependencies manually:${NC}"
        echo "  - PostgreSQL"
        echo "  - Redis"
        echo "  - Nginx"
        echo "  - SSL certificate tool (certbot)"
    fi

    echo -e "${GREEN}✓ Dependencies installed${NC}"
}

setup_database() {
    echo -e "${BLUE}Setting up PostgreSQL database...${NC}"

    # Create database and user
    sudo -u postgres psql -c "CREATE USER tpt_titan WITH PASSWORD 'secure_password_change_me';" 2>/dev/null || true
    sudo -u postgres psql -c "CREATE DATABASE tpt_titan OWNER tpt_titan;" 2>/dev/null || true
    sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE tpt_titan TO tpt_titan;" 2>/dev/null || true

    echo -e "${GREEN}✓ Database setup completed${NC}"
}

create_config() {
    echo -e "${BLUE}Creating configuration files...${NC}"

    # Create main configuration file
    cat > $CONFIG_DIR/config.yaml << EOF
# TPT Titan Configuration
server:
  host: "0.0.0.0"
  port: "8080"
  mode: "release"

database:
  host: "localhost"
  port: "5432"
  user: "tpt_titan"
  password: "secure_password_change_me"
  dbname: "tpt_titan"
  sslmode: "disable"

redis:
  host: "localhost"
  port: "6379"
  password: ""
  db: 0

jwt:
  secret: "$(openssl rand -hex 32)"
  expiry: "24h"

encryption:
  master_key: "$(openssl rand -hex 32)"

ai:
  ollama_url: "http://localhost:11434"
  default_model: "llama2"

file_storage:
  path: "$DATA_DIR/uploads"
  max_size: "100MB"

plugins:
  enabled: true
  directory: "$DATA_DIR/plugins"

logging:
  level: "info"
  file: "$LOG_DIR/tpt-titan.log"
EOF

    # Set proper permissions
    chown $USER:$GROUP $CONFIG_DIR/config.yaml
    chmod 600 $CONFIG_DIR/config.yaml

    # Create environment file
    cat > $CONFIG_DIR/.env << EOF
# Environment variables for TPT Titan
DATABASE_URL=postgresql://tpt_titan:secure_password_change_me@localhost:5432/tpt_titan
REDIS_URL=redis://localhost:6379
JWT_SECRET=$(openssl rand -hex 32)
ENCRYPTION_KEY=$(openssl rand -hex 32)
EOF

    chown $USER:$GROUP $CONFIG_DIR/.env
    chmod 600 $CONFIG_DIR/.env

    echo -e "${GREEN}✓ Configuration files created${NC}"
}

create_service() {
    echo -e "${BLUE}Creating systemd service...${NC}"

    # Use the checked-in, reviewable service template if present; otherwise fall back
    # to a generated unit (kept in sync with deploy/tpt-titan.service).
    local template="$SCRIPT_DIR/../deploy/tpt-titan.service"
    if [[ -f "$template" ]]; then
        sed -e "s#/opt/tpt-titan#$INSTALL_DIR#g" \
            -e "s#/var/lib/tpt-titan#$DATA_DIR#g" \
            -e "s#/var/log/tpt-titan#$LOG_DIR#g" \
            -e "s#/etc/tpt-titan/.env#$CONFIG_DIR/.env#g" \
            -e "s#User=tpt-titan#User=$USER#g" \
            -e "s#Group=tpt-titan#Group=$GROUP#g" \
            "$template" > /etc/systemd/system/tpt-titan.service
    else
        cat > /etc/systemd/system/tpt-titan.service << EOF
[Unit]
Description=TPT Titan
After=network.target

[Service]
Type=simple
User=$USER
Group=$GROUP
WorkingDirectory=$INSTALL_DIR
ExecStart=$INSTALL_DIR/tpt-titan
Restart=always
RestartSec=5
EnvironmentFile=$CONFIG_DIR/.env

# Security settings
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ReadWritePaths=$DATA_DIR $LOG_DIR
ProtectHome=yes

# Resource limits
MemoryLimit=1G
CPUQuota=50%

[Install]
WantedBy=multi-user.target
EOF
    fi

    # Reload systemd and enable service
    systemctl daemon-reload
    systemctl enable tpt-titan

    echo -e "${GREEN}✓ Systemd service created${NC}"
}

setup_nginx() {
    echo -e "${BLUE}Configuring Nginx...${NC}"

    # Create Nginx configuration
    cat > /etc/nginx/sites-available/tpt-titan << EOF
server {
    listen 80;
    server_name _;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "no-referrer-when-downgrade" always;
    add_header Content-Security-Policy "default-src 'self' http: https: data: blob: 'unsafe-inline'" always;

    # Handle static files
    location /static/ {
        alias $INSTALL_DIR/static/;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Main application
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;

        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # WebSocket support
    location /ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

    # Enable site
    ln -sf /etc/nginx/sites-available/tpt-titan /etc/nginx/sites-enabled/
    rm -f /etc/nginx/sites-enabled/default

    # Test configuration
    nginx -t

    echo -e "${GREEN}✓ Nginx configured${NC}"
}

setup_ssl() {
    echo -e "${BLUE}Setting up SSL certificate...${NC}"

    # Prompt for domain
    read -p "Enter your domain name (leave blank to skip SSL setup): " domain

    if [[ -n "$domain" ]]; then
        echo -e "${YELLOW}Note: Make sure your domain DNS points to this server${NC}"
        echo -e "${BLUE}Obtaining SSL certificate for $domain...${NC}"

        certbot --nginx -d $domain --non-interactive --agree-tos --email admin@$domain

        echo -e "${GREEN}✓ SSL certificate obtained${NC}"
    else
        echo -e "${YELLOW}Skipping SSL setup. You can run SSL setup later with:${NC}"
        echo "  certbot --nginx -d yourdomain.com"
    fi
}

initialize_database() {
    echo -e "${BLUE}Initializing database...${NC}"

    # Run database migrations/schema setup
    if [[ -f "$INSTALL_DIR/scripts/init-db.sql" ]]; then
        sudo -u postgres psql -d tpt_titan -f $INSTALL_DIR/scripts/init-db.sql
    fi

    echo -e "${GREEN}✓ Database initialized${NC}"
}

start_services() {
    echo -e "${BLUE}Starting services...${NC}"

    # Start and enable services
    systemctl start postgresql
    systemctl enable postgresql

    systemctl start redis
    systemctl enable redis

    systemctl start tpt-titan
    systemctl enable tpt-titan

    systemctl start nginx
    systemctl enable nginx

    echo -e "${GREEN}✓ Services started${NC}"
}

show_completion() {
    echo ""
    echo -e "${GREEN}🎉 TPT Titan installation completed successfully!${NC}"
    echo ""
    echo -e "${BLUE}What's next:${NC}"
    echo "1. Update the database password in $CONFIG_DIR/config.yaml"
    echo "2. Update JWT and encryption secrets in $CONFIG_DIR/.env"
    echo "3. Configure your domain in Nginx if needed"
    echo "4. Access TPT Titan at: http://your-server-ip"
    echo ""
    echo -e "${BLUE}Useful commands:${NC}"
    echo "  systemctl status tpt-titan    # Check service status"
    echo "  journalctl -u tpt-titan       # View logs"
    echo "  systemctl restart tpt-titan   # Restart service"
    echo ""
    echo -e "${YELLOW}Security recommendations:${NC}"
    echo "• Change the default database password"
    echo "• Set up firewall rules (ufw or firewalld)"
    echo "• Configure fail2ban for additional security"
    echo "• Regularly update the system and TPT Titan"
    echo ""
    echo -e "${GREEN}Thank you for installing TPT Titan!${NC}"
}

main() {
    detect_os
    detect_arch
    print_banner

    echo -e "${YELLOW}This script will install TPT Titan on your system.${NC}"
    echo -e "${YELLOW}Please ensure you have backups of important data.${NC}"
    echo ""
    read -p "Do you want to continue? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi

    check_requirements
    create_user
    create_directories
    install_dependencies
    download_and_extract
    setup_database
    create_config
    create_service
    setup_nginx
    setup_ssl
    initialize_database
    start_services
    show_completion
}

build_from_source() {
    echo -e "${BLUE}Building TPT Titan from local source...${NC}"

    local repo_root="$SCRIPT_DIR/.."
    local backend_dir="$repo_root/backend"

    if [[ ! -d "$backend_dir" ]]; then
        echo -e "${RED}Error: backend source not found at $backend_dir.${NC}"
        echo -e "${RED}Cannot install without either a release tarball or the source tree.${NC}"
        exit 1
    fi

    if ! command -v go &> /dev/null; then
        echo -e "${RED}Error: 'go' is required to build from source but was not found.${NC}"
        exit 1
    fi

    echo -e "${BLUE}Compiling backend (this may take a minute)...${NC}"
    ( cd "$backend_dir" && go build -o "$INSTALL_DIR/tpt-titan" . ) || {
        echo -e "${RED}Error: backend build failed.${NC}"
        exit 1
    }

    # Copy static assets / config templates if present
    [[ -d "$repo_root/deploy" ]] && cp -r "$repo_root/deploy" "$INSTALL_DIR/" 2>/dev/null || true

    chown -R $USER:$GROUP "$INSTALL_DIR"
    echo -e "${GREEN}✓ TPT Titan built and installed from source${NC}"
}

# Run main function
main "$@"
