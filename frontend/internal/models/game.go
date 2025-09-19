package models

import "time"

// GameStatus represents the current state of a game
type GameStatus string

const (
	GameStatusSetup      GameStatus = "setup"
	GameStatusInProgress GameStatus = "in_progress"
	GameStatusCompleted  GameStatus = "completed"
	GameStatusAbandoned  GameStatus = "abandoned"
)

// SideBetType represents different types of side bets
type SideBetType string

const (
	SideBetBestNine       SideBetType = "best-nine"
	SideBetPuttPuttPoker  SideBetType = "putt-putt-poker"
)

// Gender represents player gender options
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// Game represents a golf game session
type Game struct {
	ID              string        `json:"id"`
	Course          string        `json:"course"`
	Status          GameStatus    `json:"status"`
	HandicapEnabled bool          `json:"handicap_enabled"`
	SideBets        []SideBetType `json:"side_bets"`
	ShareLink       string        `json:"share_link"`
	SpectatorLink   string        `json:"spectator_link"`
	CurrentHole     *int          `json:"current_hole,omitempty"`
	CreatedAt       time.Time     `json:"created_at"`
	StartedAt       *time.Time    `json:"started_at,omitempty"`
	CompletedAt     *time.Time    `json:"completed_at,omitempty"`
	Players         []Player      `json:"players"`
	CourseInfo      *CourseInfo   `json:"course_info,omitempty"`
	FinalResults    *FinalResults `json:"final_results,omitempty"`
}

// Player represents a golf player in a game
type Player struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Handicap  float64      `json:"handicap"`
	Gender    Gender       `json:"gender,omitempty"`
	Position  int          `json:"position"`
	GameID    string       `json:"game_id"`
	CreatedAt time.Time    `json:"created_at"`
	Stats     *PlayerStats `json:"stats,omitempty"`
}

// PlayerStats represents calculated player statistics
type PlayerStats struct {
	HolesCompleted int     `json:"holes_completed"`
	CurrentScore   string  `json:"current_score"`
	TotalPutts     int     `json:"total_putts"`
	PokerCards     *int    `json:"poker_cards,omitempty"`
	BestNineScore  *string `json:"best_nine_score,omitempty"`
}

