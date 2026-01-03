# Invite Manager

A Go-based application for managing invites with a React frontend.

## Project Structure

```
.
├── backend/              # Backend Go application
│   ├── cmd/             # Main applications
│   │   ├── server/     # Main API server
│   │   └── sheets/     # Google Sheets integration tool
│   ├── internal/        # Private application code
│   │   ├── api/        # API handlers and routes
│   │   ├── config/     # Configuration management
│   │   ├── models/     # Data models
│   │   └── services/   # Business logic
│   ├── pkg/            # Public library code
│   └── test/           # Additional test files
├── web/                # Frontend React application
│   ├── src/           # React source code
│   └── public/        # Static assets
├── .github/           # GitHub Actions workflows
├── docker-compose.yml           # Production Docker compose
├── docker-compose.dev.yml       # Development Docker compose (hot-reload)
├── docker-compose.app.yml       # Production app compose (pre-built images)
├── docker-compose.sheets.yml    # Production sheets compose
└── README.md         # This file
```

## Prerequisites

- Go 1.22+
- Node.js 24+
- Docker and docker-compose
- Google Cloud project with Sheets API enabled
- Google service account credentials
- GitHub account (for container registry access)

## Environment Variables

Required environment variables:
- `GOOGLE_CREDENTIALS_FILE`: Path to your Google service account credentials JSON file
- `GOOGLE_SPREADSHEET_ID`: ID of your Google Spreadsheet
- `GOOGLE_SHEET_NAME`: Name of the sheet to use
- `EMAIL_RECIPIENT`: Email address to receive notifications (for sheets service)
- `SMTP2GO_FROM_EMAIL`: Your verified sender email address (for sheets service)
- `SMTP2GO_USERNAME`: Your SMTP2Go username (for sheets service)
- `SMTP2GO_PASSWORD`: Your SMTP2Go API key (for sheets service)
- `DASHBOARD_URL`: URL for the dashboard link in email notifications (for sheets service)
- `GITHUB_USERNAME`: Your GitHub username (for container registry)

Optional environment variables:
- `GOOGLE_TOKEN_FILE`: Path to OAuth2 token file (if using user flow instead of service account)
- `LOG_LEVEL`: Logging verbosity - `debug`, `info`, `warn`, `error` (default: `info`)

Example:
```bash
export GOOGLE_CREDENTIALS_FILE="path/to/credentials.json"
export GOOGLE_SPREADSHEET_ID="your-spreadsheet-id"
export GOOGLE_SHEET_NAME="Sheet1"
export EMAIL_RECIPIENT="notifications@example.com"
export SMTP2GO_FROM_EMAIL="your.email@yourdomain.com"
export SMTP2GO_USERNAME="your-smtp2go-username"
export SMTP2GO_PASSWORD="your-smtp2go-api-key"
export DASHBOARD_URL="https://your-dashboard-url.example.com"
export GITHUB_USERNAME="your-github-username"
export LOG_LEVEL="info"
```

## Logging

The application uses structured JSON logging optimized for Grafana Loki integration.

### Log Format

All logs are output as JSON with consistent fields:

```json
{
  "time": "2026-01-03T10:15:30.123Z",
  "level": "INFO",
  "app": "slack-invite-api",
  "msg": "request completed",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "GET",
  "path": "/api/invites",
  "status": 200,
  "duration": "45.2ms"
}
```

### Application Identifiers

- `slack-invite-api` - Backend API server
- `slack-invite-sheets` - Google Sheets sync service
- `slack-invite-web` - Frontend (errors sent to backend)

### Log Levels

Set via `LOG_LEVEL` environment variable:
- `debug` - Verbose debugging information
- `info` - General operational messages (default)
- `warn` - Warning conditions
- `error` - Error conditions

### Frontend Error Logging

Frontend errors are captured and sent to the backend `/api/logs` endpoint for centralized logging. This includes:
- API call failures
- Clipboard operation errors
- React component crashes (via ErrorBoundary)

## Development

1. Clone the repository:
   ```bash
   git clone https://github.com/stebennett/invite-manager.git
   cd invite-manager
   ```

