# Player Management API

## Overview

The Player Management API handles adding, updating, and managing players within a golf game session.

## Endpoints

### Add Player to Game

```http
POST /api/games/{gameId}/players
```

**Request Body:**
```json
{
  "name": "John Doe",
  "handicap": 18,
  "gender": "male"
}
```

**Response (201 Created):**
```json
{
  "id": "player_123abc456def",
  "name": "John Doe",
  "handicap": 18,
  "gender": "male",
  "position": 1,
  "game_id": "game_abc123def456",
  "created_at": "2025-09-18T10:45:00Z",
  "stats": {
    "holes_completed": 0,
    "current_score": 0,
    "putts": 0,
    "poker_cards": 3
  }
}
```

### Get All Players in Game

```http
GET /api/games/{gameId}/players
```

**Response (200 OK):**
```json
{
  "players": [
    {
      "id": "player_123abc456def",
      "name": "John Doe",
      "handicap": 18,
      "gender": "male",
      "position": 1,
      "stats": {
        "holes_completed": 8,
        "current_score": "+15",
        "total_putts": 18,
        "poker_cards": 4,
        "best_nine_score": "+8"
      }
    },
    {
      "id": "player_789xyz012ghi",
      "name": "Jane Smith",
      "handicap": 24,
      "gender": "female",
      "position": 2,
      "stats": {
        "holes_completed": 8,
        "current_score": "+12",
        "total_putts": 16,
        "poker_cards": 3,
        "best_nine_score": "+6"
      }
    }
  ],
  "total_count": 2
}
```

### Get Specific Player

```http
GET /api/games/{gameId}/players/{playerId}
```

**Response (200 OK):**
```json
{
  "id": "player_123abc456def",
  "name": "John Doe",
  "handicap": 18,
  "gender": "male",
  "position": 1,
  "game_id": "game_abc123def456",
  "created_at": "2025-09-18T10:45:00Z",
  "stats": {
    "holes_completed": 8,
    "current_score": "+15",
    "total_putts": 18,
    "poker_cards": 4,
    "best_nine_score": "+8"
  },
  "hole_scores": [
    {
      "hole": 1,
      "strokes": 5,
      "putts": 2,
      "par": 4,
      "score_to_par": "+1"
    }
  ]
}
```

### Update Player

```http
PUT /api/games/{gameId}/players/{playerId}
```

**Request Body:**
```json
{
  "name": "John Smith",
  "handicap": 20
}
```

**Response (200 OK):**
```json
{
  "id": "player_123abc456def",
  "name": "John Smith",
  "handicap": 20,
  "gender": "male",
  "position": 1,
  "updated_at": "2025-09-18T11:15:00Z"
}
```

### Remove Player from Game

```http
DELETE /api/games/{gameId}/players/{playerId}
```

**Response (204 No Content)**

### Get Player Leaderboard

```http
GET /api/games/{gameId}/leaderboard
```

**Response (200 OK):**
```json
{
  "overall": [
    {
      "position": 1,
      "player": {
        "id": "player_789xyz012ghi",
        "name": "Jane Smith"
      },
      "score": "+12",
      "holes_completed": 8,
      "total_putts": 16
    },
    {
      "position": 2,
      "player": {
        "id": "player_123abc456def",
        "name": "John Doe"
      },
      "score": "+15",
      "holes_completed": 8,
      "total_putts": 18
    }
  ],
  "side_bets": {
    "best_nine": [
      {
        "position": 1,
        "player": {
          "id": "player_789xyz012ghi",
          "name": "Jane Smith"
        },
        "score": "+6"
      }
    ],
    "putt_putt_poker": [
      {
        "player": {
          "id": "player_123abc456def",
          "name": "John Doe"
        },
        "cards": 4,
        "additional_bets": 2
      }
    ]
  }
}
```

## Handicap Guidelines

The API provides helpful guidance for handicap entry:

### Default Handicap Suggestions
- **Bogey Golfer (Male)**: 20 handicap
- **Bogey Golfer (Female)**: 24 handicap

### Handicap Validation
- **Range**: 0-54 (official USGA range)
- **Decimals**: Supported (e.g., 18.5)

## Player Position

Players are automatically assigned positions based on:
1. Order of joining the game
2. Can be manually reordered for tee-off sequence

## Error Responses

### Player Not Found (404)
```json
{
  "error": "player_not_found",
  "message": "Player with ID 'player_123' not found in this game"
}
```

### Game Full (400)
```json
{
  "error": "game_full",
  "message": "Maximum of 4 players allowed per game"
}
```

### Duplicate Player Name (400)
```json
{
  "error": "duplicate_player_name",
  "message": "Player with name 'John Doe' already exists in this game"
}
```

### Invalid Handicap (400)
```json
{
  "error": "invalid_handicap",
  "message": "Handicap must be between 0 and 54",
  "details": {
    "field": "handicap",
    "provided_value": -5,
    "allowed_range": "0-54"
  }
}
```

### Game Started (400)
```json
{
  "error": "game_already_started",
  "message": "Cannot modify players after game has started"
}
```