// Score represents a score for a specific hole
type Score struct {
	ID              string            `json:"id"`
	PlayerID        string            `json:"player_id"`
	GameID          string            `json:"game_id"`
	Hole            int               `json:"hole"`
	Strokes         int               `json:"strokes"`
	Putts           int               `json:"putts"`
	Par             int               `json:"par"`
	ScoreToPar      string            `json:"score_to_par"`
	HandicapStroke  bool              `json:"handicap_stroke"`
	EffectiveScore  int               `json:"effective_score"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       *time.Time        `json:"updated_at,omitempty"`
	SideBetUpdates  *SideBetUpdates   `json:"side_bet_updates,omitempty"`
}

// SideBetUpdates represents side bet calculations for a score
type SideBetUpdates struct {
	PuttPuttPoker *PuttPuttPokerUpdate `json:"putt_putt_poker,omitempty"`
}

// PuttPuttPokerUpdate represents putt putt poker updates
type PuttPuttPokerUpdate struct {
	CardsAwarded   int  `json:"cards_awarded"`
	PenaltyApplied bool `json:"penalty_applied"`
	TotalCards     int  `json:"total_cards"`
}

// CourseInfo represents golf course information
type CourseInfo struct {
	Name     string     `json:"name"`
	Holes    []HoleInfo `json:"holes"`
	TotalPar int        `json:"total_par"`
}

// HoleInfo represents information about a specific hole
type HoleInfo struct {
	Hole             int     `json:"hole"`
	Par              int     `json:"par"`
	HandicapRanking  *int    `json:"handicap_ranking,omitempty"`
	Yardage          *int    `json:"yardage,omitempty"`
	Description      *string `json:"description,omitempty"`
}

// FinalResults represents the final results of a completed game
type FinalResults struct {
	Overall       []LeaderboardEntry `json:"overall"`
	BestNine      *BestNineStandings `json:"best_nine,omitempty"`
	PuttPuttPoker *PokerFinalResults `json:"putt_putt_poker,omitempty"`
}

// LeaderboardEntry represents a player's position in the leaderboard
type LeaderboardEntry struct {
	Position       int           `json:"position"`
	Player         PlayerSummary `json:"player"`
	Score          string        `json:"score"`
	HolesCompleted int           `json:"holes_completed"`
	TotalPutts     int           `json:"total_putts"`
	Trend          *string       `json:"trend,omitempty"`
}

// PlayerSummary represents basic player information
type PlayerSummary struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Handicap *float64 `json:"handicap,omitempty"`
}

// BestNineStandings represents the best nine side bet results
type BestNineStandings struct {
	BetType         string            `json:"bet_type"`
	Status          GameStatus        `json:"status"`
	HandicapEnabled bool              `json:"handicap_enabled"`
	Standings       []BestNineResult  `json:"standings"`
	Winner          *BestNineWinner   `json:"winner,omitempty"`
}

// BestNineResult represents a player's best nine result
type BestNineResult struct {
	Player              PlayerSummary `json:"player"`
	BestNineScore       string        `json:"best_nine_score"`
	HolesCompleted      int           `json:"holes_completed"`
	BestHoles           []int         `json:"best_holes"`
	WorstHoles          []int         `json:"worst_holes"`
	RawBestNine         string        `json:"raw_best_nine"`
	HandicapAdjustment  string        `json:"handicap_adjustment"`
	FinalScore          string        `json:"final_score"`
	Position            int           `json:"position"`
}

// BestNineWinner represents the winner of the best nine bet
type BestNineWinner struct {
	PlayerID string `json:"player_id"`
	Score    string `json:"score"`
	Margin   string `json:"margin"`
}

// PokerFinalResults represents the final poker results
type PokerFinalResults struct {
	BetType    string      `json:"bet_type"`
	Status     GameStatus  `json:"status"`
	Standings  []PokerHand `json:"standings"`
	Winner     *PokerWinner `json:"winner,omitempty"`
}

// PokerHand represents a player's final poker hand
type PokerHand struct {
	Player          PlayerSummary `json:"player"`
	TotalCardsEarned int          `json:"total_cards_earned"`
	DealtCards      []string     `json:"dealt_cards"`
	BestHand        PokerHandInfo `json:"best_hand"`
	Position        int          `json:"position"`
}

// PokerHandInfo represents poker hand details
type PokerHandInfo struct {
	Cards       []string      `json:"cards"`
	HandType    PokerHandType `json:"hand_type"`
	HandRank    int          `json:"hand_rank"`
	Description string       `json:"description"`
}

// PokerHandType represents different poker hand types
type PokerHandType string

const (
	PokerHandRoyalFlush    PokerHandType = "royal_flush"
	PokerHandStraightFlush PokerHandType = "straight_flush"
	PokerHandFourOfAKind   PokerHandType = "four_of_a_kind"
	PokerHandFullHouse     PokerHandType = "full_house"
	PokerHandFlush         PokerHandType = "flush"
	PokerHandStraight      PokerHandType = "straight"
	PokerHandThreeOfAKind  PokerHandType = "three_of_a_kind"
	PokerHandTwoPair       PokerHandType = "two_pair"
	PokerHandPair          PokerHandType = "pair"
	PokerHandHighCard      PokerHandType = "high_card"
)

// PokerWinner represents the winner of the poker bet
type PokerWinner struct {
	PlayerID string `json:"player_id"`
	HandType string `json:"hand_type"`
}

// CreateGameRequest represents the request to create a new game
type CreateGameRequest struct {
	Course          string        `json:"course"`
	SideBets        []SideBetType `json:"side_bets"`
	HandicapEnabled bool          `json:"handicap_enabled"`
}

// CreatePlayerRequest represents the request to add a player to a game
type CreatePlayerRequest struct {
	Name     string  `json:"name"`
	Handicap float64 `json:"handicap"`
	Gender   Gender  `json:"gender,omitempty"`
}

// CreateScoreRequest represents the request to record a score
type CreateScoreRequest struct {
	Hole    int `json:"hole"`
	Strokes int `json:"strokes"`
	Putts   int `json:"putts"`
}

// APIError represents an API error response
type APIError struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id"`
	Timestamp time.Time              `json:"timestamp"`
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}