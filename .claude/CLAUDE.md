# Slack Invite Manager - Claude Code Instructions

## Project Overview

This is a Go-based Slack invite management application with a React TypeScript frontend. The application integrates with Google Sheets for data import/export and includes email notifications via SMTP2Go.

## Project Structure

```
.
├── backend/              # Backend Go application
│   ├── cmd/             # Main applications
│   │   ├── server/      # Main API server
│   │   └── sheets/      # Google Sheets integration tool
│   ├── internal/        # Private application code
│   │   ├── api/        # API handlers and routes
│   │   ├── config/     # Configuration management
│   │   ├── models/     # Data models
│   │   └── services/   # Business logic
│   ├── pkg/            # Public library code
│   ├── test/           # Additional test files
│   ├── Dockerfile      # Backend API Dockerfile
│   └── Dockerfile.sheets  # Sheets service Dockerfile
├── web/                # Frontend React application
│   ├── src/           # React source code
│   ├── public/        # Static assets
│   ├── Dockerfile     # Production web Dockerfile
│   └── Dockerfile.dev # Development web Dockerfile
├── .github/           # GitHub Actions workflows
├── data/              # Local data storage (SQLite)
├── docker-compose.yml           # Development environment
├── docker-compose.app.yml       # Production app services
├── docker-compose.sheets.yml    # Production sheets service
└── README.md         # Project documentation
```

## Technology Stack

### Backend
- **Language**: Go 1.22+
- **Database**: SQLite (local storage)
- **Key Dependencies**:
  - `google.golang.org/api` - Google Sheets API integration
  - `golang.org/x/oauth2` - OAuth2 authentication
- **Architecture**: Standard Go project layout with clean separation of concerns

### Frontend
- **Framework**: React 19+ with TypeScript 5.9+
- **Styling**: Tailwind CSS v4 (via Vite plugin)
- **Testing**: Vitest and React Testing Library
- **Build Tool**: Vite

### Infrastructure
- **Containerization**: Docker with multi-service docker-compose
- **Web Server**: Nginx (static file serving in production)
- **SSL**: Let's Encrypt
- **CI/CD**: GitHub Actions with automated testing and deployment
- **Email**: SMTP2Go for notifications

## Coding Standards

### Go Standards
- Use Go 1.22+ features and idioms
- Follow `gofmt`, `golint`, and `govet` guidelines
- Use explicit error handling (no silent ignores)
- Implement interfaces where appropriate
- Use dependency injection for better testability
- Keep functions small and focused
- Use meaningful variable and function names
- Document all exported functions and types
- Use table-driven tests following AAA (Arrange-Act-Assert) pattern
- Mock external dependencies in tests

### React/TypeScript Standards
- Use functional components with hooks
- Use TypeScript for all components with proper prop typing
- Keep components small and focused
- Implement error boundaries where appropriate
- Use React Testing Library for component tests
- Follow best practices for state management
- Prefer composition over inheritance
- Use Tailwind CSS for styling (utility-first approach)

### Testing Standards
- Write unit tests for all business logic
- Use table-driven tests in Go
- Follow AAA (Arrange-Act-Assert) pattern
- Mock external dependencies
- Aim for high test coverage of critical paths
- Use `testing` package for Go tests
- Use Vitest and React Testing Library for frontend tests

**Run tests**:
```bash
# Backend tests
cd backend && go test ./...

# Frontend tests
cd web && npm test
```

### Security Standards
- Never commit secrets or credentials
- Use environment variables for configuration
- Sanitize all user input
- Implement proper authentication and authorization
- Follow OWASP security guidelines
- Use HTTPS for all communications
- Implement rate limiting
- Use prepared statements for database queries

## Development Workflow

### Local Development Setup

1. **Prerequisites**:
   - Go 1.22+
   - Node.js 24+
   - Docker and docker-compose
   - Google Cloud project with Sheets API enabled
   - Google service account credentials

