# Golf Gamez API Architecture

## Overview

The Golf Gamez API is designed to support weekend warrior golf game tracking with live score updates and side bet calculations. The API follows RESTful principles and provides endpoints for game management, score tracking, and real-time side bet calculations.

## Core Entities

### Game
- Represents a golf round session
- Contains metadata like course information, game settings
- Manages player list and game state
- Generates shareable links for spectators

### Player
- Individual golfer in a game
- Contains name, handicap, and performance statistics
- Tracks cards for Putt Putt Poker side bet

### Score
- Individual hole performance for a player
- Includes stroke count, putt count, and calculated scores
- Used for side bet calculations

### Side Bets
- **Best Nine**: Calculates best 9 holes vs par with handicap
- **Putt Putt Poker**: Card-based betting game based on putting performance

## API Design Principles

1. **Stateless**: Each request contains all necessary information
2. **RESTful**: Standard HTTP methods and status codes
3. **Real-time**: WebSocket support for live updates
4. **Anonymous**: No user authentication required
5. **Shareable**: Public game links for spectators

## Course Data

The API uses Diamond Run golf course data:
- 18 holes with predefined par values
- Par 3: Holes 2, 5, 8, 11, 13
- Par 4: Holes 1, 4, 6, 7, 10, 14, 15, 16, 18
- Par 5: Holes 3, 9, 12, 17
- Total Par: 72

## Technology Stack Considerations

- **Backend**: Go with Chi router for performance
- **Database**: SQLite/PostgreSQL for persistence
- **Real-time**: WebSockets for live updates
- **Frontend**: Go templates with HTMX or React SPA
- **Deployment**: Docker containers with simple deployment