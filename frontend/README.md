# Golf Gamez Frontend

A mobile-first Progressive Web Application (PWA) built with Go/WebAssembly for tracking golf scores and side bets in real-time.

## Features

### Core Functionality
- **Anonymous Game Creation**: Create golf games without registration
- **Real-time Score Tracking**: Live score updates via WebSocket
- **Mobile-First Design**: Optimized for mobile devices with touch-friendly interfaces
- **Side Bet Support**: Best Nine and Putt Putt Poker tracking
- **Spectator Mode**: Share-able links for read-only game viewing
- **Offline Support**: Works offline with data sync when reconnected

### Technical Features
- **Go/WebAssembly**: High-performance client-side execution
- **Progressive Web App**: Install-able on mobile devices
- **Resilient Connectivity**: Automatic reconnection with exponential backoff
- **Touch Optimized**: 44px+ touch targets, haptic feedback
- **Responsive Design**: Works on phones, tablets, and desktop

## Architecture

### Frontend Stack
- **Go 1.21+**: Compiled to WebAssembly for client-side execution
- **Vanilla CSS**: Custom mobile-first framework with golf-specific components
- **Service Worker**: PWA capabilities and offline support
- **WebSocket**: Real-time communication with resilient retry logic

### Key Components
- **API Client**: HTTP client optimized for WebAssembly environment
- **WebSocket Manager**: Connection handling with automatic reconnection
- **UI Manager**: DOM manipulation and component rendering
- **Storage Manager**: Local storage for offline capabilities
- **Score Entry**: Touch-optimized input with validation

## Project Structure

```
frontend/
├── main.go                 # Application entry point
├── build.sh               # Build script for WebAssembly
├── go.mod                 # Go module definition
├── internal/
│   ├── app/               # Application core
│   ├── api/               # API client layer
│   ├── websocket/         # WebSocket management
│   ├── ui/                # User interface components
│   ├── storage/           # Local storage handling
│   └── models/            # Data models
└── web/
    ├── index.html         # Main HTML file
    ├── sw.js              # Service worker
    └── static/
        ├── css/           # Stylesheets
        ├── js/            # JavaScript files
        └── assets/        # Images and icons
```

## Building and Running

### Prerequisites
- Go 1.21 or later
- Modern web browser with WebAssembly support

### Build for Development
```bash
# Make build script executable
chmod +x build.sh

# Build WebAssembly binary
./build.sh
```

### Build for Production
```bash
# Build with optimizations
./build.sh production
```

### Serve Locally
```bash
# Simple HTTP server (Python)
cd web && python3 -m http.server 8000

# Or using Go
cd web && go run -m http.server

# Access at http://localhost:8000
```

## API Integration

The frontend is designed to work with the Golf Gamez API. Configure the API endpoint in `internal/api/client.go`:

```go
// Development
apiClient := api.NewClient("http://localhost:8080/v1")

// Production
apiClient := api.NewClient("https://api.golfgamez.com/v1")
```

## Mobile Optimization

### Touch Targets
- Minimum 44px touch targets for accessibility
- Enhanced touch feedback with haptic vibration
- Optimized tap zones for score entry

### Performance
- WebAssembly for native-speed calculations
- Minimal JavaScript footprint
- Efficient DOM updates
- Service worker caching

### Responsive Design
- Mobile-first CSS framework
- Breakpoints: 480px, 768px, 1024px
- Touch-optimized score entry interface
- Collapsible navigation for small screens

## Offline Support

The application provides robust offline capabilities:

### Cached Resources
- Application shell (HTML, CSS, JS)
- WebAssembly binary
- Static assets and icons

### Offline Functionality
- View cached game data
- Enter scores offline (synced when reconnected)
- Access recent games
- Basic calculations without network

### Data Sync
- Automatic sync when connection restored
- Conflict resolution for score updates
- Background sync for pending operations

## Real-time Features

### WebSocket Connection
- Automatic connection on game start
- Exponential backoff retry strategy
- Heartbeat mechanism for connection health
- Queue messages during disconnection

### Live Updates
- Real-time score updates
- Leaderboard position changes
- Side bet calculations
- Connection status indicator

## Security Considerations

### Client-Side Security
- No sensitive data stored locally
- Share tokens for game access
- Read-only spectator tokens
- HTTPS enforcement in production

### Data Validation
- Client-side input validation
- Server-side validation backup
- Score range checking
- Player limit enforcement

## Browser Support

### Minimum Requirements
- WebAssembly support
- ES6+ JavaScript features
- CSS Grid and Flexbox
- Service Worker API
- Local Storage API

### Tested Browsers
- Chrome/Chromium 80+
- Safari 14+
- Firefox 78+
- Edge 80+

### Mobile Browsers
- iOS Safari 14+
- Chrome Mobile 80+
- Samsung Internet 12+

## Performance Optimizations

### WebAssembly
- Optimized Go compilation flags
- Minimal runtime overhead
- Efficient memory management
- Fast startup time

### Network
- Service worker caching
- Resource preloading
- WebSocket connection pooling
- Compression support

### UI/UX
- 60fps animations
- Smooth scrolling
- Instant feedback
- Progressive loading

## Development Guidelines

### Code Organization
- Feature-based module structure
- Clear separation of concerns
- Consistent error handling
- Comprehensive documentation

### Performance
- Minimize DOM manipulation
- Use efficient data structures
- Optimize critical rendering path
- Profile WebAssembly performance

### Accessibility
- WCAG 2.1 AA compliance
- Screen reader support
- Keyboard navigation
- High contrast support

## Contributing

### Setup Development Environment
1. Install Go 1.21+
2. Clone repository
3. Run build script
4. Start local server
5. Open browser to localhost

### Code Style
- Go fmt for Go code
- Consistent CSS methodology
- JSDoc for JavaScript
- Semantic HTML structure

### Testing
- Unit tests for core logic
- Integration tests for API client
- Manual testing on target devices
- Performance benchmarking

## License

MIT License - See LICENSE file for details.

## Support

For issues and questions:
- GitHub Issues for bug reports
- API documentation in `/docs`
- Development guide in this README