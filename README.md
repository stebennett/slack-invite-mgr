# Invite Manager

A Go-based application for managing invites with a React frontend.

## Project Structure

```
.
├── backend/              # Backend Go application
│   ├── cmd/             # Main applications
│   │   ├── server/      # Main server application
│   │   └── sheets/      # Google Sheets integration tool
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
├── Dockerfile         # Main Dockerfile
├── docker-compose.yml # Docker compose configuration
└── README.md         # This file
```

## Prerequisites

- Go 1.22+
- Node.js 20+
- Docker and docker-compose
- Google Cloud project with Sheets API enabled
- Google service account credentials

## Environment Variables

Required environment variables:
- `GOOGLE_CREDENTIALS_FILE`: Path to your Google service account credentials JSON file
- `GOOGLE_SPREADSHEET_ID`: ID of your Google Spreadsheet
- `GOOGLE_SHEET_NAME`: Name of the sheet to use
- `EMAIL_RECIPIENT`: Email address to receive notifications
- `SMTP2GO_FROM_EMAIL`: Your verified sender email address
- `SMTP2GO_USERNAME`: Your SMTP2Go username
- `SMTP2GO_PASSWORD`: Your SMTP2Go API key

Example:
```bash
export GOOGLE_CREDENTIALS_FILE="path/to/credentials.json"
export GOOGLE_SPREADSHEET_ID="your-spreadsheet-id"
export GOOGLE_SHEET_NAME="Sheet1"
export EMAIL_RECIPIENT="notifications@example.com"
export SMTP2GO_FROM_EMAIL="your.email@yourdomain.com"
export SMTP2GO_USERNAME="your-smtp2go-username"
export SMTP2GO_PASSWORD="your-smtp2go-api-key"
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
   docker-compose up
   ```

4. The application will be available at:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

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

## Running

### Local Development
```bash
# Start both frontend and backend
docker-compose up

# Start only backend
docker-compose up backend

# Start only frontend
docker-compose up web
```
