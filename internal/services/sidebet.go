package services

import (
	"database/sql"

	"golf-gamez/internal/models"
	"golf-gamez/pkg/errors"
)

// SideBetService handles side bet calculations and logic
type SideBetService struct {
	db *sql.DB
}

// NewSideBetService creates a new side bet service
func NewSideBetService(db *sql.DB) *SideBetService {
	return &SideBetService{db: db}
}

// GetBestNineStandings returns Best Nine side bet standings
func (s *SideBetService) GetBestNineStandings(gameID string) (*models.BestNineStandings, error) {
	// Verify game exists and has Best Nine enabled
	game, err := s.getGameForSideBet(gameID, models.SideBetBestNine)
	if err != nil {
		return nil, err
	}

	// TODO: Implement Best Nine calculation logic
	// For now, return empty standings
	standings := &models.BestNineStandings{
		BetType:         models.SideBetBestNine,
		Status:          game.Status,
		HandicapEnabled: game.HandicapEnabled,
		Standings:       []models.BestNineResult{},
	}

	return standings, nil
}

// GetPuttPuttPokerStatus returns Putt Putt Poker side bet status
func (s *SideBetService) GetPuttPuttPokerStatus(gameID string) (*models.PuttPuttPokerStatus, error) {
	// Verify game exists and has Putt Putt Poker enabled
	game, err := s.getGameForSideBet(gameID, models.SideBetPuttPuttPoker)
	if err != nil {
		return nil, err
	}

	// TODO: Implement Putt Putt Poker logic
	// For now, return empty status
	status := &models.PuttPuttPokerStatus{
		BetType:     models.SideBetPuttPuttPoker,
		Status:      game.Status,
		CurrentHole: game.CurrentHole,
		Players:     []models.PuttPuttPokerResult{},
	}

	return status, nil
}

// DealPokerCards deals final poker cards and determines winner
func (s *SideBetService) DealPokerCards(gameID string) (*models.PokerDealResult, error) {
	// Verify game is completed and has Putt Putt Poker enabled
	game, err := s.getGameForSideBet(gameID, models.SideBetPuttPuttPoker)
	if err != nil {
		return nil, err
	}

	if game.Status != models.GameStatusCompleted {
		return nil, errors.BusinessLogicError(
			errors.ErrGameNotCompleted,
			"Cannot deal final cards until game is completed",
			string(game.Status),
			string(models.GameStatusCompleted),
		)
	}

	// Check if cards have already been dealt
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM poker_hands WHERE game_id = ?", gameID).Scan(&count)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.New(errors.ErrCardsAlreadyDealt, "Final cards have already been dealt for this game")
	}

	// TODO: Implement poker card dealing logic
	// For now, return empty result
	result := &models.PokerDealResult{
		Players: []models.PokerHand{},
	}

	return result, nil
}

// UpdateSideBetsForScore updates side bet calculations when a score is recorded
func (s *SideBetService) UpdateSideBetsForScore(gameID, playerID string, score *models.Score) (*models.SideBetUpdates, error) {
	// TODO: Implement side bet updates
	// This would be called when a score is recorded to update side bet calculations

	updates := &models.SideBetUpdates{}

	// Check if Putt Putt Poker is enabled
	if s.isSideBetEnabled(gameID, models.SideBetPuttPuttPoker) {
		puttUpdate, err := s.updatePuttPuttPokerForScore(gameID, playerID, score)
		if err != nil {
			return nil, err
		}
		updates.PuttPuttPoker = puttUpdate
	}

	return updates, nil
}

// Helper methods

func (s *SideBetService) getGameForSideBet(gameID string, sideBetType models.SideBetType) (*gameInfo, error) {
	var game gameInfo
	var sideBetsJSON string

	query := `
		SELECT status, handicap_enabled, side_bets, current_hole
		FROM games
		WHERE id = ?
	`
	err := s.db.QueryRow(query, gameID).Scan(
		&game.Status,
		&game.HandicapEnabled,
		&sideBetsJSON,
		&game.CurrentHole,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ResourceNotFoundError("Game", gameID)
		}
		return nil, err
	}

	// Parse side bets
	tempGame := &models.Game{}
	if err := tempGame.UnmarshalSideBets(sideBetsJSON); err != nil {
		return nil, err
	}
	game.SideBets = tempGame.SideBets

	// Check if side bet is enabled
	enabled := false
	for _, sb := range game.SideBets {
		if sb == sideBetType {
			enabled = true
			break
		}
	}

	if !enabled {
		return nil, errors.NewWithDetails(
			errors.ErrSideBetNotEnabled,
			"Side bet is not enabled for this game",
			map[string]interface{}{
				"bet_type": sideBetType,
			},
		)
	}

	return &game, nil
}

func (s *SideBetService) isSideBetEnabled(gameID string, sideBetType models.SideBetType) bool {
	_, err := s.getGameForSideBet(gameID, sideBetType)
	return err == nil
}

func (s *SideBetService) updatePuttPuttPokerForScore(gameID, playerID string, score *models.Score) (*models.PuttPuttPokerUpdate, error) {
	// TODO: Implement Putt Putt Poker logic for score updates
	// This would check putts and award cards or apply penalties

	update := &models.PuttPuttPokerUpdate{
		CardsAwarded:   0,
		PenaltyApplied: false,
		TotalCards:     3, // Placeholder
	}

	return update, nil
}

// gameInfo is a simplified game struct for side bet operations
type gameInfo struct {
	Status          models.GameStatus
	HandicapEnabled bool
	SideBets        []models.SideBetType
	CurrentHole     *int
}