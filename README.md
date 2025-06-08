# Slack Invite Manager

A Go-based application for managing Slack channel invites with a React frontend.

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
- Slack API token with appropriate permissions
- Google Cloud project with Sheets API enabled
- Google service account credentials

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
# Slack Configuration
SLACK_TOKEN=your-slack-token
SLACK_CHANNEL_ID=your-channel-id
DATABASE_PATH=/app/data/slack-invite.db

# Google Sheets Configuration
GOOGLE_CREDENTIALS_FILE=path/to/your/service-account.json
GOOGLE_SPREADSHEET_ID=your-google-sheet-id
GOOGLE_SHEET_NAME=Sheet1
```

## Development

1. Start the development environment:
   ```bash
   docker-compose up
   ```

2. The application will be available at:
   - Backend: http://localhost:8080
   - Frontend: http://localhost:3000

## Google Sheets Integration

The application includes a command-line tool for testing the Google Sheets integration. To use it:

1. Set up Google Cloud:
   - Create a project in Google Cloud Console
   - Enable the Google Sheets API
   - Create a service account and download the credentials JSON file
   - Share your Google Sheet with the service account email

2. Configure environment variables:
   ```bash
   export GOOGLE_CREDENTIALS_FILE=path/to/your/service-account.json
   export GOOGLE_SPREADSHEET_ID=your-google-sheet-id
   export GOOGLE_SHEET_NAME=Sheet1
   ```

3. Run the sheets command:
   ```bash
   cd backend
   go run cmd/sheets/main.go
   ```

The tool will read columns A-J from the specified sheet and output the data in JSON format.

## Building

1. Build the application:
   ```bash
   docker build -t slack-invite-mgr .
   ```

2. Run the application:
   ```bash
   docker run -p 8080:8080 slack-invite-mgr
   ```

## Testing

Run the backend tests:
```bash
cd backend
go test ./...
```

## License

MIT 