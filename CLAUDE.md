# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

TPT Titan is a complete open-source alternative to Microsoft Office 365, built from scratch under GNU AGPL v3.0. It's a comprehensive, self-hosted productivity suite with office applications, communication tools, file management, AI integration, and enterprise security features.

**Tech Stack:**
- Backend: Go (1.19+) with Gin framework, GORM, WebSockets
- Frontend: SvelteKit with Tailwind CSS
- Desktop: Tauri (Rust) wrapping the web frontend
- Databases: SQLite (default) or PostgreSQL
- Cache: Redis (optional but recommended for production)
- AI: Ollama integration for local models

## Essential Commands

### Development

**Backend:**
```bash
cd backend
go mod tidy                    # Install/update dependencies
go run main.go                 # Run development server (port 8080)
go test ./...                  # Run all tests
go test -v ./tests/...         # Run tests with verbose output
```

**Frontend:**
```bash
cd frontend
npm install                    # Install dependencies
npm run dev                    # Run development server (port 3000)
npm run build                  # Build for production
npm run preview                # Preview production build
npm run check                  # Type-check without building
```

**Desktop:**
```bash
cd desktop
npm install
npm run tauri dev              # Run desktop app in development
npm run tauri build            # Build production desktop app
```

### Docker

```bash
docker-compose up -d           # Start all services (SQLite default)
docker-compose logs -f         # Follow logs
docker-compose down            # Stop all services
```

For PostgreSQL instead of SQLite, copy `docker-compose.override.yml.example` to `docker-compose.override.yml` and run docker-compose again.

### Testing

```bash
# Backend integration tests
cd backend
go test -v ./tests/integration_test.go

# Run specific test
go test -v -run TestSpecificFunction ./tests/...
```

## Architecture

### Service-Oriented Backend

The Go backend follows a **service-oriented architecture** where business logic is isolated in services (in `backend/services/`) and routes (in `backend/routes/`) handle HTTP concerns only:

- **Routes** (`backend/routes/*.go`): HTTP handlers, request parsing, response formatting
- **Services** (`backend/services/*.go`): Business logic, database operations, external integrations
- **Models** (`backend/models/*.go`): Database models with GORM tags
- **Middleware** (`backend/middleware/security.go`): Security, auth, logging, rate limiting
- **Config** (`backend/config/`): Configuration management and database connection

**Service Initialization Pattern:**
Services are initialized in `main.go` and passed to route handlers. Example:
```go
// In main.go
chatService := services.NewChatService(db)
routes.InitChatService(chatService)

// In routes/chat.go
var chatServiceInstance *services.ChatService
func InitChatService(service *services.ChatService) {
    chatServiceInstance = service
}
```

### Real-Time Communication

**WebSocket Hub Architecture** (`backend/websocket.go`):
- Central Hub manages all WebSocket connections
- Clients are organized by UserID (one user can have multiple clients)
- Message broadcasting based on room participants
- Supports message types: WSMessageSent, WSReactionAdded, WSTypingStart, WSUserStatus
- Integrated with ChatService for participant lookup

### Database Architecture

**Dual Database Support:**
The app supports both SQLite (default) and PostgreSQL via environment variable `DB_TYPE`:
- SQLite: Ideal for development and small deployments
- PostgreSQL: Required for production scale deployments

**Key Models:**
- User authentication with JWT + 2FA/TOTP support
- End-to-end encryption models (`models/encrypted.go`)
- Spreadsheet with formula evaluation, charts, version control
- Forms with visual query builder, relationships, workflows
- Email with IMAP/SMTP, PGP encryption, attachments
- Calendar with events, sharing, notifications
- Contacts with vCard import/export, groups
- Chat rooms (direct/group/channel) with messages, reactions
- Video conferencing with WebRTC support
- File sync with P2P protocol and versioning

**Database Connection:**
Connection is managed globally via `config.GetDatabase()`. Services receive db instance during initialization.

### Security Middleware Stack

Located in `backend/middleware/security.go`, applied globally to all routes:
1. **RequestIDMiddleware**: Unique request IDs for tracing
2. **CORSMiddleware**: Cross-origin resource sharing
3. **SecurityHeadersMiddleware**: X-Frame-Options, X-Content-Type-Options, CSP, HSTS
4. **RateLimitMiddleware**: Token bucket rate limiting per IP
5. **InputValidationMiddleware**: XSS, SQL injection prevention
6. **SQLInjectionProtectionMiddleware**: Additional SQL injection checks
7. **AuditMiddleware**: Security event logging

