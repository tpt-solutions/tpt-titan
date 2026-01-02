# TPT Titan Documentation

Welcome to the official documentation for TPT Titan, the complete open-source alternative to Microsoft Office 365.

## 📖 Documentation Overview

This documentation covers everything you need to know about TPT Titan, from installation and setup to advanced usage and development.

## 🚀 Quick Start

### Installation

Choose your preferred installation method:

#### Docker (Recommended)
```bash
# Download and run with Docker Compose
wget https://raw.githubusercontent.com/tpt-titan/tpt-titan/main/docker-compose.yml
docker-compose up -d
```

#### Manual Installation
```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/tpt-titan/tpt-titan/main/scripts/install.sh | bash

# Windows
# Download the Windows installer from releases
```

### First Login

1. Open your browser to `http://localhost:8080`
2. Create your admin account
3. Start exploring the features!

## 📋 Features

### Core Applications
- **📊 Spreadsheet**: Excel-compatible spreadsheet with formulas, charts, and collaboration
- **📝 Document Editor**: Word-compatible editor with DOCX export and rich formatting
- **📋 Forms & Database**: Access-style database with visual form builder
- **📅 Calendar**: Full calendar with sharing and notifications
- **📧 Email Client**: Complete email client with PGP encryption
- **👥 Contacts**: Contact management with vCard import/export
- **💬 Chat**: Real-time messaging and collaboration
- **✅ Tasks**: Kanban-style task management

### Advanced Features
- **🔐 End-to-End Encryption**: Zero-knowledge encryption for all data
- **🔌 Plugin System**: Extensible architecture with plugin marketplace
- **📱 Cross-Platform**: Web, desktop (Windows/Mac/Linux), and mobile support
- **☁️ File Sync**: P2P file synchronization across devices
- **🎥 Video Conferencing**: WebRTC-based meetings and screen sharing
- **🤖 AI Integration**: Built-in AI assistance for writing and analysis
- **📊 GDPR Compliance**: Complete privacy controls and data portability

## 📚 Documentation Sections

### User Guides
- [Getting Started](getting-started.md) - Basic setup and first steps
- [Spreadsheet Guide](spreadsheet.md) - Using the spreadsheet application
- [Document Editor](editor.md) - Writing and formatting documents
- [Forms & Database](forms.md) - Creating forms and managing data
- [Calendar](calendar.md) - Managing events and schedules
- [Email](email.md) - Using the email client
- [Tasks](tasks.md) - Project management and task tracking

### Administration
- [Installation](installation.md) - Detailed installation instructions
- [Configuration](configuration.md) - Server configuration options
- [Security](security.md) - Security best practices
- [Backup & Recovery](backup.md) - Data backup and restore procedures
- [Monitoring](monitoring.md) - System monitoring and alerts
- [Troubleshooting](troubleshooting.md) - Common issues and solutions

### Development
- [API Reference](api.md) - REST API documentation
- [Plugin Development](plugins.md) - Creating custom plugins
- [Contributing](contributing.md) - How to contribute to TPT Titan
- [Architecture](architecture.md) - System architecture overview

## 🆘 Support

### Community Support
- **Forum**: [Community Forum](https://forum.tpt-titan.com)
- **Discord**: [Discord Server](https://discord.gg/tpt-titan)
- **Matrix**: `#tpt-titan:matrix.org`

### Professional Support
- **Enterprise Support**: Contact sales@tpt-titan.com
- **Security Issues**: Report to security@tpt-titan.com
- **Bug Reports**: [GitHub Issues](https://github.com/tpt-titan/tpt-titan/issues)

## 🤝 Contributing

We welcome contributions from the community! Please see our [Contributing Guide](contributing.md) for details on:

- Reporting bugs
- Suggesting features
- Submitting code changes
- Writing documentation
- Translating the interface

## 📄 License

TPT Titan is licensed under the **AGPL v3.0 or later**. This ensures that TPT Titan remains free and open-source forever.

## 🙏 Acknowledgments

TPT Titan is built on the shoulders of many amazing open-source projects. Special thanks to:

- **Go** - The programming language that powers our backend
- **PostgreSQL** - Our robust database
- **Redis** - High-performance caching and sessions
- **Svelte** - Our responsive frontend framework
- **Tauri** - Cross-platform desktop app framework
- **Gin** - HTTP web framework for Go

---

**Ready to get started?** [Install TPT Titan](installation.md) and begin your productivity journey!

*TPT Titan - Freedom in Productivity*
