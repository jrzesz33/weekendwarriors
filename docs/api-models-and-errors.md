# API Response Models and Error Handling

## Overview

This document defines the standard response models, error formats, and data structures used throughout the Golf Gamez API.

## Standard Response Envelope

### Success Response Format
All successful API responses follow a consistent structure:

```json
{
  "data": {},
  "meta": {
    "timestamp": "2025-09-18T10:30:00Z",
    "version": "v1",
    "request_id": "req_abc123def456"
  }
}
```

For simple responses, the envelope may be omitted and data returned directly.

### Error Response Format
All error responses use a standard format:

```json
{
  "error": {
    "code": "validation_error",
    "message": "Invalid input data",
    "details": {
      "field": "handicap",
      "value": -5,
      "constraint": "must be between 0 and 54"
    },
    "request_id": "req_abc123def456",
    "timestamp": "2025-09-18T10:30:00Z"
  }
}
```

## Core Data Models

### Game Model
```typescript
interface Game {
  id: string;                    // game_abc123def456
  course: string;                // 'diamond-run'
  status: GameStatus;            // 'setup' | 'in_progress' | 'completed' | 'abandoned'
  handicap_enabled: boolean;
  side_bets: SideBetType[];      // ['best-nine', 'putt-putt-poker']
  share_link: string;            // public game URL
  spectator_link: string;        // spectator view URL
  current_hole?: number;         // 1-18, null if not started
  created_at: string;            // ISO 8601 timestamp
  started_at?: string;           // ISO 8601 timestamp
  completed_at?: string;         // ISO 8601 timestamp
  players: Player[];
  course_info?: CourseInfo;
  final_results?: FinalResults;
}

type GameStatus = 'setup' | 'in_progress' | 'completed' | 'abandoned';
type SideBetType = 'best-nine' | 'putt-putt-poker';
```

### Player Model
```typescript
interface Player {
  id: string;                    // player_123abc456def
  name: string;                  // 1-100 characters
  handicap: number;              // 0.0 to 54.0
  gender?: 'male' | 'female' | 'other';
  position: number;              // tee-off order 1,2,3,4
  game_id: string;
  created_at: string;
  stats?: PlayerStats;
}

interface PlayerStats {
  holes_completed: number;       // 0-18
  current_score: string;         // "+5", "E", "-2"
  total_putts: number;
  poker_cards?: number;          // for putt putt poker
  best_nine_score?: string;      // best nine calculation
}
```

### Score Model
```typescript
interface Score {
  id: string;                    // score_abc123def456
  player_id: string;
  game_id: string;
  hole: number;                  // 1-18
  strokes: number;               // 1+
  putts: number;                 // 0+ (0 for hole-in-one)
  par: number;                   // 3, 4, or 5
  score_to_par: string;          // "+1", "E", "-1"
  handicap_stroke: boolean;      // did player receive handicap stroke
  effective_score: number;       // score after handicap adjustment
  created_at: string;
  updated_at?: string;
  side_bet_updates?: SideBetUpdates;
}

interface SideBetUpdates {
  putt_putt_poker?: {
    cards_awarded: number;
    penalty_applied: boolean;
    total_cards: number;
  };
}
```

### Course Model
```typescript
interface CourseInfo {
  name: string;                  // 'Diamond Run'
  holes: HoleInfo[];
  total_par: number;             // sum of all hole pars
}

interface HoleInfo {
  hole: number;                  // 1-18
  par: number;                   // 3, 4, or 5
  handicap_ranking?: number;     // 1-18 difficulty ranking
  yardage?: number;
  description?: string;
}
```

### Side Bet Models

#### Best Nine Model
```typescript
interface BestNineResult {
  player: PlayerSummary;
  best_nine_score: string;       // "+2", "E", "-1"
  holes_completed: number;
  best_holes: number[];          // array of hole numbers used
  worst_holes: number[];         // array of hole numbers discarded
  raw_best_nine: string;         // score before handicap
  handicap_adjustment: string;   // handicap strokes applied
  final_score: string;           // final calculated score
  position: number;
}

interface BestNineStandings {
  bet_type: 'best_nine';
  status: GameStatus;
  handicap_enabled: boolean;
  standings: BestNineResult[];
  winner?: {
    player_id: string;
    score: string;
    margin: string;
  };
}
```

#### Putt Putt Poker Model
```typescript
interface PuttPuttPokerResult {
  player: PlayerSummary;
  total_cards: number;
  starting_cards: number;        // always 3
  cards_earned: number;
  penalties: number;             // count of 3+ putt penalties
  putting_stats: PuttingStats;
  position: number;
}

interface PuttingStats {
  one_putts: number;
  hole_in_ones: number;
  three_putts: number;
  average_putts: number;
}

interface PokerHand {
  player: PlayerSummary;
  total_cards_earned: number;
  dealt_cards: string[];         // ["AS", "KH", "QS", "JD", "10C"]
  best_hand: {
    cards: string[];             // best 5-card combination
    hand_type: PokerHandType;
    hand_rank: number;           // 1-10 for comparison
    description: string;         // human readable
  };
  position: number;
}

type PokerHandType = 'royal_flush' | 'straight_flush' | 'four_of_a_kind' |
                     'full_house' | 'flush' | 'straight' | 'three_of_a_kind' |
                     'two_pair' | 'pair' | 'high_card';
```

