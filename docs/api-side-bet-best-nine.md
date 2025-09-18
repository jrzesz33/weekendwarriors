# Best Nine Side Bet API

## Overview

The Best Nine side bet calculates each player's score using their best 9 holes compared to par, with handicap adjustments applied. This creates a more forgiving scoring system where poor holes don't count against the final bet result.

## Calculation Rules

### Best Nine Logic
1. Player completes all 18 holes
2. System identifies the 9 holes with the **worst** scores relative to par
3. These worst 9 holes are **discarded**
4. The remaining 9 holes (best 9) are used for the final score
5. Handicap is applied proportionally (50% of full handicap for 9 holes)

### Example Calculation
- Player shoots: `+1, +2, +3, E, -1, +4, +2, +1, +5, +2, +1, +3, E, +2, +1, +4, +2, +1`
- Worst 9 holes: `+5, +4, +4, +3, +3, +2, +2, +2, +2` (discarded)
- Best 9 holes: `+1, +1, +1, +1, +1, E, E, -1` = `+4`
- With 18 handicap: 9 strokes applied proportionally = `+4 - 9 = -5`

## Endpoints

### Get Best Nine Standings

```http
GET /api/games/{gameId}/side-bets/best-nine
```

**Response (200 OK):**
```json
{
  "bet_type": "best_nine",
  "status": "in_progress",
  "handicap_enabled": true,
  "standings": [
    {
      "position": 1,
      "player": {
        "id": "player_123",
        "name": "John Doe",
        "handicap": 18
      },
      "best_nine_score": "-2",
      "holes_completed": 18,
      "best_holes": [1, 2, 4, 7, 9, 11, 14, 16, 18],
      "worst_holes": [3, 5, 6, 8, 10, 12, 13, 15, 17],
      "raw_best_nine": "+7",
      "handicap_adjustment": "-9",
      "final_score": "-2"
    },
    {
      "position": 2,
      "player": {
        "id": "player_456",
        "name": "Jane Smith",
        "handicap": 12
      },
      "best_nine_score": "E",
      "holes_completed": 18,
      "best_holes": [1, 3, 4, 6, 8, 10, 13, 15, 17],
      "worst_holes": [2, 5, 7, 9, 11, 12, 14, 16, 18],
      "raw_best_nine": "+6",
      "handicap_adjustment": "-6",
      "final_score": "E"
    }
  ],
  "winner": {
    "player_id": "player_123",
    "score": "-2",
    "margin": "2 strokes"
  }
}
```

### Get Player's Best Nine Details

```http
GET /api/games/{gameId}/players/{playerId}/side-bets/best-nine
```

**Response (200 OK):**
```json
{
  "player": {
    "id": "player_123",
    "name": "John Doe",
    "handicap": 18
  },
  "best_nine_analysis": {
    "holes_completed": 18,
    "status": "completed",
    "best_holes": [
      {
        "hole": 1,
        "par": 4,
        "strokes": 4,
        "score_to_par": "E",
        "handicap_stroke": true,
        "effective_score": "-1"
      },
      {
        "hole": 2,
        "par": 3,
        "strokes": 3,
        "score_to_par": "E",
        "handicap_stroke": false,
        "effective_score": "E"
      }
    ],
    "worst_holes": [
      {
        "hole": 3,
        "par": 5,
        "strokes": 8,
        "score_to_par": "+3",
        "handicap_stroke": true,
        "effective_score": "+2"
      }
    ],
    "calculations": {
      "raw_best_nine_total": "+7",
      "handicap_strokes_applied": 9,
      "final_best_nine_score": "-2"
    },
    "current_position": 1,
    "lead_margin": "+2 strokes"
  }
}
```

### Get Live Best Nine Updates

```http
GET /api/games/{gameId}/side-bets/best-nine/live
```

**Response (200 OK):**
```json
{
  "current_hole": 12,
  "projections": [
    {
      "player": {
        "id": "player_123",
        "name": "John Doe"
      },
      "holes_completed": 12,
      "current_best_nine": "+3",
      "projected_final": "-1",
      "position": 1,
      "confidence": "high"
    },
    {
      "player": {
        "id": "player_456",
        "name": "Jane Smith"
      },
      "holes_completed": 12,
      "current_best_nine": "+5",
      "projected_final": "+1",
      "position": 2,
      "confidence": "medium"
    }
  ],
  "hole_impact": {
    "hole": 12,
    "players_completed": 2,
    "significant_changes": [
      {
        "player_id": "player_123",
        "previous_position": 2,
        "new_position": 1,
        "score_change": "-2"
      }
    ]
  }
}
```

### Calculate Hypothetical Best Nine

```http
POST /api/games/{gameId}/players/{playerId}/side-bets/best-nine/calculate
```

**Request Body:**
```json
{
  "hypothetical_scores": {
    "remaining_holes": [
      {
        "hole": 13,
        "projected_strokes": 4
      },
      {
        "hole": 14,
        "projected_strokes": 5
      }
    ]
  }
}
```

**Response (200 OK):**
```json
{
  "current_best_nine": "+3",
  "projected_best_nine": "-1",
  "scenarios": {
    "best_case": "-3",
    "worst_case": "+2",
    "most_likely": "-1"
  },
  "holes_analysis": {
    "completed": 12,
    "remaining": 6,
    "current_worst_holes": [3, 5, 8, 10, 11, 12],
    "holes_at_risk": [6, 7]
  }
}
```

## Handicap Application

### Full Handicap (18 holes)
- Players receive strokes on holes based on hole handicap ranking
- Hole handicap 1-18 determines stroke allocation priority

### Best Nine Handicap (9 holes)
- 50% of player's handicap rounded to nearest whole number
- Applied proportionally across the 9 best holes
- Example: 18 handicap = 9 strokes for best nine

### Stroke Allocation Priority
1. Strokes applied to holes with handicap ranking ≤ player's reduced handicap
2. If player has 9 strokes, they get strokes on handicap holes 1-9
3. Applied only to the best 9 holes used in final calculation

## Game States

### In Progress
- Players still completing holes
- Live projections available
- Standings may change significantly

### Completed
- All players finished 18 holes
- Final best nine calculated
- Winner determined
- Payouts can be calculated

## Error Responses

### Insufficient Holes (400)
```json
{
  "error": "insufficient_holes",
  "message": "Best Nine calculation requires at least 9 completed holes",
  "details": {
    "holes_completed": 6,
    "holes_required": 9
  }
}
```

### Best Nine Not Enabled (400)
```json
{
  "error": "side_bet_not_enabled",
  "message": "Best Nine side bet is not enabled for this game"
}
```

### Game Not Completed (400)
```json
{
  "error": "game_not_completed",
  "message": "Final Best Nine results available only after game completion",
  "details": {
    "current_status": "in_progress",
    "completion_percentage": 67
  }
}
```

## WebSocket Updates

Subscribe to real-time Best Nine updates:

```javascript
// WebSocket endpoint
ws://api.example.com/ws/games/{gameId}/side-bets/best-nine

// Example update message
{
  "type": "best_nine_update",
  "game_id": "game_abc123",
  "hole": 8,
  "player_id": "player_123",
  "new_score": "+2",
  "position_change": {
    "previous": 2,
    "current": 1
  },
  "standings": [...]
}
```