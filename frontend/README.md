# Golf Gamez Frontend PWA

A Progressive Web Application built with Go and WebAssembly for tracking golf scores and side bets. This mobile-first application provides real-time score tracking, anonymous game creation, and offline capabilities.

## Features

### Core Functionality
- **Anonymous Game Creation**: No registration required, instant game setup
- **Real-time Updates**: WebSocket-powered live score and leaderboard updates
- **Mobile-First Design**: Touch-optimized UI for outdoor use
- **Offline Support**: PWA capabilities with service worker caching
- **Responsive Layout**: Works seamlessly on phones, tablets, and desktops

### Golf Features
- **Diamond Run Golf Course**: Pre-configured 18-hole course (Par 71)
- **Scorecard Interface**: Intuitive hole-by-hole score entry
- **Handicap System**: Support for player handicaps with guidance
- **Live Leaderboard**: Real-time standings and player statistics

### Side Bets
- **Best Nine**: Calculates best 9 holes vs par with handicap adjustments
- **Putt Putt Poker**: Card-based betting game based on putting performance

## Technology Stack

- **Framework**: [go-app](https://go-app.dev/) for WebAssembly compilation
- **Language**: Go 1.21+
- **Architecture**: Component-based with clean separation of concerns
- **Styling**: Mobile-first CSS with touch optimization
- **PWA Features**: Service Worker, Web App Manifest, offline caching
- **Real-time**: WebSocket integration for live updates

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Make (for build scripts)
- Backend API running on `http://localhost:8080`

### Development Setup

1. **Clone and navigate to frontend directory**:
```bash
cd /workspaces/golf_gamez/frontend
```

2. **Install dependencies**:
```bash
make deps
```

3. **Start development server**:
```bash
make dev
```

4. **Open in browser**:
   - Web: `http://localhost:8000`
   - Mobile: `http://[your-ip]:8000`

### Production Build

```bash
# Build optimized version
make build

# Serve production build
make serve-production
```

## Project Structure

```
frontend/
├── cmd/                    # Application entry points
│   ├── main.go            # Main WebAssembly application
│   ├── game.go            # Game management component
│   ├── scorecard.go       # Scorecard interface component
│   ├── leaderboard.go     # Leaderboard display component
│   └── server.go          # Development server
├── internal/              # Internal packages
│   ├── components/        # Reusable UI components
│   ├── services/          # Business logic services
│   │   ├── api.go         # Backend API client
│   │   └── websocket.go   # WebSocket service
│   ├── models/            # Data models
│   └── utils/             # Utility functions
├── web/                   # Static web assets
│   ├── css/
│   │   └── app.css        # Main stylesheet
│   ├── js/
│   │   └── app.js         # PWA enhancements
│   ├── static/            # Icons and images
│   ├── manifest.json      # PWA manifest
│   └── sw.js              # Service worker
├── build/                 # Build output (generated)
├── Makefile              # Build automation
├── Dockerfile            # Container configuration
├── docker-compose.yml    # Development environment
└── README.md             # This file
```

## Available Commands

### Development
```bash
make dev          # Build and serve with auto-reload
make build-dev    # Build development version
make watch        # Watch for changes (requires entr)
make serve        # Serve built application
```

### Production
```bash
make build        # Build production version
make optimize     # Build with optimizations
make prod         # Production build with all optimizations
```

### Quality Assurance
```bash
make test         # Run tests
make test-coverage # Run tests with coverage
make lint         # Run linter
make fmt          # Format code
make check        # Run all quality checks
```

### Utilities
```bash
make clean        # Clean build artifacts
make deps         # Install dependencies
make info         # Show build information
make size         # Show build size information
```

### Docker
```bash
make docker-build # Build Docker image
make docker-run   # Run in container
docker-compose up frontend-dev  # Development with Docker
```

## PWA Features

### Service Worker
- **Caching Strategy**: Cache-first for static assets, network-first for API
- **Offline Support**: Graceful degradation when network unavailable
- **Background Sync**: Queues score updates when offline
- **Update Management**: Automatic updates with user notification

### App Manifest
- **Installation**: Can be installed as native app on mobile devices
- **Icons**: Complete icon set for all device sizes
- **Theme**: Golf-themed green color scheme
- **Shortcuts**: Quick access to common actions

### Performance
- **WebAssembly**: Fast execution with small bundle size
- **Lazy Loading**: Components loaded on demand
- **Caching**: Aggressive caching of static assets
- **Compression**: Gzip compression for all text assets

## Mobile Optimization

### Touch Interface
- **Touch Targets**: Minimum 44px touch targets
- **Gesture Support**: Swipe navigation where appropriate
- **Feedback**: Visual feedback for all interactions
- **Accessibility**: Screen reader and keyboard navigation support

### Layout
- **Responsive Design**: Mobile-first approach with progressive enhancement
- **Viewport Handling**: Proper viewport meta tags and CSS units
- **Orientation**: Works in both portrait and landscape modes
- **Safe Areas**: Respects device safe areas and notches

## API Integration

### Authentication
- **Token-based**: Uses share tokens extracted from game creation
- **Anonymous**: No user registration required
- **Security**: Secure token handling and storage

### Real-time Features
- **WebSocket**: Live updates for scores and leaderboard
- **Reconnection**: Automatic reconnection with exponential backoff
- **Offline Queuing**: Queue updates when offline for later sync

### Error Handling
- **Graceful Degradation**: App continues to work with limited functionality
- **User Feedback**: Clear error messages and recovery suggestions
- **Retry Logic**: Automatic retries for transient failures

## Development Guidelines

### Component Architecture
- **Single Responsibility**: Each component has a clear, focused purpose
- **State Management**: Minimal, local state with clear data flow
- **Event Handling**: Consistent event handling patterns
- **Reusability**: Components designed for reuse across the application

### Code Style
- **Go Standards**: Follows standard Go conventions and best practices
- **Documentation**: Comprehensive comments for public APIs
- **Error Handling**: Explicit error handling throughout
- **Testing**: Unit tests for business logic

### Performance
- **Bundle Size**: Optimized WebAssembly output
- **Memory Usage**: Efficient memory management
- **Network**: Minimized API calls and optimized caching
- **Rendering**: Efficient DOM updates

## Deployment

### Static Hosting
The application builds to static files that can be hosted on any web server:

```bash
# Build for production
make build

# Deploy files from build/ directory
```

### Docker Deployment
```bash
# Build production container
docker build -t golf-gamez-frontend .

# Run container
docker run -p 80:80 golf-gamez-frontend
```

### Platform-Specific Deployment
```bash
# Netlify
make deploy-netlify

# Vercel
make deploy-vercel

# Create deployment package
make deploy-prepare
```

## Browser Support

- **Modern Browsers**: Chrome 57+, Firefox 52+, Safari 11+, Edge 16+
- **WebAssembly**: Required for application execution
- **Service Workers**: Required for offline functionality
- **ES6**: Modern JavaScript features used throughout

## Performance Metrics

- **First Contentful Paint**: < 1.5s on 3G
- **Time to Interactive**: < 3s on 3G
- **Bundle Size**: < 2MB total (including WebAssembly)
- **Lighthouse Score**: 90+ across all categories

## Contributing

1. **Development Setup**: Follow the Quick Start guide
2. **Code Standards**: Run `make check` before committing
3. **Testing**: Add tests for new functionality
4. **Documentation**: Update README for significant changes

## Troubleshooting

### Common Issues

**Build Fails**:
```bash
# Ensure Go version is 1.21+
go version

# Clean and rebuild
make clean && make build
```

**Server Won't Start**:
```bash
# Check if port is in use
lsof -i :8000

# Try different port
make serve -port=8001
```

**WebAssembly Loading Issues**:
```bash
# Ensure wasm_exec.js is present
make ensure-wasm-exec

# Check browser console for MIME type errors
```

**API Connection Issues**:
- Ensure backend is running on `http://localhost:8080`
- Check CORS settings in development
- Verify network connectivity

### Debug Mode
Enable debug logging by setting environment variables:
```bash
GOLFGAMEZ_DEBUG=true make dev
```

## License

This project is part of the Golf Gamez application suite. See the main project README for license information.

## Support

For issues and questions:
- Check the troubleshooting section above
- Review the main project documentation
- Open an issue on the project repository

---

Built with ❤️ for weekend golfers who love their side bets!