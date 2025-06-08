# Slack Invite Manager

A Go-based application for managing Slack channel invites with a React frontend.

## Project Structure

```
.
├── backend/              # Backend Go application
│   ├── cmd/             # Main applications
│   │   └── server/      # Main server application
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

## Environment Variables

Create a `.env` file in the root directory with the following variables:

```
SLACK_TOKEN=your-slack-token
SLACK_CHANNEL_ID=your-channel-id
DATABASE_PATH=/app/data/slack-invite.db
```

## Development

1. Start the development environment:
   ```bash
   docker-compose up
   ```

2. The application will be available at:
   - Backend: http://localhost:8080
   - Frontend: http://localhost:3000

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