# Contributing to TPT Titan

Thank you for your interest in contributing to TPT Titan! We welcome contributions from everyone. By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### 1. Fork the Repository
Fork the TPT Titan repository on GitHub to your own account.

### 2. Clone Your Fork
```bash
git clone https://github.com/your-username/tpt-titan.git
cd tpt-titan
```

### 3. Create a Branch
Create a new branch for your contribution:
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b bugfix/issue-number
```

### 4. Set Up Development Environment

#### Backend (Go)
```bash
cd backend
go mod download
go run main.go
```

#### Frontend (Svelte)
```bash
cd frontend
npm install
npm run dev
```

#### Desktop (Tauri)
```bash
cd desktop
npm install
npm run tauri dev
```

#### Docker Development
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f
```

### 5. Make Changes
- Follow the existing code style
- Write clear, concise commit messages
- Add tests for new features
- Update documentation as needed

### 6. Test Your Changes
Run the appropriate tests for your changes:

#### Backend Tests
```bash
cd backend
go test ./...
```

#### Frontend Tests
```bash
cd frontend
npm run test
```

#### End-to-End Tests
```bash
npm run test:e2e
```

### 7. Commit Your Changes
```bash
git add .
git commit -m "feat: add your feature description"
```

Follow conventional commit format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation
- `style:` for formatting
- `refactor:` for code restructuring
- `test:` for adding tests
- `chore:` for maintenance

### 8. Push and Create Pull Request
```bash
git push origin feature/your-feature-name
```
Then create a pull request on GitHub.

## Development Guidelines

### Code Style
- **Go**: Follow standard Go conventions and use `gofmt`
- **JavaScript/TypeScript**: Use ESLint and Prettier
- **Rust**: Follow standard Rust formatting with `rustfmt`
- **SQL**: Use consistent naming and formatting

### Testing
- Write unit tests for all new functionality
- Aim for good test coverage (>80%)
- Include integration tests for API endpoints
- Add E2E tests for critical user flows

### Documentation
- Update README.md for significant changes
- Add JSDoc/TSDoc comments for public APIs
- Update user documentation for new features

### Security
- Never commit sensitive information
- Use environment variables for secrets
- Follow security best practices
- Report security issues privately

## Project Structure

```
tpt-titan/
├── backend/          # Go backend API
├── frontend/         # Svelte frontend
├── desktop/          # Tauri desktop app
├── docs/            # Documentation
├── scripts/         # Build and deployment scripts
├── docker-compose.yml
├── LICENSE
├── README.md
└── TODO.md
```

## Communication
- Use GitHub Issues for bug reports and feature requests
- Join our Discord/Slack for discussions
- Check existing issues before creating new ones

## License
By contributing to TPT Titan, you agree that your contributions will be dual-licensed under the Apache License, Version 2.0 and the MIT License, without any additional terms or conditions.

## Recognition
Contributors will be recognized in our README.md and potentially in release notes.

Thank you for helping make TPT Titan better! 🚀
