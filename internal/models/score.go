package models

import (
	"fmt"
	"time"
)

// Score represents a player's score for a specific hole
type Score struct {
	ID              string           `json:"id" db:"id"`
	PlayerID        string           `json:"player_id" db:"player_id"`
	GameID          string           `json:"game_id" db:"game_id"`
	Hole            int              `json:"hole" db:"hole"`
	Strokes         int              `json:"strokes" db:"strokes"`
	Putts           int              `json:"putts" db:"putts"`
	Par             int              `json:"par" db:"par"`
	ScoreToPar      string           `json:"score_to_par"`
	HandicapStroke  bool             `json:"handicap_stroke" db:"handicap_stroke"`
	EffectiveScore  int              `json:"effective_score" db:"effective_score"`
	CreatedAt       time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time       `json:"updated_at,omitempty" db:"updated_at"`
	SideBetUpdates  *SideBetUpdates  `json:"side_bet_updates,omitempty"`
}

// ScoreRequest represents the request to record a score
type ScoreRequest struct {
	Hole    int `json:"hole" validate:"required,min=1,max=18"`
	Strokes int `json:"strokes" validate:"required,min=1,max=20"`
	Putts   int `json:"putts" validate:"required,min=0,max=10"`
}

// UpdateScoreRequest represents the request to update a score
type UpdateScoreRequest struct {
	Strokes *int `json:"strokes,omitempty" validate:"omitempty,min=1,max=20"`
	Putts   *int `json:"putts,omitempty" validate:"omitempty,min=0,max=10"`
}

// SideBetUpdates represents side bet updates when a score is recorded
type SideBetUpdates struct {
	PuttPuttPoker *PuttPuttPokerUpdate `json:"putt_putt_poker,omitempty"`
}

// PuttPuttPokerUpdate represents poker updates for a score
type PuttPuttPokerUpdate struct {
	CardsAwarded   int  `json:"cards_awarded"`
	PenaltyApplied bool `json:"penalty_applied"`
	TotalCards     int  `json:"total_cards"`
}

// GameScorecard represents the complete scorecard for all players
type GameScorecard struct {
	Game       GameSummary        `json:"game"`
	CourseInfo []HoleInfo         `json:"course_info"`
	Players    []ScorecardPlayer  `json:"players"`
}

// GameSummary represents basic game information for scorecard
type GameSummary struct {
	ID          string `json:"id"`
	Course      string `json:"course"`
	CurrentHole *int   `json:"current_hole"`
}

// ScorecardPlayer represents player information with scores for scorecard
type ScorecardPlayer struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Position int          `json:"position"`
	Scores   []Score      `json:"scores"`
	Totals   *PlayerStats `json:"totals"`
}

// Leaderboard represents current game standings
type Leaderboard struct {
	Overall   []LeaderboardEntry `json:"overall"`
	SideBets  *SideBetLeaderboard `json:"side_bets,omitempty"`
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

// SideBetLeaderboard represents side bet standings
type SideBetLeaderboard struct {
	BestNine      []BestNineResult      `json:"best_nine,omitempty"`
	PuttPuttPoker []PuttPuttPokerResult `json:"putt_putt_poker,omitempty"`
}

// FormatScoreToPar formats a score relative to par
func FormatScoreToPar(strokes, par int) string {
	diff := strokes - par
	switch {
	case diff == 0:
		return "E"
	case diff > 0:
		return "+" + formatPositiveInt(diff)
	default:
		return formatNegativeInt(diff)
	}
}

// CalculateEffectiveScore calculates score after handicap adjustment
func CalculateEffectiveScore(strokes, par int, handicapStroke bool) int {
	if handicapStroke {
		return strokes - 1 - par
	}
	return strokes - par
}

// Helper functions
func formatPositiveInt(n int) string {
	return fmt.Sprintf("%d", n)
}

func formatNegativeInt(n int) string {
	return fmt.Sprintf("%d", n)
}