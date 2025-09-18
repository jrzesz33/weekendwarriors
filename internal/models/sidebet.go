package models

import (
	"encoding/json"
	"time"
)

// BestNineResult represents a player's Best Nine side bet result
type BestNineResult struct {
	Player             PlayerSummary `json:"player"`
	BestNineScore      string        `json:"best_nine_score"`
	HolesCompleted     int           `json:"holes_completed"`
	BestHoles          []int         `json:"best_holes"`
	WorstHoles         []int         `json:"worst_holes"`
	RawBestNine        string        `json:"raw_best_nine"`
	HandicapAdjustment string        `json:"handicap_adjustment"`
	FinalScore         string        `json:"final_score"`
	Position           int           `json:"position"`
}

// BestNineStandings represents the Best Nine standings for a game
type BestNineStandings struct {
	BetType         SideBetType    `json:"bet_type"`
	Status          GameStatus     `json:"status"`
	HandicapEnabled bool           `json:"handicap_enabled"`
	Standings       []BestNineResult `json:"standings"`
	Winner          *BestNineWinner  `json:"winner,omitempty"`
}

// BestNineWinner represents the winner of Best Nine
type BestNineWinner struct {
	PlayerID string `json:"player_id"`
	Score    string `json:"score"`
	Margin   string `json:"margin"`
}

// PuttPuttPokerResult represents a player's Putt Putt Poker status
type PuttPuttPokerResult struct {
	Player       PlayerSummary `json:"player"`
	TotalCards   int           `json:"total_cards"`
	StartingCards int          `json:"starting_cards"`
	CardsEarned  int           `json:"cards_earned"`
	Penalties    int           `json:"penalties"`
	PuttingStats PuttingStats  `json:"putting_stats"`
	Position     int           `json:"position"`
}

// PuttPuttPokerStatus represents the overall Putt Putt Poker status
type PuttPuttPokerStatus struct {
	BetType     SideBetType           `json:"bet_type"`
	Status      GameStatus            `json:"status"`
	CurrentHole *int                  `json:"current_hole"`
	Players     []PuttPuttPokerResult `json:"players"`
	PotInfo     *PotInfo              `json:"pot_info,omitempty"`
}

// PuttingStats represents putting performance statistics
type PuttingStats struct {
	OnePutts     int     `json:"one_putts"`
	HoleInOnes   int     `json:"hole_in_ones"`
	ThreePutts   int     `json:"three_putts"`
	AveragePutts float64 `json:"average_putts"`
}

// PotInfo represents betting pot information
type PotInfo struct {
	BaseBet          float64 `json:"base_bet"`
	PenaltyAdditions float64 `json:"penalty_additions"`
	TotalPot         float64 `json:"total_pot"`
}

// PokerHandType represents poker hand types
type PokerHandType string

const (
	RoyalFlush    PokerHandType = "royal_flush"
	StraightFlush PokerHandType = "straight_flush"
	FourOfAKind   PokerHandType = "four_of_a_kind"
	FullHouse     PokerHandType = "full_house"
	Flush         PokerHandType = "flush"
	Straight      PokerHandType = "straight"
	ThreeOfAKind  PokerHandType = "three_of_a_kind"
	TwoPair       PokerHandType = "two_pair"
	Pair          PokerHandType = "pair"
	HighCard      PokerHandType = "high_card"
)

// PokerHand represents a dealt poker hand
type PokerHand struct {
	Player            PlayerSummary `json:"player"`
	TotalCardsEarned  int           `json:"total_cards_earned"`
	DealtCards        []string      `json:"dealt_cards"`
	BestHand          BestPokerHand `json:"best_hand"`
	Position          int           `json:"position"`
}

// BestPokerHand represents the best 5-card poker hand
type BestPokerHand struct {
	Cards       []string      `json:"cards"`
	HandType    PokerHandType `json:"hand_type"`
	HandRank    int           `json:"hand_rank"`
	Description string        `json:"description"`
}

