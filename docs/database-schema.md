# Database Schema Documentation

## Overview

The Golf Gamez database schema is designed to support real-time golf game tracking with side bet calculations. The schema uses a relational approach optimized for read-heavy operations during active games.

## Core Tables

### games

Stores golf game sessions and metadata.

```sql
CREATE TABLE games (
    id VARCHAR(50) PRIMARY KEY,              -- game_abc123def456
    course VARCHAR(100) NOT NULL,            -- 'diamond-run'
    status VARCHAR(20) NOT NULL,             -- 'setup', 'in_progress', 'completed', 'abandoned'
    handicap_enabled BOOLEAN NOT NULL DEFAULT true,
    side_bets JSON,                          -- ['best-nine', 'putt-putt-poker']
    share_token VARCHAR(100) UNIQUE NOT NULL, -- for public sharing
    spectator_token VARCHAR(100) UNIQUE NOT NULL, -- for spectator access
    current_hole INTEGER,                    -- 1-18, NULL if not started
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    final_results JSON,                      -- computed results after completion

    INDEX idx_games_status (status),
    INDEX idx_games_share_token (share_token),
    INDEX idx_games_spectator_token (spectator_token),
    INDEX idx_games_created_at (created_at)
);
```

### players

Stores player information within each game.

```sql
CREATE TABLE players (
    id VARCHAR(50) PRIMARY KEY,              -- player_123abc456def
    game_id VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    handicap DECIMAL(4,1) NOT NULL,          -- 0.0 to 54.0
    gender ENUM('male', 'female', 'other'),
    position INTEGER NOT NULL,               -- tee-off order 1,2,3,4
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    UNIQUE KEY unique_player_name_per_game (game_id, name),
    UNIQUE KEY unique_position_per_game (game_id, position),
    INDEX idx_players_game_id (game_id)
);
```

### scores

Stores individual hole performance for each player.

```sql
CREATE TABLE scores (
    id VARCHAR(50) PRIMARY KEY,              -- score_abc123def456
    player_id VARCHAR(50) NOT NULL,
    game_id VARCHAR(50) NOT NULL,
    hole INTEGER NOT NULL,                   -- 1-18
    strokes INTEGER NOT NULL,                -- actual strokes taken
    putts INTEGER NOT NULL,                  -- putts taken
    par INTEGER NOT NULL,                    -- hole par (from course data)
    handicap_stroke BOOLEAN NOT NULL DEFAULT false, -- did player get handicap stroke
    score_to_par INTEGER NOT NULL,          -- raw score vs par
    effective_score INTEGER NOT NULL,       -- score after handicap adjustment
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),

    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    UNIQUE KEY unique_player_hole_score (player_id, hole),
    INDEX idx_scores_player_id (player_id),
    INDEX idx_scores_game_id (game_id),
    INDEX idx_scores_hole (hole),
    INDEX idx_scores_game_hole (game_id, hole)
);
```

### side_bet_calculations

Stores computed side bet results and standings.

```sql
CREATE TABLE side_bet_calculations (
    id VARCHAR(50) PRIMARY KEY,
    game_id VARCHAR(50) NOT NULL,
    player_id VARCHAR(50) NOT NULL,
    bet_type ENUM('best_nine', 'putt_putt_poker') NOT NULL,
    calculation_data JSON NOT NULL,          -- bet-specific calculation details
    current_position INTEGER,
    final_position INTEGER,
    is_winner BOOLEAN DEFAULT false,
    calculated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
    UNIQUE KEY unique_player_bet_type (player_id, bet_type),
    INDEX idx_sidebet_game_bet (game_id, bet_type),
    INDEX idx_sidebet_player_bet (player_id, bet_type)
);
```

### putt_putt_poker_cards

Tracks card earnings and penalties for Putt Putt Poker side bet.

```sql
CREATE TABLE putt_putt_poker_cards (
    id VARCHAR(50) PRIMARY KEY,
    player_id VARCHAR(50) NOT NULL,
    game_id VARCHAR(50) NOT NULL,
    hole INTEGER,                            -- NULL for starting cards
    action ENUM('starting', 'one_putt', 'hole_in_one', 'penalty') NOT NULL,
    cards_change INTEGER NOT NULL,           -- +1, +2, or 0 for penalty
    penalty_amount DECIMAL(10,2),            -- dollar amount for penalties
    total_cards INTEGER NOT NULL,            -- running total after this action
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    INDEX idx_poker_cards_player (player_id),
    INDEX idx_poker_cards_game (game_id),
    INDEX idx_poker_cards_hole (hole)
);
```

### poker_hands

Stores final dealt cards and poker hand results.

