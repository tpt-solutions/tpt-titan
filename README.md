# TPT Titan

TPT Titan is a complete open source productivity suite designed specifically for **small businesses (1-7 employees)** who need professional tools but want to **avoid expensive subscriptions** and maintain **complete data privacy**. Built with offline-first architecture, it runs entirely on your local machine with zero internet dependency for core functionality.

## 🎯 Perfect for SMEs

- ✅ **Free & Open Source** - No subscriptions, no hidden costs
- ✅ **Offline-First** - Works completely without internet
- ✅ **Data Privacy** - Your data stays on your computer
- ✅ **Professional Features** - Excel-class spreadsheet, Word-class editor
- ✅ **Zero Vendor Lock-in** - Full data export and control

## Features

### ✅ **Implemented & Ready**
- **Text Editor**: Professional rich text editor with formatting, find/replace, templates, and PDF export
- **Spreadsheet**: Excel-class spreadsheet with 50+ formulas, data visualization, and import/export
- **Forms Builder**: Visual drag-and-drop form creation with 12+ field types, conditional logic, and templates
- **Task Management**: Kanban board with drag-and-drop, project tracking, and team collaboration
- **Desktop App**: Cross-platform desktop application (Windows/Mac/Linux) with offline-first architecture

### 🔄 **Backend Integration In Progress**
- **End-to-End Encryption**: AES-256-GCM encryption system with zero-knowledge architecture
- **Local Database**: SQLite integration for offline data storage
- **API Integration**: RESTful APIs connecting frontend components to encrypted backend

### 🚀 **Upcoming Features**
- **AI Integration**: Local AI models for writing assistance and data analysis
- **Real-time Collaboration**: Multi-user editing with conflict resolution
- **Email Client**: IMAP/SMTP with PGP encryption
- **Calendar & Contacts**: Personal information management
- **Plugin System**: Extensible architecture for custom features

### Communication & Collaboration
- **Email Client**: Full IMAP/SMTP support with PGP encryption, attachments, search, folders, and contact integration
- **Calendar**: Event management, multiple calendar support, notifications, sharing with permissions, and email integration
- **Contacts**: Contact management with vCard import/export, groups, categories, and integration with email/calendar
- **Chat**: Real-time messaging with WebSocket support, rooms (direct/group/channel), message reactions, threading, and participant management
- **Video Conferencing**: WebRTC-based meetings with screen sharing, participant management, and recording framework

### File Management & Sync
- **File Synchronization**: P2P file sync with versioning, conflict resolution, selective sync, and bandwidth optimization
- **Document Export**: Export capabilities for various formats (PDF, DOCX, Excel, etc.)

### AI & Intelligence
- **AI Integration**: Hardware-optimized AI with automatic model selection (Qwen 3, Qwen 2.5, quantized models), smart writing assistance, data analysis, OCR, and task management
- **Task Management**: Kanban board with AI task suggestions, integration with forms and email, project tracking, and workflow automation

### Security & Privacy
- **End-to-End Encryption**: Cryptographic key management with Shamir's secret sharing, hardware recovery (USB/face ID), and zero-knowledge backup
- **GDPR Compliance**: Complete GDPR implementation including right to erasure, data portability, consent management, and privacy controls
- **Authentication**: JWT-based auth with 2FA/MFA (TOTP), secure session management, and admin console
- **Security Features**: Input validation, CSRF protection, rate limiting, CORS, security headers, and audit logging

### System Features
- **Plugin System**: Extensible architecture with Go plugin system, event-driven communication, sandboxing, and marketplace
- **Backup & Recovery**: Comprehensive backup service with compression, retention, and restore capabilities
- **Performance**: Redis caching, database optimization, monitoring, and horizontal scaling support
- **Cross-Platform**: Web and desktop clients for Windows, Linux, and macOS

## Tech Stack

- **Backend**: Go with Gin framework, WebSockets, PostgreSQL/SQLite
- **Frontend**: Svelte with Tailwind CSS
- **Desktop**: Tauri (Rust)
- **Database**: SQLite (default) or PostgreSQL
- **Cache**: Redis
- **Deployment**: Docker Compose, systemd services, binary releases
- **Real-time**: WebSocket connections
- **Encryption**: Custom implementation with standard cryptographic libraries
- **AI**: Ollama integration with local model support

## Installation

### Prerequisites
- Docker and Docker Compose (recommended for production)
- Go 1.19+ (for development)
- Node.js 18+ (for frontend/desktop development)
- SQLite (built-in default) or PostgreSQL (optional)

### Quick Start with Docker
```bash
git clone https://github.com/yourorg/tpt-titan.git
cd tpt-titan
docker-compose up -d
```

### Manual Installation
See the complete [Installation Guide](docs/installation.md) for development setup, production deployment, and troubleshooting.

## Development

