# Score Tracking API

## Overview

The Score Tracking API handles recording and retrieving golf scores for each hole, including stroke counts and putt counts for side bet calculations.

## Endpoints

### Record Score for Hole

```http
POST /api/games/{gameId}/players/{playerId}/scores
```

**Request Body:**
```json
{
  "hole": 1,
  "strokes": 5,
  "putts": 2
}
```

**Response (201 Created):**
```json
{
  "id": "score_abc123def456",
  "player_id": "player_123abc456def",
  "game_id": "game_abc123def456",
  "hole": 1,
  "strokes": 5,
  "putts": 2,
  "par": 4,
  "score_to_par": "+1",
  "handicap_adjusted": true,
  "effective_score": 4,
  "created_at": "2025-09-18T11:30:00Z",
  "side_bet_updates": {
    "putt_putt_poker": {
      "cards_awarded": 0,
      "penalty_applied": false,
      "total_cards": 3
    }
  }
}
```

### Update Score for Hole

```http
PUT /api/games/{gameId}/players/{playerId}/scores/{hole}
```

**Request Body:**
```json
{
  "strokes": 4,
  "putts": 1
}
```

**Response (200 OK):**
```json
{
  "id": "score_abc123def456",
  "player_id": "player_123abc456def",
  "game_id": "game_abc123def456",
  "hole": 1,
  "strokes": 4,
  "putts": 1,
  "par": 4,
  "score_to_par": "E",
  "handicap_adjusted": true,
  "effective_score": 3,
  "updated_at": "2025-09-18T11:35:00Z",
  "side_bet_updates": {
    "putt_putt_poker": {
      "cards_awarded": 1,
      "penalty_applied": false,
      "total_cards": 4
    }
  }
}
```

### Get Player's Scorecard

```http
GET /api/games/{gameId}/players/{playerId}/scorecard
```

**Response (200 OK):**
```json
{
  "player": {
    "id": "player_123abc456def",
    "name": "John Doe",
    "handicap": 18
  },
  "holes": [
    {
      "hole": 1,
      "par": 4,
      "strokes": 4,
      "putts": 1,
      "score_to_par": "E",
      "handicap_stroke": true,
      "effective_score": 3
    },
    {
      "hole": 2,
      "par": 3,
      "strokes": 4,
      "putts": 2,
      "score_to_par": "+1",
      "handicap_stroke": false,
      "effective_score": 4
    }
  ],
  "totals": {
    "holes_completed": 2,
    "total_strokes": 8,
    "total_putts": 3,
    "score_to_par": "+1",
    "handicap_adjusted_score": "+1"
  },
  "side_bet_status": {
    "best_nine": {
      "current_best_score": "+1",
      "holes_used": [1, 2]
    },
    "putt_putt_poker": {
      "total_cards": 4,
      "additional_bets": 0
    }
  }
}
```

### Get Game Scorecard (All Players)

```http
GET /api/games/{gameId}/scorecard
```

**Response (200 OK):**
```json
{
  "game": {
    "id": "game_abc123def456",
    "course": "diamond-run",
    "current_hole": 3
  },
  "course_info": [
    {"hole": 1, "par": 4},
    {"hole": 2, "par": 3},
    {"hole": 3, "par": 5}
  ],
  "players": [
    {
      "id": "player_123abc456def",
      "name": "John Doe",
      "position": 1,
      "scores": [
        {
          "hole": 1,
          "strokes": 4,
          "putts": 1,
          "score_to_par": "E"
        },
        {
          "hole": 2,
          "strokes": 4,
          "putts": 2,
          "score_to_par": "+1"
        }
      ],
      "totals": {
        "holes_completed": 2,
        "score_to_par": "+1",
        "total_putts": 3
      }
    }
  ]
}
```

### Get Hole-by-Hole Leaderboard

```http
GET /api/games/{gameId}/leaderboard/hole/{holeNumber}
```

**Response (200 OK):**
```json
{
  "hole": 5,
  "par": 3,
  "results": [
    {
      "position": 1,
      "player": {
        "id": "player_789xyz012ghi",
        "name": "Jane Smith"
      },
      "strokes": 2,
      "putts": 1,
      "score_to_par": "-1",
      "is_best_score": true
    },
    {
      "position": 2,
      "player": {
        "id": "player_123abc456def",
        "name": "John Doe"
      },
      "strokes": 3,
      "putts": 2,
      "score_to_par": "E",
      "is_best_score": false
    }
  ],
  "statistics": {
    "average_score": 2.5,
    "best_score": 2,
    "worst_score": 3,
    "eagles": 0,
    "birdies": 1,
    "pars": 1,
    "bogeys": 0
  }
}
```

### Delete Score

```http
DELETE /api/games/{gameId}/players/{playerId}/scores/{hole}
```

**Response (204 No Content)**

### Bulk Score Entry

```http
POST /api/games/{gameId}/players/{playerId}/scores/bulk
```

**Request Body:**
```json
{
  "scores": [
    {
      "hole": 1,
      "strokes": 5,
      "putts": 2
    },
    {
      "hole": 2,
      "strokes": 3,
      "putts": 1
    },
    {
      "hole": 3,
      "strokes": 6,
      "putts": 3
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "scores_created": 3,
  "scores": [
    {
      "hole": 1,
      "strokes": 5,
      "putts": 2,
      "score_to_par": "+1"
    }
  ],
  "updated_totals": {
    "total_strokes": 14,
    "total_putts": 6,
    "score_to_par": "+2"
  },
  "side_bet_updates": {
    "putt_putt_poker": {
      "total_cards": 4,
      "cards_awarded_this_batch": 1,
      "penalties_this_batch": 1
    }
  }
}
```

## Score Calculation Rules

### Handicap Application
- Handicap strokes are applied to specific holes based on hole difficulty
- Each hole has a handicap ranking (1-18)
- Players receive strokes on holes equal to their handicap

### Score to Par Notation
- **Eagle**: `-2` (2 under par)
- **Birdie**: `-1` (1 under par)
- **Par**: `E` (even with par)
- **Bogey**: `+1` (1 over par)
- **Double Bogey**: `+2` (2 over par)
- **Triple Bogey+**: `+3+` (3 or more over par)

### Putt Tracking Rules
- **One Putt**: Awards 1 additional card in Putt Putt Poker
- **Hole-in-One**: Awards 2 additional cards
- **Three Putt+**: Triggers penalty notification (add $1 to bet)

## Error Responses

### Invalid Hole Number (400)
```json
{
  "error": "invalid_hole",
  "message": "Hole number must be between 1 and 18",
  "details": {
    "provided_hole": 19,
    "valid_range": "1-18"
  }
}
```

### Score Already Exists (409)
```json
{
  "error": "score_already_exists",
  "message": "Score for hole 5 already recorded. Use PUT to update."
}
```

### Invalid Score Values (400)
```json
{
  "error": "invalid_score",
  "message": "Invalid stroke or putt count",
  "details": {
    "strokes": {
      "value": 0,
      "error": "Must be at least 1"
    },
    "putts": {
      "value": 5,
      "error": "Cannot exceed stroke count"
    }
  }
}
```

### Game Not Started (400)
```json
{
  "error": "game_not_started",
  "message": "Cannot record scores before game starts"
}
```

### Future Hole (400)
```json
{
  "error": "future_hole",
  "message": "Cannot record scores for future holes",
  "details": {
    "current_hole": 5,
    "requested_hole": 8
  }
}
```