```sql
CREATE TABLE poker_hands (
    id VARCHAR(50) PRIMARY KEY,
    game_id VARCHAR(50) NOT NULL,
    player_id VARCHAR(50) NOT NULL,
    total_cards_earned INTEGER NOT NULL,
    dealt_cards JSON NOT NULL,               -- array of card strings ["AS", "KH", ...]
    best_hand_cards JSON NOT NULL,           -- best 5-card combination
    hand_type VARCHAR(50) NOT NULL,          -- 'straight', 'full_house', etc.
    hand_rank INTEGER NOT NULL,              -- 1-10 for comparison
    hand_description VARCHAR(200),           -- human readable description
    position INTEGER NOT NULL,
    deal_timestamp TIMESTAMP NOT NULL,
    random_seed VARCHAR(100) NOT NULL,       -- for verification

    FOREIGN KEY (game_id) REFERENCES games(id) ON DELETE CASCADE,
    FOREIGN KEY (player_id) REFERENCES players(id) ON DELETE CASCADE,
    UNIQUE KEY unique_player_hand (game_id, player_id),
    INDEX idx_poker_hands_game (game_id),
    INDEX idx_poker_hands_rank (hand_rank)
);
```

### course_data

Stores golf course hole information (prepopulated for Diamond Run).

```sql
CREATE TABLE course_data (
    id VARCHAR(50) PRIMARY KEY,
    course_name VARCHAR(100) NOT NULL,
    hole INTEGER NOT NULL,                   -- 1-18
    par INTEGER NOT NULL,                    -- 3, 4, or 5
    handicap_ranking INTEGER NOT NULL,       -- 1-18 hole difficulty ranking
    yardage INTEGER,                         -- optional yardage info
    description TEXT,                        -- optional hole description

    UNIQUE KEY unique_course_hole (course_name, hole),
    INDEX idx_course_data_course (course_name)
);
```

## Indexes and Performance

### Primary Indexes
- All primary keys are indexed automatically
- Foreign key relationships are indexed for join performance

### Query-Specific Indexes
- `idx_games_status`: For filtering active games
- `idx_scores_game_hole`: For leaderboard queries by hole
- `idx_sidebet_game_bet`: For side bet leaderboards
- `idx_poker_cards_player`: For card history queries

### Composite Indexes
- `unique_player_hole_score`: Ensures one score per player per hole
- `idx_scores_game_hole`: Optimizes hole-by-hole leaderboard queries

## Data Types and Constraints

### ID Generation
- All IDs use prefixed format: `{entity}_{random}`
- Example: `game_abc123def456`, `player_123abc456def`
- 50 character limit accommodates prefix + random string

### JSON Fields
- `games.side_bets`: Array of enabled side bet types
- `games.final_results`: Computed final standings and winners
- `side_bet_calculations.calculation_data`: Bet-specific computed data
- `poker_hands.dealt_cards`: Array of card strings
- `poker_hands.best_hand_cards`: Best 5-card poker hand

### Decimal Precision
- `players.handicap`: DECIMAL(4,1) supports 0.0 to 99.9 handicaps
- `putt_putt_poker_cards.penalty_amount`: DECIMAL(10,2) for currency

### Enums
- `games.status`: Prevents invalid game states
- `players.gender`: For handicap suggestions
- `putt_putt_poker_cards.action`: Card earning action types

## Sample Data

### Diamond Run Course Data
```sql
INSERT INTO course_data (id, course_name, hole, par, handicap_ranking) VALUES
('hole_dr_01', 'diamond-run', 1, 4, 10),
('hole_dr_02', 'diamond-run', 2, 3, 18),
('hole_dr_03', 'diamond-run', 3, 5, 2),
('hole_dr_04', 'diamond-run', 4, 4, 8),
('hole_dr_05', 'diamond-run', 5, 3, 16),
('hole_dr_06', 'diamond-run', 6, 4, 12),
('hole_dr_07', 'diamond-run', 7, 4, 6),
('hole_dr_08', 'diamond-run', 8, 3, 14),
('hole_dr_09', 'diamond-run', 9, 5, 4),
('hole_dr_10', 'diamond-run', 10, 4, 9),
('hole_dr_11', 'diamond-run', 11, 3, 17),
('hole_dr_12', 'diamond-run', 12, 5, 1),
('hole_dr_13', 'diamond-run', 13, 3, 15),
('hole_dr_14', 'diamond-run', 14, 4, 11),
('hole_dr_15', 'diamond-run', 15, 4, 7),
('hole_dr_16', 'diamond-run', 16, 4, 13),
('hole_dr_17', 'diamond-run', 17, 5, 3),
('hole_dr_18', 'diamond-run', 18, 4, 5);
```

## Migration Strategy

### Version 1.0 (Initial Schema)
- Core tables: games, players, scores, course_data
- Basic game and score tracking functionality

### Version 1.1 (Side Bets)
- Add: side_bet_calculations, putt_putt_poker_cards, poker_hands
- Migrate existing games to include side bet columns

### Future Considerations
- Partitioning scores table by game_id for large datasets
- Read replicas for spectator queries
- Archival strategy for completed games
- Additional course support (normalized course management)