### Local Development Setup
```bash
git clone https://github.com/yourorg/tpt-titan.git
cd tpt-titan

# Backend
cd backend
go mod tidy
cp ../.env.example .env
# Edit .env with your configuration
go run main.go

# Frontend (new terminal)
cd frontend
npm install
cp ../.env.example .env.local
npm run dev

# Desktop (optional, new terminal)
cd desktop
npm install
npm run tauri dev
```

### Project Structure
```
tpt-titan/
├── backend/                 # Go/Gin API server
│   ├── config/             # Configuration management
│   ├── middleware/         # Security and request middleware
│   ├── models/             # Database models
│   ├── routes/             # API endpoints
│   ├── services/           # Business logic services
│   ├── utils/              # Utility functions
│   ├── tests/              # Integration tests
│   └── websocket.go        # Real-time communication
├── frontend/                # Svelte web client
│   ├── src/
│   │   ├── lib/
│   │   │   ├── components/ # Reusable UI components
│   │   │   ├── stores.js   # State management
│   │   │   └── api.js      # API client
│   │   └── routes/         # Page routes
│   └── package.json
├── desktop/                 # Tauri desktop application
│   ├── src-tauri/          # Rust backend
│   └── package.json
├── docs/                    # Documentation
├── scripts/                 # Installation and utility scripts
│   ├── install.sh          # Production installer
│   ├── init-db.sql         # Database initialization
│   └── create-db-schema.sql
├── docker-compose.yml       # Container orchestration
├── .env.example            # Environment configuration template
└── README.md
```

## Contributing

We welcome contributions from the community! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- Development workflow
- Code standards and guidelines
- Testing requirements
- Pull request process

### Code of Conduct
Please read our [Code of Conduct](CODE_OF_CONDUCT.md) to understand our community standards.

## Community & Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/yourorg/tpt-titan/issues)
- **Documentation**: [Complete user and developer docs](https://docs.tpt-titan.org)
- **Community Forum**: Join discussions and get help
- **Discord**: Real-time chat with the community

## License

TPT Titan is free and open source, dual-licensed under your choice of either:

- **Apache License, Version 2.0** — see [LICENSE-APACHE](LICENSE-APACHE)
- **MIT License** — see [LICENSE-MIT](LICENSE-MIT)

Use, modify, and distribute the software (including in closed-source or
commercial products) under the terms of whichever license you prefer.

## Current Development Status

### ✅ **Phase 1: Core Components Complete**
- **Text Editor**: Professional rich text editor with full formatting, find/replace, and PDF export
- **Spreadsheet**: Excel-class spreadsheet with 50+ formulas, charts, and import/export
- **Forms Builder**: Visual form creation with 12+ field types, conditional logic, templates
- **Task Management**: Kanban board with drag-and-drop functionality
- **Desktop App**: Cross-platform Tauri application framework ready

### 🔄 **Phase 2: Backend Integration (In Progress)**
- **Encryption Integration**: Connect AES-256-GCM crypto system to document storage
- **SQLite Database**: Implement local data persistence for all components
- **API Connections**: Wire frontend components to encrypted backend
- **Data Portability**: Export/import functionality for user data

### 🚀 **Phase 3: Production Release (Next)**
- **Desktop App Bundling**: Package Go backend with Tauri for seamless offline experience
- **User Onboarding**: Setup wizards and initial configuration
- **Documentation**: User guides and video tutorials
- **Community Launch**: GitHub release with desktop binaries

### 📅 **Phase 4: Advanced Features (Future)**
- **AI Integration**: Local AI models for writing assistance and data analysis
- **Real-time Collaboration**: Multi-user editing with conflict resolution
- **Email Client**: IMAP/SMTP with PGP encryption
- **Calendar & Contacts**: Personal information management
- **Plugin System**: Extensible architecture for custom features

### 🎯 **Immediate Priorities**
1. **Encryption Integration**: Make documents save/load with proper encryption
2. **Desktop App Release**: Package as downloadable application
3. **User Testing**: Get feedback from SME users
4. **Community Building**: Grow user base and contributor community

## Security

TPT Titan takes security seriously:

- **Encryption**: All data is encrypted at rest and in transit
- **Access Control**: Role-based permissions and secure authentication
- **Audit Logging**: Comprehensive security event tracking
- **Regular Updates**: Security patches and dependency updates
- **Privacy by Design**: GDPR compliance and data protection

For security-related issues, please see our [Security Policy](SECURITY.md).

## Support the Project

TPT Titan is free and open source. If you'd like to support development:

- ⭐ Star the repository on GitHub
- 🐛 Report bugs and suggest features
- 💻 Contribute code improvements
- 📖 Help improve documentation
- 💬 Participate in community discussions
- 💰 [Donate](https://donate.tpt-titan.org) to support ongoing development

---

**TPT Titan** - Complete open source productivity suite for the privacy-conscious user.