2. **Environment Variables** (create `.env` file):
   ```bash
   GOOGLE_CREDENTIALS_FILE="path/to/credentials.json"
   GOOGLE_SPREADSHEET_ID="your-spreadsheet-id"
   GOOGLE_SHEET_NAME="Sheet1"
   EMAIL_RECIPIENT="notifications@example.com"
   SMTP2GO_FROM_EMAIL="your.email@yourdomain.com"
   SMTP2GO_USERNAME="your-smtp2go-username"
   SMTP2GO_PASSWORD="your-smtp2go-api-key"
   GITHUB_USERNAME="your-github-username"
   ```

3. **Start Development Environment**:
   ```bash
   docker compose up
   ```
   - Frontend: http://localhost:3000 (with hot-reloading)
   - Backend API: http://localhost:8080

### Production Deployment

- **Application runs on Ubuntu home server**
- **Docker containers for all services**
- Uses pre-built images from GitHub Container Registry:
  - `ghcr.io/<username>/slack-invite-mgr-backend` - Backend API
  - `ghcr.io/<username>/slack-invite-mgr-web` - Frontend web
  - `ghcr.io/<username>/slack-invite-mgr-sheets` - Sheets service

**Start production services**:
```bash
# Main application (API + Web)
docker compose -f docker-compose.app.yml up -d

# Sheets service
docker compose -f docker-compose.sheets.yml up -d
```

### CI/CD Pipeline

- GitHub Actions workflows handle CI/CD
- Automated testing on pull requests
- Automated building and publishing of Docker images
- Images pushed to GitHub Container Registry
- Deployment to home server (requires manual trigger or automated webhook)

## Key Configuration Files

- `.env` / `.env.example` - Environment variable templates
- `.github/workflows/` - CI/CD pipeline definitions
- `docker-compose.yml` - Development environment orchestration
- `docker-compose.app.yml` - Production app services
- `docker-compose.sheets.yml` - Production sheets service
- `backend/Dockerfile` - Backend API container definition
- `backend/Dockerfile.sheets` - Sheets service container definition
- `web/Dockerfile` - Frontend production container
- `web/Dockerfile.dev` - Frontend development container
- `web/nginx.conf.template` - Nginx configuration template (static file serving only)
- `web/public/config.js` - Runtime configuration template (generated at container startup with `PUBLIC_URL` and `API_URL`)
- `backend/go.mod` - Go dependencies
- `web/package.json` - Node.js dependencies

## Common Development Tasks

### Adding New API Endpoints
1. Define handler in `backend/internal/api/handlers.go`
2. Add route in `backend/internal/api/router.go`
3. Implement business logic in `backend/internal/services/`
4. Add tests in `*_test.go` files
5. Update API documentation

### Adding New Frontend Components
1. Create component in `web/src/components/`
2. Use TypeScript with proper prop types
3. Apply Tailwind CSS for styling
4. Add tests using React Testing Library
5. Import and use in parent components

### Working with Google Sheets
- Service code is in `backend/internal/services/sheets.go`
- Configuration in `backend/internal/config/sheets.go`
- Separate sheets integration tool in `backend/cmd/sheets/`

### Database Changes
- SQLite database stored in `data/` directory
- Schema changes should be handled carefully
- Regular backups are automated in production

## Monitoring and Maintenance

- Application logs are collected and monitored
- Error tracking and alerting configured
- Regular security updates required
- Automated database backups in production
- Cron jobs for maintenance tasks (sheets sync, etc.)

## Important Notes

- **Hot reloading** is enabled for both Go and React in development
- **Source code is mounted** in development containers for live updates
- **Production uses optimized builds** from GitHub Container Registry
- **Frontend serves on port 80** in production, port 3000 in development
- **Backend API** always on port 8080
- **Database** is SQLite for both dev and production (in `data/` directory)
- **Subpath deployment** is fully configurable at runtime via `PUBLIC_URL` environment variable - no rebuild required
- **API URL** is configured via `API_URL` environment variable - must be browser-accessible (not internal docker hostname)
- **Frontend makes direct API calls** - nginx does not proxy API requests, making the container portable for k8s deployments

## Getting Help

- See `README.md` for detailed setup instructions
- Check `.github/workflows/` for CI/CD pipeline details
- Review existing tests for examples of testing patterns
- Consult Go and React documentation for framework-specific questions