### Authentication & Authorization

**JWT-based Authentication:**
- Public routes: `/api/v1/auth/register`, `/api/v1/auth/login`, password reset
- Protected routes: All others require JWT token in Authorization header
- Admin routes: Require both valid JWT and admin role
- 2FA/TOTP: Optional per-user basis

**AuthMiddleware** (`backend/routes/auth/middleware.go`):
Sets `user_id` in Gin context, available via `c.Get("user_id")` in handlers.

### Frontend State Management

**Svelte Stores** (`frontend/src/lib/stores.js`):
- Writable stores for user, token, spreadsheet data
- API client in `frontend/src/lib/api.js` handles authentication headers

**Component Organization:**
- Reusable components in `frontend/src/lib/components/`
- Route-based pages in `frontend/src/routes/`
- Each major feature has dedicated components (Spreadsheet, Calendar, EmailInbox, etc.)

### AI Integration

**Hardware-Optimized AI** (`backend/services/ai.go`):
- Automatic hardware detection (CPU, GPU, Metal, CUDA)
- Model recommendation based on hardware capabilities
- Support for Qwen 3, Qwen 2.5, quantized models
- Ollama integration for local models
- Optional OpenRouter for cloud models

**AI Features:**
- Writing assistance in text editor
- Formula suggestions in spreadsheet
- Task suggestions and prioritization
- Handwriting recognition for math equations
- OCR capabilities

### Encryption Architecture

**End-to-End Encryption** (`backend/utils/crypto.go`, `backend/models/encrypted.go`):
- AES-256-GCM encryption for all sensitive data
- Argon2 key derivation from user passwords
- Shamir's secret sharing for key recovery
- Hardware recovery keys (USB/face ID support)
- Zero-knowledge backup architecture
- All models with encrypted fields implement EncryptedModel interface

### Office Suite Features

**Spreadsheet:**
- Formula evaluation engine with 50+ functions (SUM, AVERAGE, SIN, COS, etc.)
- Real-time collaboration with version control and cell locking
- Chart generation with auto-suggestions (bar, line, pie, scatter, area)
- Excel import/export with full .xlsx support including formulas and styles
- Service: `backend/services/spreadsheet_math.go`, `backend/services/excel.go`

**Text Editor:**
- Notion-style block editing
- Markdown and WYSIWYG dual mode
- Natural math notation (better than LaTeX)
- Handwriting recognition for equations (`backend/services/handwriting_recognition.go`)
- Export to PDF, DOCX, LaTeX, MathML, SVG (`backend/services/docx_export.go`)

**Forms & Templates:**
- Drag-and-drop form builder with 12+ field types
- MS Access-style visual query builder (`backend/services/form_query_builder.go`)
- Relationship management between forms (`backend/services/form_relationships.go`)
- Report generation and dashboards (`backend/services/form_reporting.go`)
- Workflow automation with approval chains (`backend/services/form_workflow.go`)
- Digital signatures and email distribution

### Communication Features

**Email:**
- IMAP/SMTP integration (`backend/services/email.go`)
- PGP encryption (`backend/services/pgp_encryption.go`)
- Attachment handling with virus scanning (`backend/services/email_attachments.go`)
- Contact integration

**Calendar:**
- Multi-calendar support
- Sharing with granular permissions (`backend/services/calendar_sharing.go`)
- Notifications via email/SMS/push (`backend/services/calendar_notifications.go`)
- Event attendees and RSVP

**Contacts:**
- vCard import/export (`backend/services/contact_import_export.go`)
- Groups and categories
- Integration with email and calendar

**Chat:**
- WebSocket-based real-time messaging
- Room types: direct, group, channel
- Message reactions and threading
- Participant management with roles

**Video Conferencing:**
- WebRTC-based meetings (`backend/services/videoconf.go`)
- Screen sharing
- Recording framework
- Participant management

### File Synchronization

**P2P Sync** (`backend/services/filesync.go`):
- File versioning and conflict resolution
- Selective sync
- Bandwidth optimization
- Cross-platform file watching (fsnotify)

