# Game Management API

## Overview

The Game Management API handles creating, retrieving, and managing golf game sessions.

## Endpoints

### Create New Game

```http
POST /api/games
```

**Request Body:**
```json
{
  "course": "diamond-run",
  "side_bets": ["best-nine", "putt-putt-poker"],
  "handicap_enabled": true
}
```

**Response (201 Created):**
```json
{
  "id": "game_abc123def456",
  "course": "diamond-run",
  "side_bets": ["best-nine", "putt-putt-poker"],
  "handicap_enabled": true,
  "status": "setup",
  "share_link": "https://api.example.com/games/abc123def456",
  "spectator_link": "https://api.example.com/spectate/abc123def456",
  "created_at": "2025-09-18T10:30:00Z",
  "players": [],
  "current_hole": null
}
```

### Get Game Details

```http
GET /api/games/{gameId}
```

**Response (200 OK):**
```json
{
  "id": "game_abc123def456",
  "course": "diamond-run",
  "side_bets": ["best-nine", "putt-putt-poker"],
  "handicap_enabled": true,
  "status": "in_progress",
  "share_link": "https://api.example.com/games/abc123def456",
  "spectator_link": "https://api.example.com/spectate/abc123def456",
  "created_at": "2025-09-18T10:30:00Z",
  "started_at": "2025-09-18T11:00:00Z",
  "players": [
    {
      "id": "player_123",
      "name": "John Doe",
      "handicap": 18,
      "position": 1
    }
  ],
  "current_hole": 3,
  "course_info": {
    "name": "Diamond Run",
    "holes": [
      {"hole": 1, "par": 4},
      {"hole": 2, "par": 3},
      {"hole": 3, "par": 5}
    ]
  }
}
```

### Start Game

```http
POST /api/games/{gameId}/start
```

**Response (200 OK):**
```json
{
  "id": "game_abc123def456",
  "status": "in_progress",
  "started_at": "2025-09-18T11:00:00Z",
  "current_hole": 1
}
```

### End Game

```http
POST /api/games/{gameId}/complete
```

**Response (200 OK):**
```json
{
  "id": "game_abc123def456",
  "status": "completed",
  "completed_at": "2025-09-18T15:30:00Z",
  "final_results": {
    "best_nine_winner": {
      "player_id": "player_123",
      "score": "+2"
    },
    "putt_putt_poker_winner": {
      "player_id": "player_456",
      "hand": "full_house",
      "cards": ["AS", "AH", "AC", "KS", "KH"]
    }
  }
}
```

### Get Game Status

```http
GET /api/games/{gameId}/status
```

**Response (200 OK):**
```json
{
  "status": "in_progress",
  "current_hole": 8,
  "players_completed_hole": 3,
  "total_players": 4,
  "estimated_completion": "2025-09-18T15:45:00Z"
}
```

### Delete Game

```http
DELETE /api/games/{gameId}
```

**Response (204 No Content)**

## Game Status Values

- `setup`: Game created but not started
- `in_progress`: Game is active and players are recording scores
- `completed`: All 18 holes completed and final results calculated
- `abandoned`: Game was ended early or cancelled

## Error Responses

### Game Not Found (404)
```json
{
  "error": "game_not_found",
  "message": "Game with ID 'abc123def456' not found"
}
```

### Invalid Game State (400)
```json
{
  "error": "invalid_game_state",
  "message": "Cannot start game without players"
}
```

### Validation Error (400)
```json
{
  "error": "validation_error",
  "message": "Invalid side bet type",
  "details": {
    "field": "side_bets",
    "allowed_values": ["best-nine", "putt-putt-poker"]
  }
}
```