package models

import (
	"encoding/json"
	"time"
)

// GameStatus represents the current state of a golf game
type GameStatus string

const (
	GameStatusSetup       GameStatus = "setup"
	GameStatusInProgress  GameStatus = "in_progress"
	GameStatusCompleted   GameStatus = "completed"
	GameStatusAbandoned   GameStatus = "abandoned"
)

// SideBetType represents available side bet types
type SideBetType string

const (
	SideBetBestNine      SideBetType = "best-nine"
	SideBetPuttPuttPoker SideBetType = "putt-putt-poker"
)

// Game represents a golf game session
type Game struct {
	ID             string       `json:"id" db:"id"`
	Course         string       `json:"course" db:"course"`
	Status         GameStatus   `json:"status" db:"status"`
	HandicapEnabled bool        `json:"handicap_enabled" db:"handicap_enabled"`
	SideBets       []SideBetType `json:"side_bets"`
	ShareLink      string       `json:"share_link"`
	SpectatorLink  string       `json:"spectator_link"`
	ShareToken     string       `json:"-" db:"share_token"`
	SpectatorToken string       `json:"-" db:"spectator_token"`
	CurrentHole    *int         `json:"current_hole" db:"current_hole"`
	CreatedAt      time.Time    `json:"created_at" db:"created_at"`
	StartedAt      *time.Time   `json:"started_at" db:"started_at"`
	CompletedAt    *time.Time   `json:"completed_at" db:"completed_at"`
	Players        []Player     `json:"players,omitempty"`
	CourseInfo     *CourseInfo  `json:"course_info,omitempty"`
	FinalResults   *FinalResults `json:"final_results,omitempty"`
}

// CreateGameRequest represents the request to create a new game
type CreateGameRequest struct {
	Course          string        `json:"course" validate:"required,oneof=diamond-run"`
	SideBets        []SideBetType `json:"side_bets,omitempty"`
	HandicapEnabled bool          `json:"handicap_enabled"`
}

// FinalResults represents the final game results
type FinalResults struct {
	OverallWinner       *Winner `json:"overall_winner,omitempty"`
	BestNineWinner      *Winner `json:"best_nine_winner,omitempty"`
	PuttPuttPokerWinner *Winner `json:"putt_putt_poker_winner,omitempty"`
}

// Winner represents a game winner
type Winner struct {
	PlayerID string `json:"player_id"`
	Score    string `json:"score"`
	Hand     string `json:"hand,omitempty"`     // For poker
	Cards    []string `json:"cards,omitempty"`  // For poker
}

// CourseInfo represents golf course information
type CourseInfo struct {
	Name     string     `json:"name"`
	Holes    []HoleInfo `json:"holes"`
	TotalPar int        `json:"total_par"`
}

// HoleInfo represents information about a specific hole
type HoleInfo struct {
	Hole            int    `json:"hole"`
	Par             int    `json:"par"`
	HandicapRanking int    `json:"handicap_ranking"`
	Yardage         *int   `json:"yardage,omitempty"`
	Description     string `json:"description,omitempty"`
}

// GameCompletionResult represents the result of completing a game
type GameCompletionResult struct {
	ID           string        `json:"id"`
	Status       GameStatus    `json:"status"`
	CompletedAt  time.Time     `json:"completed_at"`
	FinalResults *FinalResults `json:"final_results"`
}

// SpectatorView represents the spectator view of a game
type SpectatorView struct {
	Game           *Game        `json:"game"`
	Leaderboard    *Leaderboard `json:"leaderboard"`
	LiveUpdates    bool         `json:"live_updates"`
	SpectatorCount int          `json:"spectator_count"`
}

// MarshalSideBets converts side bets slice to JSON string for database storage
func (g *Game) MarshalSideBets() (string, error) {
	if len(g.SideBets) == 0 {
		return "[]", nil
	}
	data, err := json.Marshal(g.SideBets)
	return string(data), err
}

// UnmarshalSideBets converts JSON string from database to side bets slice
func (g *Game) UnmarshalSideBets(data string) error {
	if data == "" || data == "[]" {
		g.SideBets = []SideBetType{}
		return nil
	}
	return json.Unmarshal([]byte(data), &g.SideBets)
}

// MarshalFinalResults converts final results to JSON string for database storage
func (g *Game) MarshalFinalResults() (string, error) {
	if g.FinalResults == nil {
		return "", nil
	}
	data, err := json.Marshal(g.FinalResults)
	return string(data), err
}

// UnmarshalFinalResults converts JSON string from database to final results
func (g *Game) UnmarshalFinalResults(data string) error {
	if data == "" {
		g.FinalResults = nil
		return nil
	}
	g.FinalResults = &FinalResults{}
	return json.Unmarshal([]byte(data), g.FinalResults)
}