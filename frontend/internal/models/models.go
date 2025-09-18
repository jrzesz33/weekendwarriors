package models

import "time"

// Game represents a golf game session
type Game struct {
	ID             string           `json:"id"`
	Course         string           `json:"course"`
	Status         string           `json:"status"`
	HandicapEnabled bool            `json:"handicap_enabled"`
	SideBets       []string         `json:"side_bets"`
	ShareLink      string           `json:"share_link"`
	SpectatorLink  string           `json:"spectator_link"`
	CurrentHole    *int             `json:"current_hole"`
	CreatedAt      time.Time        `json:"created_at"`
	StartedAt      *time.Time       `json:"started_at"`
	CompletedAt    *time.Time       `json:"completed_at"`
	Players        []Player         `json:"players"`
	CourseInfo     CourseInfo       `json:"course_info"`
	FinalResults   *FinalResults    `json:"final_results"`
}

// Player represents a golfer in the game
type Player struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Handicap  float64     `json:"handicap"`
	Gender    string      `json:"gender"`
	Position  int         `json:"position"`
	GameID    string      `json:"game_id"`
	CreatedAt time.Time   `json:"created_at"`
	Stats     PlayerStats `json:"stats"`
}

// PlayerStats contains player performance statistics
type PlayerStats struct {
	HolesCompleted int    `json:"holes_completed"`
	CurrentScore   string `json:"current_score"`
	TotalPutts     int    `json:"total_putts"`
	PokerCards     int    `json:"poker_cards"`
	BestNineScore  string `json:"best_nine_score"`
}

// Score represents a player's score for a specific hole
type Score struct {
	ID             string    `json:"id"`
	PlayerID       string    `json:"player_id"`
	GameID         string    `json:"game_id"`
	Hole           int       `json:"hole"`
	Strokes        int       `json:"strokes"`
	Putts          int       `json:"putts"`
	Par            int       `json:"par"`
	ScoreToPar     string    `json:"score_to_par"`
	HandicapStroke bool      `json:"handicap_stroke"`
	EffectiveScore int       `json:"effective_score"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CourseInfo contains golf course information
type CourseInfo struct {
	Name     string     `json:"name"`
	Holes    []HoleInfo `json:"holes"`
	TotalPar int        `json:"total_par"`
}

// HoleInfo contains information about a specific hole
type HoleInfo struct {
	Hole             int    `json:"hole"`
	Par              int    `json:"par"`
	HandicapRanking  int    `json:"handicap_ranking"`
	Yardage          int    `json:"yardage"`
	Description      string `json:"description"`
}

// Leaderboard contains current game standings
type Leaderboard struct {
	Overall   []LeaderboardEntry `json:"overall"`
	SideBets  SideBetLeaderboard `json:"side_bets"`
}

// LeaderboardEntry represents a player's position in the leaderboard
type LeaderboardEntry struct {
	Position       int           `json:"position"`
	Player         PlayerSummary `json:"player"`
	Score          string        `json:"score"`
	HolesCompleted int           `json:"holes_completed"`
	TotalPutts     int           `json:"total_putts"`
	Trend          string        `json:"trend"`
}

// SideBetLeaderboard contains side bet standings
type SideBetLeaderboard struct {
	BestNine      []BestNineResult      `json:"best_nine"`
	PuttPuttPoker []PuttPuttPokerResult `json:"putt_putt_poker"`
}

// BestNineResult represents Best Nine side bet results
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

// PuttPuttPokerResult represents Putt Putt Poker side bet results
type PuttPuttPokerResult struct {
	Player       PlayerSummary `json:"player"`
	TotalCards   int           `json:"total_cards"`
	StartingCards int          `json:"starting_cards"`
	CardsEarned  int           `json:"cards_earned"`
	Penalties    int           `json:"penalties"`
	PuttingStats PuttingStats  `json:"putting_stats"`
	Position     int           `json:"position"`
}

// PuttingStats contains putting performance statistics
type PuttingStats struct {
	OnePutts     int     `json:"one_putts"`
	HoleInOnes   int     `json:"hole_in_ones"`
	ThreePutts   int     `json:"three_putts"`
	AveragePutts float64 `json:"average_putts"`
}

// PlayerSummary contains basic player information
type PlayerSummary struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Handicap float64 `json:"handicap"`
}

// FinalResults contains end-of-game results
type FinalResults struct {
	OverallWinner      *GameWinner       `json:"overall_winner"`
	BestNineWinner     *GameWinner       `json:"best_nine_winner"`
	PuttPuttPokerWinner *PokerGameWinner `json:"putt_putt_poker_winner"`
}

// GameWinner represents a game winner
type GameWinner struct {
	PlayerID string `json:"player_id"`
	Score    string `json:"score"`
}

// PokerGameWinner represents a poker game winner
type PokerGameWinner struct {
	PlayerID string   `json:"player_id"`
	Hand     string   `json:"hand"`
	Cards    []string `json:"cards"`
}

// CreateGameRequest represents a request to create a new game
type CreateGameRequest struct {
	Course          string   `json:"course"`
	SideBets        []string `json:"side_bets"`
	HandicapEnabled bool     `json:"handicap_enabled"`
}

// CreatePlayerRequest represents a request to add a player
type CreatePlayerRequest struct {
	Name     string  `json:"name"`
	Handicap float64 `json:"handicap"`
	Gender   string  `json:"gender"`
}

// ScoreRequest represents a request to record a score
type ScoreRequest struct {
	Hole    int `json:"hole"`
	Strokes int `json:"strokes"`
	Putts   int `json:"putts"`
}

// WebSocketMessage represents a real-time update message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	GameID  string      `json:"game_id"`
	Payload interface{} `json:"payload"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details"`
	RequestID string      `json:"request_id"`
	Timestamp time.Time   `json:"timestamp"`
}