### Leaderboard Models
```typescript
interface Leaderboard {
  overall: LeaderboardEntry[];
  side_bets?: {
    best_nine?: BestNineResult[];
    putt_putt_poker?: PuttPuttPokerResult[];
  };
}

interface LeaderboardEntry {
  position: number;
  player: PlayerSummary;
  score: string;                 // "+15", "E", "-3"
  holes_completed: number;
  total_putts: number;
  trend?: 'up' | 'down' | 'same'; // position change from last update
}

interface PlayerSummary {
  id: string;
  name: string;
  handicap?: number;
}
```

## Standard HTTP Status Codes

### Success Codes
- **200 OK**: Successful GET, PUT, or DELETE operation
- **201 Created**: Successful POST operation creating new resource
- **204 No Content**: Successful DELETE operation with no response body

### Client Error Codes
- **400 Bad Request**: Invalid request data or parameters
- **401 Unauthorized**: Authentication required (future use)
- **403 Forbidden**: Access denied (future use)
- **404 Not Found**: Resource not found
- **409 Conflict**: Resource conflict (e.g., duplicate name)
- **422 Unprocessable Entity**: Valid JSON but business logic error

### Server Error Codes
- **500 Internal Server Error**: Unexpected server error
- **503 Service Unavailable**: Service temporarily unavailable

## Error Types and Codes

### Validation Errors (400)
```typescript
interface ValidationError {
  code: 'validation_error';
  message: string;
  details: {
    field: string;
    value: any;
    constraint: string;
    allowed_values?: any[];
  };
}
```

**Common validation errors:**
- `invalid_handicap`: Handicap outside 0-54 range
- `invalid_hole_number`: Hole number outside 1-18 range
- `invalid_score_values`: Strokes/putts with impossible values
- `missing_required_field`: Required field not provided

### Business Logic Errors (400, 409, 422)
```typescript
interface BusinessLogicError {
  code: string;
  message: string;
  details?: {
    current_state?: string;
    required_state?: string;
    additional_info?: any;
  };
}
```

**Common business logic errors:**
- `game_not_started`: Cannot perform action before game starts
- `game_already_completed`: Cannot modify completed game
- `player_limit_exceeded`: Maximum 4 players per game
- `duplicate_player_name`: Player name already exists in game
- `score_already_exists`: Score for hole already recorded
- `future_hole`: Cannot record scores for future holes
- `invalid_game_state`: Game state prevents requested operation

### Resource Errors (404)
```typescript
interface ResourceError {
  code: 'resource_not_found';
  message: string;
  details: {
    resource_type: string;
    resource_id: string;
  };
}
```

**Resource error types:**
- `game_not_found`: Game ID does not exist
- `player_not_found`: Player ID does not exist in game
- `score_not_found`: Score for player/hole combination not found

### Side Bet Errors (400)
```typescript
interface SideBetError {
  code: string;
  message: string;
  details?: {
    bet_type: SideBetType;
    requirement?: string;
    current_value?: any;
  };
}
```

**Side bet error types:**
- `side_bet_not_enabled`: Side bet not enabled for this game
- `insufficient_holes`: Not enough holes completed for calculation
- `cards_already_dealt`: Poker cards already dealt for game
- `invalid_putt_count`: Putt count exceeds stroke count

## Pagination

For endpoints returning large datasets:

```typescript
interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    per_page: number;
    total_count: number;
    total_pages: number;
    has_next: boolean;
    has_prev: boolean;
  };
}
```

**Query parameters:**
- `page`: Page number (default: 1)
- `per_page`: Items per page (default: 20, max: 100)
- `sort`: Sort field and direction (e.g., `created_at:desc`)

## Rate Limiting

Rate limit headers included in all responses:

```http
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642781400
X-RateLimit-Window: 3600
```

## Request/Response Examples

### Successful Game Creation
```http
POST /api/games
Content-Type: application/json

{
  "course": "diamond-run",
  "side_bets": ["best-nine", "putt-putt-poker"],
  "handicap_enabled": true
}
```

```http
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": "game_abc123def456",
  "course": "diamond-run",
  "status": "setup",
  "handicap_enabled": true,
  "side_bets": ["best-nine", "putt-putt-poker"],
  "share_link": "https://api.example.com/games/abc123def456",
  "spectator_link": "https://api.example.com/spectate/abc123def456",
  "created_at": "2025-09-18T10:30:00Z",
  "players": [],
  "current_hole": null
}
```

### Validation Error Example
```http
POST /api/games/abc123/players
Content-Type: application/json

{
  "name": "John Doe",
  "handicap": -5
}
```

```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": {
    "code": "validation_error",
    "message": "Handicap must be between 0 and 54",
    "details": {
      "field": "handicap",
      "value": -5,
      "constraint": "must be between 0 and 54"
    },
    "request_id": "req_abc123def456",
    "timestamp": "2025-09-18T10:30:00Z"
  }
}
```