// PokerDealResult represents the result of dealing final poker cards
type PokerDealResult struct {
	DealTimestamp   time.Time        `json:"deal_timestamp"`
	RandomSeed      string           `json:"random_seed"`
	Players         []PokerHand      `json:"players"`
	Winner          *PokerWinner     `json:"winner,omitempty"`
	PotDistribution *PotDistribution `json:"pot_distribution,omitempty"`
}

// PokerWinner represents the poker hand winner
type PokerWinner struct {
	PlayerID     string   `json:"player_id"`
	HandType     string   `json:"hand_type"`
	WinningCards []string `json:"winning_cards"`
}

// PotDistribution represents how the pot is distributed
type PotDistribution struct {
	TotalPot   float64         `json:"total_pot"`
	WinnerTake float64         `json:"winner_take"`
	Breakdown  *PotBreakdown   `json:"breakdown,omitempty"`
}

// PotBreakdown shows pot composition
type PotBreakdown struct {
	BaseBets          float64 `json:"base_bets"`
	PenaltyAdditions  float64 `json:"penalty_additions"`
}

// PuttPuttPokerCard represents a card earning/penalty event
type PuttPuttPokerCard struct {
	ID           string    `json:"id" db:"id"`
	PlayerID     string    `json:"player_id" db:"player_id"`
	GameID       string    `json:"game_id" db:"game_id"`
	Hole         *int      `json:"hole" db:"hole"`
	Action       string    `json:"action" db:"action"`
	CardsChange  int       `json:"cards_change" db:"cards_change"`
	PenaltyAmount *float64 `json:"penalty_amount" db:"penalty_amount"`
	TotalCards   int       `json:"total_cards" db:"total_cards"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// SideBetCalculation represents stored side bet calculation data
type SideBetCalculation struct {
	ID               string          `json:"id" db:"id"`
	GameID           string          `json:"game_id" db:"game_id"`
	PlayerID         string          `json:"player_id" db:"player_id"`
	BetType          SideBetType     `json:"bet_type" db:"bet_type"`
	CalculationData  json.RawMessage `json:"calculation_data" db:"calculation_data"`
	CurrentPosition  *int            `json:"current_position" db:"current_position"`
	FinalPosition    *int            `json:"final_position" db:"final_position"`
	IsWinner         bool            `json:"is_winner" db:"is_winner"`
	CalculatedAt     time.Time       `json:"calculated_at" db:"calculated_at"`
}

// BestNineCalculationData represents the stored calculation data for Best Nine
type BestNineCalculationData struct {
	BestHoles          []int   `json:"best_holes"`
	WorstHoles         []int   `json:"worst_holes"`
	RawScore           int     `json:"raw_score"`
	HandicapAdjustment int     `json:"handicap_adjustment"`
	FinalScore         int     `json:"final_score"`
	UsedScores         []HoleScore `json:"used_scores"`
}

// HoleScore represents a score used in calculations
type HoleScore struct {
	Hole           int  `json:"hole"`
	Strokes        int  `json:"strokes"`
	Par            int  `json:"par"`
	ScoreToPar     int  `json:"score_to_par"`
	HandicapStroke bool `json:"handicap_stroke"`
}

// PuttPuttPokerCalculationData represents stored calculation data for Putt Putt Poker
type PuttPuttPokerCalculationData struct {
	TotalCards     int              `json:"total_cards"`
	CardsEarned    int              `json:"cards_earned"`
	Penalties      int              `json:"penalties"`
	PuttingStats   PuttingStats     `json:"putting_stats"`
	CardHistory    []CardEvent      `json:"card_history"`
}

// CardEvent represents a card earning or penalty event
type CardEvent struct {
	Hole          *int     `json:"hole"`
	Action        string   `json:"action"`
	CardsChange   int      `json:"cards_change"`
	PenaltyAmount *float64 `json:"penalty_amount"`
	TotalCards    int      `json:"total_cards"`
}