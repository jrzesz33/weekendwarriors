# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Essential Commands

### Development Workflow
```bash
# Start development with auto-rebuild and serve
make dev

# Build production version
make build

# Clean rebuild (use when assets are stale)
make clean && make build

# Serve built application
make serve

# Run tests
make test

# Run linter and format code
make lint && make fmt
```

### Key Build Understanding
The build system copies files from `web/` to `build/` directory. When making CSS or static file changes, always edit the source files in `web/` directory, not `build/`, as `build/` gets overwritten on every build.

## Architecture Overview

### Go WebAssembly + go-app v10 Framework
This is a Progressive Web App built with Go that compiles to WebAssembly. It uses the go-app v10 framework for component-based UI development.

**Critical Architecture Pattern:**
- `cmd/main.go` - Main WebAssembly application with state management and routing
- `cmd/views.go` - UI rendering functions for all application states
- `cmd/handlers.go` - Event handlers for user interactions
- `internal/services/api.go` - Backend API client with authentication
- `internal/models/game.go` - Data models matching backend API
- `web/` - Source static assets (CSS, JS, manifest, service worker)
- `build/` - Generated build output (never edit directly)

### State Management Pattern
The app uses a single main component `GolfGamezApp` with:
- `currentState` enum for navigation (StateHome, StateCreateGame, StateGameActive, etc.)
- `currentGame` for active game data
- `apiClient` for backend communication
- Event handlers that update state and trigger re-renders

### Authentication & Tokens
- Uses share tokens (format: `gt_xxxxxxxx`) for anonymous game access
- Tokens extracted from share links via `apiClient.ExtractTokenFromShareLink()`
- Deep linking supports `/game/{token}` and `/spectate/{token}` URLs

### go-app v10 Specific Patterns
```go
// Component rendering
func (g *GolfGamezApp) Render() app.UI {
    return app.Div().Class("golf-app").Body(
        g.renderHeader(),
        g.renderMain(),
    )
}

// Event handling
app.Button().OnClick(g.onCreateGameStart)

// Conditional rendering
app.If(condition, func() app.UI {
    return app.Div().Text("content")
})

// Dynamic lists
app.Range(items).Slice(func(i int) app.UI {
    return app.Div().Text(items[i])
})
```

### API Integration Architecture
- `services/api.go` handles all backend communication
- Uses share tokens for authentication
- Supports game creation, player management, score recording
- WebSocket service in `services/websocket.go` for real-time updates
- All API models defined in `internal/models/game.go`

### PWA Architecture
- Service worker (`web/sw.js`) handles offline caching
- Web manifest (`web/manifest.json`) enables app installation
- CSS is mobile-first with touch optimization
- WebAssembly loads via `web/index.html` with loading screens

## Important Development Notes

### CSS/Styling Changes
Always edit `web/css/app.css` (source) not `build/css/app.css` (generated). The Makefile copies `web/*` to `build/*` on every build, overwriting any changes made directly to build files.

### go-app Framework Version
Uses go-app v10. Key differences from earlier versions:
- No `app.Dispatch()` - use direct state updates
- Event handlers have signature `func(ctx app.Context, e app.Event)`
- Conditional rendering requires lambda functions with `app.If()`

### Real-time Features
WebSocket integration is prepared but needs completion:
- Connect in `OnMount()` lifecycle
- Handle score updates, leaderboard changes
- Graceful reconnection on network issues

### Game Flow States
1. `StateHome` - Welcome screen with create/join options
2. `StateCreateGame` - Game setup form
3. `StateJoinGame` - Enter game token
4. `StateGameSetup` - Add players before starting
5. `StateGameActive` - Live scorecard and scoring
6. `StateGameResults` - Final results and side bet winners
7. `StateSpectator` - Read-only game viewing

### Testing & Quality
- Backend API should run on `localhost:8080`
- Frontend serves on `localhost:8000`
- Use `make check` for comprehensive quality checks
- Mobile testing important - app optimized for golf course use

### Side Bets Implementation
Two types supported:
- **Best Nine**: Player's best 9 holes vs par with handicap
- **Putt Putt Poker**: Card-based game using putting performance

Data structures and calculations defined in `internal/models/game.go`.