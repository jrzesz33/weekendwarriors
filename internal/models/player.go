package models

import (
	"time"
)

// Gender represents player gender for handicap calculations
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// Player represents a golfer in a game
type Player struct {
	ID        string       `json:"id" db:"id"`
	GameID    string       `json:"game_id" db:"game_id"`
	Name      string       `json:"name" db:"name"`
	Handicap  float64      `json:"handicap" db:"handicap"`
	Gender    *Gender      `json:"gender,omitempty" db:"gender"`
	Position  int          `json:"position" db:"position"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	Stats     *PlayerStats `json:"stats,omitempty"`
}

// CreatePlayerRequest represents the request to add a player to a game
type CreatePlayerRequest struct {
	Name     string  `json:"name" validate:"required,min=1,max=100"`
	Handicap float64 `json:"handicap" validate:"required,min=0,max=54"`
	Gender   *Gender `json:"gender,omitempty" validate:"omitempty,oneof=male female other"`
}

// UpdatePlayerRequest represents the request to update a player
type UpdatePlayerRequest struct {
	Name     *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Handicap *float64 `json:"handicap,omitempty" validate:"omitempty,min=0,max=54"`
}

// PlayerStats represents player statistics during a game
type PlayerStats struct {
	HolesCompleted int    `json:"holes_completed"`
	CurrentScore   string `json:"current_score"`
	TotalPutts     int    `json:"total_putts"`
	PokerCards     *int   `json:"poker_cards,omitempty"`
	BestNineScore  string `json:"best_nine_score,omitempty"`
}

// PlayerSummary represents a condensed player view
type PlayerSummary struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Handicap *float64 `json:"handicap,omitempty"`
}

// PlayerDetail represents detailed player information including scores
type PlayerDetail struct {
	Player
	HoleScores []Score `json:"hole_scores"`
}

// PlayersResponse represents the response when getting all players
type PlayersResponse struct {
	Players    []Player `json:"players"`
	TotalCount int      `json:"total_count"`
}