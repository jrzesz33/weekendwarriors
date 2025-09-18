# Putt Putt Poker Side Bet API

## Overview

Putt Putt Poker is a card-based side bet where players earn or lose playing cards based on their putting performance. At the end of the round, cards are randomly dealt and the best poker hand wins. Handicap does **not** apply to this side bet.

## Game Rules

### Starting Cards
- Each player begins with **3 playing cards**

### Card Awards
- **One Putt**: +1 card
- **Hole-in-One**: +2 cards
- **Three Putt or Worse**: Player must add $1 to the bet (penalty notification)

### Final Poker Hand
- After 18 holes, each player receives random cards equal to their total earned
- Best 5-card poker hand wins
- Standard poker hand rankings apply

## Endpoints

### Get Putt Putt Poker Status

```http
GET /api/games/{gameId}/side-bets/putt-putt-poker
```

**Response (200 OK):**
```json
{
  "bet_type": "putt_putt_poker",
  "status": "in_progress",
  "current_hole": 12,
  "players": [
    {
      "player": {
        "id": "player_123",
        "name": "John Doe"
      },
      "total_cards": 5,
      "starting_cards": 3,
      "cards_earned": 2,
      "penalties": 1,
      "putting_stats": {
        "one_putts": 2,
        "hole_in_ones": 0,
        "three_putts": 1,
        "average_putts": 1.8
      },
      "position": 1
    },
    {
      "player": {
        "id": "player_456",
        "name": "Jane Smith"
      },
      "total_cards": 4,
      "starting_cards": 3,
      "cards_earned": 1,
      "penalties": 0,
      "putting_stats": {
        "one_putts": 1,
        "hole_in_ones": 0,
        "three_putts": 0,
        "average_putts": 1.9
      },
      "position": 2
    }
  ],
  "pot_info": {
    "base_bet": 5.00,
    "penalty_additions": 1.00,
    "total_pot": 6.00
  }
}
```

### Get Player's Putt Putt Poker Details

```http
GET /api/games/{gameId}/players/{playerId}/side-bets/putt-putt-poker
```

**Response (200 OK):**
```json
{
  "player": {
    "id": "player_123",
    "name": "John Doe"
  },
  "card_history": [
    {
      "hole": 1,
      "putts": 1,
      "action": "card_awarded",
      "cards_change": +1,
      "total_cards": 4
    },
    {
      "hole": 3,
      "putts": 3,
      "action": "penalty",
      "penalty_amount": 1.00,
      "total_cards": 4
    },
    {
      "hole": 7,
      "putts": 1,
      "action": "card_awarded",
      "cards_change": +1,
      "total_cards": 5
    }
  ],
  "current_status": {
    "total_cards": 5,
    "cards_earned": 2,
    "penalties": 1,
    "holes_completed": 12,
    "position": 1
  },
  "putting_performance": {
    "best_putting_holes": [1, 7, 9],
    "worst_putting_holes": [3],
    "putting_average": 1.8,
    "putting_distribution": {
      "1_putt": 2,
      "2_putt": 9,
      "3_putt": 1,
      "4_putt": 0
    }
  }
}
```

### Deal Final Cards and Calculate Winner

```http
POST /api/games/{gameId}/side-bets/putt-putt-poker/deal
```

**Response (201 Created):**
```json
{
  "deal_timestamp": "2025-09-18T15:30:00Z",
  "random_seed": "a1b2c3d4e5f6",
  "players": [
    {
      "player": {
        "id": "player_123",
        "name": "John Doe"
      },
      "total_cards_earned": 5,
      "dealt_cards": ["AS", "KH", "QS", "JD", "10C"],
      "best_hand": {
        "cards": ["AS", "KH", "QS", "JD", "10C"],
        "hand_type": "straight",
        "hand_rank": 5,
        "description": "Ace-high straight"
      },
      "position": 1
    },
    {
      "player": {
        "id": "player_456",
        "name": "Jane Smith"
      },
      "total_cards_earned": 4,
      "dealt_cards": ["9H", "9D", "7S", "4C"],
      "best_hand": {
        "cards": ["9H", "9D", "7S", "4C"],
        "hand_type": "pair",
        "hand_rank": 2,
        "description": "Pair of nines"
      },
      "position": 2
    }
  ],
  "winner": {
    "player_id": "player_123",
    "hand_type": "straight",
    "winning_cards": ["AS", "KH", "QS", "JD", "10C"]
  },
  "pot_distribution": {
    "total_pot": 21.00,
    "winner_take": 21.00,
    "breakdown": {
      "base_bets": 20.00,
      "penalty_additions": 1.00
    }
  }
}
```