// GameCompletionResult represents the result of completing a game
type GameCompletionResult struct {
	ID           string        `json:"id"`
	Status       string        `json:"status"`
	CompletedAt  time.Time     `json:"completed_at"`
	FinalResults *FinalResults `json:"final_results"`
}

// GameScorecard represents the full scorecard for a game
type GameScorecard struct {
	Game       GameSummary     `json:"game"`
	CourseInfo []HoleInfo      `json:"course_info"`
	Players    []PlayerDetail  `json:"players"`
}

// GameSummary contains basic game information
type GameSummary struct {
	ID          string `json:"id"`
	Course      string `json:"course"`
	CurrentHole *int   `json:"current_hole"`
}

// PlayerDetail contains detailed player information with scores
type PlayerDetail struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Position int         `json:"position"`
	Scores   []Score     `json:"scores"`
	Totals   PlayerStats `json:"totals"`
}

// BestNineStandings represents Best Nine side bet standings
type BestNineStandings struct {
	BetType         string           `json:"bet_type"`
	Status          string           `json:"status"`
	HandicapEnabled bool             `json:"handicap_enabled"`
	Standings       []BestNineResult `json:"standings"`
	Winner          *GameWinner      `json:"winner"`
}

// PuttPuttPokerStatus represents Putt Putt Poker side bet status
type PuttPuttPokerStatus struct {
	BetType     string                `json:"bet_type"`
	Status      string                `json:"status"`
	CurrentHole int                   `json:"current_hole"`
	Players     []PuttPuttPokerResult `json:"players"`
	PotInfo     PotInfo               `json:"pot_info"`
}

// PotInfo contains pot information for Putt Putt Poker
type PotInfo struct {
	BaseBet          float64 `json:"base_bet"`
	PenaltyAdditions float64 `json:"penalty_additions"`
	TotalPot         float64 `json:"total_pot"`
}

// PokerDealResult represents the result of dealing poker cards
type PokerDealResult struct {
	DealTimestamp time.Time         `json:"deal_timestamp"`
	RandomSeed    string            `json:"random_seed"`
	Players       []PokerHand       `json:"players"`
	Winner        *PokerGameWinner  `json:"winner"`
	PotDistribution PotDistribution `json:"pot_distribution"`
}

// PokerHand represents a player's final poker hand
type PokerHand struct {
	Player           PlayerSummary `json:"player"`
	TotalCardsEarned int           `json:"total_cards_earned"`
	DealtCards       []string      `json:"dealt_cards"`
	BestHand         BestPokerHand `json:"best_hand"`
	Position         int           `json:"position"`
}

// BestPokerHand represents the best 5-card poker hand
type BestPokerHand struct {
	Cards       []string `json:"cards"`
	HandType    string   `json:"hand_type"`
	HandRank    int      `json:"hand_rank"`
	Description string   `json:"description"`
}

// PotDistribution contains pot distribution information
type PotDistribution struct {
	TotalPot   float64               `json:"total_pot"`
	WinnerTake float64               `json:"winner_take"`
	Breakdown  PotDistributionDetail `json:"breakdown"`
}

// PotDistributionDetail contains detailed pot breakdown
type PotDistributionDetail struct {
	BaseBets         float64 `json:"base_bets"`
	PenaltyAdditions float64 `json:"penalty_additions"`
}

// SpectatorView represents the spectator view of a game
type SpectatorView struct {
	Game           Game        `json:"game"`
	Leaderboard    Leaderboard `json:"leaderboard"`
	LiveUpdates    bool        `json:"live_updates"`
	SpectatorCount int         `json:"spectator_count"`
}