2. Set up environment variables as described above

3. Start the development environment:
   ```bash
   docker compose -f docker-compose.dev.yml up
   ```

4. The application will be available at:
   - Frontend: http://localhost:5173 (Vite dev server)
   - Backend API: http://localhost:8080

## Running in Production

### Main Application (API + Web)
```bash
# Start the main application using pre-built images
docker compose -f docker-compose.app.yml up -d
```

The application will be available at:
- Frontend: http://localhost:80
- Backend API: http://localhost:8080

### Sheets Service

The sheets service is designed to run on-demand or via a cron job. You can run it using:

**Option 1: Using docker compose**
```bash
docker compose -f docker-compose.sheets.yml up
```

**Option 2: Using the provided script**
```bash
./run-sheets.sh
```

The `run-sheets.sh` script:
- Validates all required environment variables
- Checks for credentials file existence
- Copies credentials to the data directory
- Runs the sheets service container
- Shows logs in real-time

## Docker Images

The application uses three Docker images from GitHub Container Registry:
- `ghcr.io/<username>/slack-invite-mgr-backend`: Backend API service (built from `backend/Dockerfile`)
- `ghcr.io/<username>/slack-invite-mgr-web`: Frontend web service with Nginx (built from `web/Dockerfile`)
- `ghcr.io/<username>/slack-invite-mgr-sheets`: Google Sheets integration service (built from `backend/Dockerfile.sheets`)

These images are automatically built and published by GitHub Actions (`.github/workflows/ci.yml`):
- **Test Job**: Runs on all PRs and pushes
  - Backend tests: `go test -v ./...`
  - Frontend tests: `npm test` (Vitest)
- **Build Job**: Runs on push to main branch
  - Builds and pushes all three Docker images
  - Tags: `latest`, `main`, and commit SHA
- **Release Job**: Runs on version tags (v*)
  - Builds and pushes with version tags
  - Creates GitHub release with image information

## Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd web
npm test        # Run tests once
npm run test:watch  # Run tests in watch mode
```

The frontend uses Vitest with React Testing Library for component testing.

## Development vs Production

### Development Environment (`docker-compose.dev.yml`)
- Builds images locally from source code
- Frontend runs on port 5173 with Vite hot-reloading (HMR)
- Backend runs on port 8080 with Air hot-reloading
- Source code is mounted as volumes for live updates
- Uses `npm run dev` for Vite development server
- Backend uses Air for automatic rebuilds on file changes

### Production Environment
There are two options for production deployment:

**Option 1: Build locally (`docker-compose.yml`)**
- Builds images from source code
- Frontend runs on port 80 with Nginx
- Backend API runs on port 8080

**Option 2: Pre-built images (`docker-compose.app.yml`)**
- Uses pre-built images from GitHub Container Registry
- Frontend runs on port 80 with Nginx
- Backend API runs on port 8080
- Optimized production builds

**Sheets Service (`docker-compose.sheets.yml`)**
- Runs as a standalone service
- Typically executed on-demand or via cron job
- Sends email notifications when complete

**Note on Frontend Deployment:**
The web frontend uses relative asset paths, allowing deployment at any URL path without rebuilding.

**Environment variables:**
- `API_URL`: Full URL to the backend API (must be browser-accessible)

**How it works:**
- The Docker image is built with relative asset paths (Vite `base: './'`)
- At container startup, `API_URL` is injected into `config.js` for the frontend
- The frontend makes direct API calls to `API_URL` (no nginx proxy)
- This design is portable for Kubernetes deployments where Ingress handles routing

**Important:** `API_URL` must be accessible from the user's browser, not just within the container network. For docker-compose, use `http://localhost:8080` or the external hostname. For Kubernetes, use the Ingress URL or external service endpoint.

**Configuration example:**

Edit `docker-compose.app.yml` and set the API URL:
```yaml
web:
  environment:
    - API_URL=http://localhost:8080   # Browser-accessible API URL
```

Then start (or restart) the containers:
```bash
docker compose -f docker-compose.app.yml up -d
```
