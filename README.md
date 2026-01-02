# TPT Titan

TPT Titan is a complete open source alternative to Microsoft Office 365, recreated entirely from scratch under GNU AGPL v3.0. It provides a comprehensive, self-hosted productivity suite designed for privacy-conscious users and organizations who want full control over their data.

## Features

### Office Suite
- **Text Editor**: Rich text editing with AI writing assistance, Notion-style block editing, Markdown support, natural math notation, handwriting recognition, and export to PDF/DOCX/LaTeX/MathML/SVG
- **Spreadsheet**: Advanced grid interface with mathematical functions (50+ formulas), AI formula suggestions, real-time collaboration, data visualization (charts), Excel import/export, and version control
- **Forms & Templates**: Visual drag-and-drop form builder (12+ field types), MS Access-style database features, visual query builder, relationship management, report generation, workflow automation, digital signatures, and email distribution

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

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

The AGPL ensures that modifications to the software remain open source and available to the community, maintaining freedom for all users.

## Roadmap

TPT Titan is actively developed with regular releases. See [TODO.md](TODO.md) for the current development status and upcoming features.

### Recent Major Releases
- ✅ End-to-end encryption system
- ✅ Complete office suite (text editor, spreadsheet, forms)
- ✅ Email client with PGP encryption
- ✅ Real-time collaboration features
- ✅ Plugin system and extensibility
- ✅ GDPR compliance implementation
- ✅ Cross-platform desktop applications

### Upcoming
- Mobile applications
- Advanced AI integrations
- Enhanced collaboration features
- Performance optimizations

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