### Get Putting Statistics

```http
GET /api/games/{gameId}/side-bets/putt-putt-poker/stats
```

**Response (200 OK):**
```json
{
  "game_putting_stats": {
    "total_putts": 68,
    "average_putts_per_hole": 1.89,
    "best_putting_hole": {
      "hole": 5,
      "average_putts": 1.25
    },
    "worst_putting_hole": {
      "hole": 12,
      "average_putts": 2.75
    }
  },
  "distribution": {
    "hole_in_ones": 0,
    "one_putts": 12,
    "two_putts": 44,
    "three_putts": 10,
    "four_plus_putts": 2
  },
  "card_awards": {
    "total_cards_awarded": 12,
    "total_penalties": 12,
    "net_card_change": 0
  },
  "leaderboard": [
    {
      "player": {
        "id": "player_123",
        "name": "John Doe"
      },
      "cards": 5,
      "one_putts": 3,
      "penalties": 1
    }
  ]
}
```

### Get Live Card Updates

```http
GET /api/games/{gameId}/side-bets/putt-putt-poker/live
```

**Response (200 OK):**
```json
{
  "current_hole": 8,
  "recent_updates": [
    {
      "hole": 8,
      "player_id": "player_123",
      "putts": 1,
      "action": "card_awarded",
      "new_total": 5,
      "timestamp": "2025-09-18T12:15:00Z"
    },
    {
      "hole": 8,
      "player_id": "player_456",
      "putts": 3,
      "action": "penalty",
      "penalty_amount": 1.00,
      "timestamp": "2025-09-18T12:18:00Z"
    }
  ],
  "current_standings": [
    {
      "position": 1,
      "player_id": "player_123",
      "total_cards": 5
    },
    {
      "position": 2,
      "player_id": "player_456",
      "total_cards": 4
    }
  ]
}
```

## Poker Hand Rankings

### Hand Types (Highest to Lowest)
1. **Royal Flush** - A, K, Q, J, 10 all same suit
2. **Straight Flush** - 5 consecutive cards same suit
3. **Four of a Kind** - 4 cards same rank
4. **Full House** - 3 of a kind + pair
5. **Flush** - 5 cards same suit
6. **Straight** - 5 consecutive cards
7. **Three of a Kind** - 3 cards same rank
8. **Two Pair** - 2 pairs of same rank
9. **Pair** - 2 cards same rank
10. **High Card** - Highest single card

### Tie-Breaking Rules
- Standard poker tie-breaking rules apply
- Ace can be high (A-K-Q-J-10) or low (A-2-3-4-5)
- Suits have no ranking for tie-breaking

## Card Management

### Maximum Cards
- No maximum limit on cards earned
- Theoretical maximum: 3 starting + 36 possible (2 per hole-in-one) = 39 cards

### Minimum Cards
- Players cannot go below 0 cards
- No cards are lost for poor putting (only penalties added to pot)

### Random Dealing
- Uses cryptographically secure random number generation
- Deal is deterministic based on stored seed for verification
- Each player receives exactly their earned card count

## Error Responses

### Game Not Completed (400)
```json
{
  "error": "game_not_completed",
  "message": "Cannot deal final cards until all players complete 18 holes"
}
```

### Already Dealt (409)
```json
{
  "error": "cards_already_dealt",
  "message": "Final cards have already been dealt for this game",
  "deal_timestamp": "2025-09-18T15:30:00Z"
}
```

### Side Bet Not Enabled (400)
```json
{
  "error": "side_bet_not_enabled",
  "message": "Putt Putt Poker side bet is not enabled for this game"
}
```

### Invalid Putt Count (400)
```json
{
  "error": "invalid_putt_count",
  "message": "Putt count cannot exceed stroke count",
  "details": {
    "strokes": 4,
    "putts": 5
  }
}
```

## WebSocket Updates

Subscribe to real-time Putt Putt Poker updates:

```javascript
// WebSocket endpoint
ws://api.example.com/ws/games/{gameId}/side-bets/putt-putt-poker

// Example update message
{
  "type": "putt_putt_poker_update",
  "game_id": "game_abc123",
  "hole": 8,
  "player_id": "player_123",
  "putts": 1,
  "action": "card_awarded",
  "cards_change": +1,
  "new_total": 5,
  "standings": [...]
}
```