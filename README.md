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
├── docker-compose.yml           # Development Docker compose
├── docker-compose.app.yml       # Production app compose
├── docker-compose.sheets.yml    # Production sheets compose
└── README.md         # This file
```

## Prerequisites

- Go 1.22+
- Node.js 20+
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
- `GITHUB_USERNAME`: Your GitHub username (for container registry)

Optional environment variables:
- `GOOGLE_TOKEN_FILE`: Path to OAuth2 token file (if using user flow instead of service account)

Example:
```bash
export GOOGLE_CREDENTIALS_FILE="path/to/credentials.json"
export GOOGLE_SPREADSHEET_ID="your-spreadsheet-id"
export GOOGLE_SHEET_NAME="Sheet1"
export EMAIL_RECIPIENT="notifications@example.com"
export SMTP2GO_FROM_EMAIL="your.email@yourdomain.com"
export SMTP2GO_USERNAME="your-smtp2go-username"
export SMTP2GO_PASSWORD="your-smtp2go-api-key"
export GITHUB_USERNAME="your-github-username"
```

## Development

1. Clone the repository:
   ```bash
   git clone https://github.com/stebennett/invite-manager.git
   cd invite-manager
   ```

2. Set up environment variables as described above

3. Start the development environment:
   ```bash
   docker compose up
   ```

4. The application will be available at:
   - Frontend: http://localhost:3000
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
  - Frontend tests: `npm test -- --passWithNoTests`
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
npm test
```

## Development vs Production

### Development Environment (`docker-compose.yml`)
- Builds images locally from source code
- Frontend runs on port 3000 with hot-reloading
- Backend runs on port 8080
- Source code is mounted as volumes for live updates
- Uses `npm start` for React development server
- Backend runs with `go run cmd/server/main.go`

### Production Environment
- Uses pre-built images from GitHub Container Registry
- `docker-compose.app.yml`: Main application (API + Web frontend)
  - Frontend runs on port 80 with Nginx
  - Backend API runs on port 8080
  - Optimized production builds
- `docker-compose.sheets.yml`: Sheets integration service
  - Runs as a standalone service
  - Typically executed on-demand or via cron job
  - Sends email notifications when complete

**Note on Frontend Deployment:**
The web frontend supports fully configurable subpath deployment at **runtime** via the `PUBLIC_URL` environment variable. No rebuild is required to change the deployment path.

**How it works:**
- The Docker image is built with relative asset paths, allowing deployment at any subpath
- At container startup, `PUBLIC_URL` is injected into:
  - `config.js` - runtime configuration for React API calls
  - `nginx.conf` - proxy routes for the backend API
- This allows serving from any subdirectory (e.g., `https://example.com/my-app/`) without rebuilding

**To deploy at a custom subpath:**

Edit `docker-compose.app.yml` and set `PUBLIC_URL` to your desired path:
```yaml
web:
  environment:
    - API_URL=http://app:8080
    - PUBLIC_URL=/your-custom-path  # Change this to your desired subpath
```

Then start (or restart) the containers:
```bash
docker compose -f docker-compose.app.yml up -d
```

**To deploy at root path:**

Set `PUBLIC_URL` to an empty string:
```yaml
web:
  environment:
    - API_URL=http://app:8080
    - PUBLIC_URL=  # Empty for root path deployment
```
