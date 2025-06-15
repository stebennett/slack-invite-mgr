# Invite Manager

A Go-based application for managing invites with a React frontend.

## Project Structure

```
.
├── backend/              # Backend Go application
│   ├── cmd/             # Main applications
│   │   ├── api/        # Main API server
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
- `EMAIL_RECIPIENT`: Email address to receive notifications
- `SMTP2GO_FROM_EMAIL`: Your verified sender email address
- `SMTP2GO_USERNAME`: Your SMTP2Go username
- `SMTP2GO_PASSWORD`: Your SMTP2Go API key
- `GITHUB_USERNAME`: Your GitHub username (for container registry)

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
```bash
# Start the sheets service using pre-built image
docker compose -f docker-compose.sheets.yml up -d
```

## Docker Images

The application uses three Docker images from GitHub Container Registry:
- `ghcr.io/<username>/slack-invite-mgr-backend`: Backend API service
- `ghcr.io/<username>/slack-invite-mgr-web`: Frontend web service
- `ghcr.io/<username>/slack-invite-mgr-sheets`: Google Sheets integration service

These images are automatically built and published by GitHub Actions when changes are pushed to the main branch.

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

### Development Environment
- Uses local builds with hot-reloading
- Frontend runs on port 3000
- Source code is mounted for live updates
- Uses development-specific configurations

### Production Environment
- Uses pre-built images from GitHub Container Registry
- Frontend runs on port 80
- Optimized for production use
- Separate compose files for different services