### Performance & Monitoring

**Caching** (`backend/services/cache.go`):
- Redis-based caching with TTL
- Pub/sub support for real-time updates
- Pipeline operations for batch requests

**Database Optimization** (`backend/services/db_optimizer.go`):
- Connection pooling
- VACUUM for SQLite
- Index optimization
- Query performance monitoring

**Monitoring** (`backend/services/monitor.go`):
- System metrics (CPU, memory, disk)
- Performance reports
- Health checks at `/api/v1/monitoring/health`
- Prometheus metrics at `/api/v1/metrics`

### Plugin System

**Extensibility** (`backend/services/plugin_system.go`):
- Go plugin system with hot-reloading
- Event-driven inter-plugin communication
- Sandboxing with resource limits
- Plugin marketplace support
- Lifecycle management (install/enable/disable/uninstall)

### GDPR Compliance

**Privacy Features** (`backend/services/privacy_gdpr.go`):
- Right to erasure (Article 17)
- Data portability (Article 20)
- Consent management with granular controls
- Data Subject Access Requests (DSAR)
- Data breach notification system
- Audit trails and compliance reporting

## Important Development Notes

### Environment Configuration

Always copy `.env.example` to `.env` before running:
```bash
cp .env.example backend/.env      # For backend
cp .env.example frontend/.env.local   # For frontend
```

Key variables:
- `DB_TYPE`: "sqlite" (default) or "postgres"
- `DB_PATH`: SQLite database path (when using SQLite)
- `JWT_SECRET`: Change in production
- `OLLAMA_HOST`: For local AI models
- `ENABLE_LOCAL_AI`: true/false for local AI features

### Database Migrations

GORM handles auto-migration on startup. Models are defined in `backend/models/*.go` with GORM tags. When adding new models:
1. Define the model struct with GORM tags
2. Add to auto-migration in `backend/config/database.go`
3. GORM will create tables and columns automatically

### API Route Patterns

All API routes follow `/api/v1/{resource}` pattern:
- Public: `/api/v1/auth/*`
- Protected: Require JWT token
- Admin: `/api/v1/admin/*` require admin role

WebSocket endpoint: `/api/v1/ws` (protected, requires JWT in query param or header)

### Error Handling

Use Gin's JSON response helpers:
```go
c.JSON(http.StatusBadRequest, gin.H{"error": "descriptive message"})
c.JSON(http.StatusOK, gin.H{"data": result})
```

Services should return `(data, error)` pattern, routes convert to HTTP responses.

### Adding New Features

1. Define model in `backend/models/`
2. Create service in `backend/services/`
3. Create route handlers in `backend/routes/`
4. Register routes in `backend/main.go`
5. Create Svelte components in `frontend/src/lib/components/`
6. Add route pages in `frontend/src/routes/`

### Testing Approach

Integration tests in `backend/tests/integration_test.go` test full request/response cycles with test database. When adding tests:
- Use `setupTestDB()` to create isolated test database
- Clean up with `teardownTestDB()` after tests
- Test happy path and error cases
- Test authentication and authorization

## Common Pitfalls

1. **Variable Initialization Order in main.go**: Services must be initialized before being used. Watch for the order - cacheService must exist before authService initialization.

2. **WebSocket Authentication**: WebSocket upgrade happens after HTTP auth middleware, so user_id is available in Gin context.

3. **Database Type**: Code must handle both SQLite and PostgreSQL. Use GORM abstraction, avoid database-specific SQL.

4. **Encryption**: All sensitive data must use the encryption utilities. Check `models/encrypted.go` for the interface.

5. **CORS in Development**: Frontend (port 3000) connects to backend (port 8080). CORS middleware allows this in development but should be restricted in production.

6. **Redis Optional**: Code must handle Redis being unavailable. CacheService can be nil, services should check before using.

## Project Status

See TODO.md for detailed completion status. Major completed features:
- Complete office suite (text editor, spreadsheet, forms)
- Email client with PGP encryption
- Calendar with sharing and notifications
- Real-time chat and video conferencing
- End-to-end encryption system
- File synchronization
- AI integration with Ollama
- GDPR compliance implementation
- Plugin system
- Admin console

Next phases focus on mobile apps, enhanced collaboration, and performance optimization.
