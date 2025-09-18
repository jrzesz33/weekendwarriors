# Golf Gamez API

A comprehensive Go backend API service for tracking golf game scores and side bets for weekend golfers. Features real-time updates, anonymous game creation, and support for Best Nine and Putt Putt Poker side bets.

## Features

### Core Functionality
- **Anonymous Game Creation**: No registration required, shareable game links
- **Real-time Score Tracking**: Live updates via WebSocket connections
- **Multi-player Support**: Up to 4 players per game
- **Diamond Run Golf Course**: Pre-configured 18-hole course data
- **Token-based Access**: Separate share and spectator tokens for security

### Side Bets
- **Best Nine**: Calculate best 9 holes vs par with handicap adjustments
- **Putt Putt Poker**: Card-based poker betting based on putting performance

### API Features
- **RESTful Design**: Standard HTTP methods and status codes
- **Rate Limiting**: Per-IP limits for game creation and API calls
- **Structured Logging**: Comprehensive request/response logging
- **Error Handling**: Standardized error responses with request IDs
- **Health Monitoring**: Health check endpoints for monitoring

## Architecture

### Technology Stack
- **Language**: Go 1.21+
- **Database**: SQLite with WAL mode
- **Router**: Chi v5 for HTTP routing
- **WebSockets**: Gorilla WebSocket for real-time updates
- **Logging**: Zerolog for structured logging
- **Authentication**: Token-based anonymous access

### Project Structure
```
golf-gamez/
├── cmd/api/           # Application entry point
├── internal/
│   ├── config/        # Configuration management
│   ├── database/      # Database connection and migrations
│   ├── handlers/      # HTTP request handlers
│   ├── middleware/    # HTTP middleware components
│   ├── models/        # Data models and structs
│   └── services/      # Business logic services
├── pkg/
│   ├── auth/          # Authentication utilities
│   ├── errors/        # Error handling
│   └── validation/    # Input validation
├── migrations/        # Database migration files
├── docs/             # API documentation
└── data/             # SQLite database storage
```

## Quick Start

### Prerequisites
- Go 1.21 or higher
- SQLite3

### Installation
```bash
# Clone the repository
git clone <repository-url>
cd golf-gamez

# Install dependencies
go mod tidy

# Build the application
go build -o bin/golf-gamez ./cmd/api
```

### Running the Server
```bash
# Start the server (default port 8080)
./bin/golf-gamez

# Or with custom configuration
PORT=9000 DATABASE_URL=custom.db ./bin/golf-gamez
```

### Environment Variables
```bash
PORT=8080                    # Server port (default: 8080)
DATABASE_URL=data/golf_gamez.db  # SQLite database path
ENVIRONMENT=development      # Environment (development/production)
LOG_LEVEL=info              # Log level (debug/info/warn/error)
CORS_ORIGINS=*              # Allowed CORS origins
```

## API Usage

### Create a Game
```bash
curl -X POST http://localhost:8080/v1/games \
  -H "Content-Type: application/json" \
  -d '{
    "course": "diamond-run",
    "side_bets": ["best-nine", "putt-putt-poker"],
    "handicap_enabled": true
  }'
```

### Add Players
```bash
curl -X POST http://localhost:8080/v1/games/{shareToken}/players \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "handicap": 18.5,
    "gender": "male"
  }'
```

### Start Game
```bash
curl -X POST http://localhost:8080/v1/games/{shareToken}/start
```

### Record Scores
```bash
curl -X POST http://localhost:8080/v1/games/{shareToken}/players/{playerId}/scores \
  -H "Content-Type: application/json" \
  -d '{
    "hole": 1,
    "strokes": 5,
    "putts": 2
  }'
```

### Get Leaderboard
```bash
curl http://localhost:8080/v1/games/{shareToken}/leaderboard
```

### Spectator Access
```bash
curl http://localhost:8080/v1/spectate/{spectatorToken}
```

## Authentication

Golf Gamez uses a token-based authentication system without requiring user registration:

### Token Types
- **Share Token** (`gt_*`): Full read/write access for players
- **Spectator Token** (`st_*`): Read-only access for viewing

### Usage
- Tokens are generated when creating a game
- Share tokens allow game modification and score recording
- Spectator tokens provide read-only access to game data
- Tokens can be used in URL paths or Authorization headers

## Database Schema

The application uses SQLite with the following key tables:
- `games`: Game sessions and metadata
- `players`: Player information within games
- `scores`: Individual hole scores
- `course_data`: Golf course hole information
- `side_bet_calculations`: Side bet calculations and standings
- `putt_putt_poker_cards`: Poker card tracking
- `poker_hands`: Final poker hand results

## Rate Limiting

The API implements several rate limiting strategies:
- **Game Creation**: 10 games per hour per IP
- **Score Updates**: 100 requests per minute per IP
- **General API**: 1000 requests per hour per IP

## Real-time Updates

WebSocket connections provide real-time updates for:
- Score changes
- Leaderboard updates
- Side bet standings
- Game state changes
- Player additions/removals

### WebSocket Connection
```javascript
const ws = new WebSocket('ws://localhost:8080/v1/ws/games/{gameId}');
ws.onmessage = function(event) {
  const update = JSON.parse(event.data);
  // Handle real-time updates
};
```

## Side Bet Logic

### Best Nine
- Uses player's best 9 holes out of 18
- Applies 50% of handicap for 9-hole calculation
- Discards worst 9 holes from final score
- Supports live projections during play

### Putt Putt Poker
- Players start with 3 cards
- Earn +1 card for one-putts
- Earn +2 cards for hole-in-ones
- $1 penalty for 3+ putts (added to pot)
- Final poker hands dealt at game completion

## API Documentation

Complete API documentation is available in the `docs/` directory:
- **OpenAPI Specification**: `docs/openapi-spec.yaml`
- **Architecture Overview**: `docs/api-architecture.md`
- **Game Management**: `docs/api-game-management.md`
- **Player Management**: `docs/api-player-management.md`
- **Score Tracking**: `docs/api-score-tracking.md`
- **Side Bet Details**: `docs/api-side-bet-*.md`
- **Data Models**: `docs/api-models-and-errors.md`
- **Security**: `docs/security-and-authentication.md`

## Development

### Building
```bash
go build -o bin/golf-gamez ./cmd/api
```

### Testing
```bash
go test ./...
```

### Database Migrations
Migrations run automatically on startup. Manual migration control:
```bash
# View current migration status
sqlite3 data/golf_gamez.db "SELECT * FROM migrations;"
```

### Logging
The application uses structured logging with zerolog:
- Request/response logging with request IDs
- Error tracking with stack traces
- Performance metrics and timing
- Security event logging

## Production Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o golf-gamez ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/golf-gamez .
CMD ["./golf-gamez"]
```

### Environment Configuration
```bash
# Production settings
ENVIRONMENT=production
LOG_LEVEL=info
PORT=8080
DATABASE_URL=/data/golf_gamez.db
CORS_ORIGINS=https://yourdomain.com
```

### Health Checks
```bash
curl http://localhost:8080/v1/health
```

## Security Considerations

- Token-based access without PII collection
- Input validation and sanitization
- Rate limiting to prevent abuse
- CORS configuration for web clients
- SQL injection prevention via prepared statements
- Graceful error handling without information disclosure

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement changes with tests
4. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:
- Check the documentation in `docs/`
- Review the OpenAPI specification
- Create an issue on GitHub

---

Built with ❤️ for weekend golfers who love their side bets!