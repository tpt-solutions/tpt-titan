# TPT Titan

TPT Titan is a complete open source alternative to Microsoft Office 365, recreated entirely from scratch under GPL v3. It provides a modular, server-optimized suite of productivity tools including office applications, email, file synchronization, calendar, contacts, and video conferencing. Designed for self-hosting on VPS or local servers, with cross-platform desktop and web clients.

## Features

- **Office Suite**: Text editor, spreadsheet, presentation tools
- **Email Client**: Full-featured email with SMTP/IMAP support
- **File Sync**: Secure, peer-to-peer file synchronization
- **Calendar & Contacts**: Integrated scheduling and contact management
- **Video Conferencing**: Real-time video calls with WebRTC
- **Security & Privacy**: End-to-end encryption, self-hosted data control
- **Modular Architecture**: Docker-based deployment with simple install alternatives
- **Cross-Platform**: Web and desktop clients for Windows, Linux, Mac

## Tech Stack

- **Backend**: Go with Gin framework
- **Frontend**: Svelte
- **Desktop**: Tauri
- **Database**: PostgreSQL
- **Deployment**: Docker Compose (primary), install scripts (alternative)
- **Real-time**: WebSockets
- **Encryption**: Custom implementation with standard libraries

## Installation

### Prerequisites
- Docker and Docker Compose (for containerized deployment)
- Go 1.19+ (for development)
- Node.js 18+ (for frontend/desktop)
- PostgreSQL (for database)

### Quick Start with Docker
```bash
git clone https://github.com/yourorg/tpt-titan.git
cd tpt-titan
docker-compose up -d
```

### Manual Installation
Follow the [Installation Guide](docs/installation.md) for non-Docker setup.

## Development

### Setup
```bash
git clone https://github.com/yourorg/tpt-titan.git
cd tpt-titan
# Backend
cd backend
go mod tidy
go run main.go
# Frontend
cd frontend
npm install
npm run dev
# Desktop
cd desktop
npm install
npm run tauri dev
```

### Project Structure
```
tpt-titan/
├── backend/          # Go/Gin API server
├── frontend/         # Svelte web client
├── desktop/          # Tauri desktop app
├── shared/           # Shared utilities
├── docker/           # Docker configurations
├── docs/             # Documentation
├── scripts/          # Install/deployment scripts
└── tests/            # Test suites
```

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md).

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Submit a pull request

### Community
- [GitHub Issues](https://github.com/yourorg/tpt-titan/issues)
- [Discord Server](https://discord.gg/tpt-titan)
- [Documentation](https://docs.tpt-titan.org)

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Roadmap

See [TODO.md](TODO.md) for the complete development checklist and roadmap.

## Support

- [Documentation](https://docs.tpt-titan.org)
- [FAQ](docs/faq.md)
- [Donate](https://donate.tpt-titan.org) to support development

TPT Titan - Freedom through open source productivity.
