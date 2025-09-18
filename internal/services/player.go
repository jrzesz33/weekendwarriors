package services

import (
	"database/sql"
	"fmt"
	"time"

	"golf-gamez/internal/models"
	"golf-gamez/pkg/auth"
	"golf-gamez/pkg/errors"

	"github.com/rs/zerolog/log"
)

// PlayerService handles player-related business logic
type PlayerService struct {
	db *sql.DB
}

// NewPlayerService creates a new player service
func NewPlayerService(db *sql.DB) *PlayerService {
	return &PlayerService{db: db}
}

// AddPlayer adds a new player to a game
func (s *PlayerService) AddPlayer(gameID string, req *models.CreatePlayerRequest) (*models.Player, error) {
	// Validate request
	if err := s.validateCreatePlayerRequest(req); err != nil {
		return nil, err
	}

	// Check game exists and is in setup state
	game, err := s.getGameForUpdate(gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != models.GameStatusSetup {
		return nil, errors.BusinessLogicError(
			errors.ErrInvalidGameState,
			"Cannot add players after game has started",
			string(game.Status),
			string(models.GameStatusSetup),
		)
	}

	// Check player limit (max 4 players)
	playerCount, err := s.getPlayerCount(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to get player count: %w", err)
	}

	if playerCount >= 4 {
		return nil, errors.New(errors.ErrPlayerLimitExceeded, "Maximum 4 players allowed per game")
	}

	// Check for duplicate name
	exists, err := s.playerNameExists(gameID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check player name: %w", err)
	}

	if exists {
		return nil, errors.New(errors.ErrDuplicatePlayerName, "Player name already exists in this game")
	}

	// Generate player ID
	playerID, err := auth.GeneratePlayerID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate player ID: %w", err)
	}

	// Determine position (next available position)
	position := playerCount + 1

	// Create player
	player := &models.Player{
		ID:        playerID,
		GameID:    gameID,
		Name:      req.Name,
		Handicap:  req.Handicap,
		Gender:    req.Gender,
		Position:  position,
		CreatedAt: time.Now(),
	}

	// Insert into database
	query := `
		INSERT INTO players (id, game_id, name, handicap, gender, position, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	var genderValue interface{}
	if player.Gender != nil {
		genderValue = string(*player.Gender)
	}

	_, err = s.db.Exec(
		query,
		player.ID,
		player.GameID,
		player.Name,
		player.Handicap,
		genderValue,
		player.Position,
		player.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert player: %w", err)
	}

	log.Info().
		Str("player_id", playerID).
		Str("game_id", gameID).
		Str("name", req.Name).
		Msg("Player added successfully")

	return player, nil
}

// GetPlayers retrieves all players for a game
func (s *PlayerService) GetPlayers(gameID string) (*models.PlayersResponse, error) {
	// Verify game exists
	if _, err := s.getGameForUpdate(gameID); err != nil {
		return nil, err
	}

	query := `
		SELECT id, game_id, name, handicap, gender, position, created_at
		FROM players
		WHERE game_id = ?
		ORDER BY position
	`
	rows, err := s.db.Query(query, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to query players: %w", err)
	}
	defer rows.Close()

	var players []models.Player
	for rows.Next() {
		var player models.Player
		var gender sql.NullString

		err := rows.Scan(
			&player.ID,
			&player.GameID,
			&player.Name,
			&player.Handicap,
			&gender,
			&player.Position,
			&player.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan player: %w", err)
		}

		if gender.Valid {
			g := models.Gender(gender.String)
			player.Gender = &g
		}

		// Load player stats
		stats, err := s.getPlayerStats(player.ID)
		if err != nil {
			log.Warn().Err(err).Str("player_id", player.ID).Msg("Failed to load player stats")
		} else {
			player.Stats = stats
		}

		players = append(players, player)
	}

	return &models.PlayersResponse{
		Players:    players,
		TotalCount: len(players),
	}, nil
}

// GetPlayer retrieves a specific player with detailed information
func (s *PlayerService) GetPlayer(gameID, playerID string) (*models.PlayerDetail, error) {
	// Verify game exists
	if _, err := s.getGameForUpdate(gameID); err != nil {
		return nil, err
	}

	query := `
		SELECT id, game_id, name, handicap, gender, position, created_at
		FROM players
		WHERE id = ? AND game_id = ?
	`

	var player models.Player
	var gender sql.NullString

	err := s.db.QueryRow(query, playerID, gameID).Scan(
		&player.ID,
		&player.GameID,
		&player.Name,
		&player.Handicap,
		&gender,
		&player.Position,
		&player.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ResourceNotFoundError("Player", playerID)
		}
		return nil, fmt.Errorf("failed to get player: %w", err)
	}

	if gender.Valid {
		g := models.Gender(gender.String)
		player.Gender = &g
	}

	// Load player stats
	stats, err := s.getPlayerStats(player.ID)
	if err != nil {
		log.Warn().Err(err).Str("player_id", player.ID).Msg("Failed to load player stats")
	} else {
		player.Stats = stats
	}

	// Load hole scores
	scores, err := s.getPlayerScores(playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to load player scores: %w", err)
	}

	return &models.PlayerDetail{
		Player:     player,
		HoleScores: scores,
	}, nil
}

// UpdatePlayer updates player information
func (s *PlayerService) UpdatePlayer(gameID, playerID string, req *models.UpdatePlayerRequest) (*models.Player, error) {
	// Validate request
	if err := s.validateUpdatePlayerRequest(req); err != nil {
		return nil, err
	}

	// Check game exists and player exists
	if _, err := s.getGameForUpdate(gameID); err != nil {
		return nil, err
	}

	// Get current player
	player, err := s.GetPlayer(gameID, playerID)
	if err != nil {
		return nil, err
	}

	// Check if name is being changed and if it would create a duplicate
	if req.Name != nil && *req.Name != player.Name {
		exists, err := s.playerNameExists(gameID, *req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check player name: %w", err)
		}
		if exists {
			return nil, errors.New(errors.ErrDuplicatePlayerName, "Player name already exists in this game")
		}
	}

	// Build update query dynamically
	setParts := []string{}
	args := []interface{}{}

	if req.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *req.Name)
	}

	if req.Handicap != nil {
		setParts = append(setParts, "handicap = ?")
		args = append(args, *req.Handicap)
	}

	if len(setParts) == 0 {
		// No updates needed, return current player
		return &player.Player, nil
	}

	// Add WHERE clause parameters
	args = append(args, playerID, gameID)

	query := fmt.Sprintf(
		"UPDATE players SET %s WHERE id = ? AND game_id = ?",
		fmt.Sprintf("%s", setParts[0]),
	)
	for i := 1; i < len(setParts); i++ {
		query = fmt.Sprintf("%s, %s", query[:len(query)-20], setParts[i]) + " WHERE id = ? AND game_id = ?"
	}

	_, err = s.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update player: %w", err)
	}

	// Return updated player
	updatedPlayer, err := s.GetPlayer(gameID, playerID)
	if err != nil {
		return nil, err
	}

	log.Info().
		Str("player_id", playerID).
		Str("game_id", gameID).
		Msg("Player updated successfully")

	return &updatedPlayer.Player, nil
}

// RemovePlayer removes a player from a game
func (s *PlayerService) RemovePlayer(gameID, playerID string) error {
	// Check game exists and is in setup state
	game, err := s.getGameForUpdate(gameID)
	if err != nil {
		return err
	}

	if game.Status != models.GameStatusSetup {
		return errors.BusinessLogicError(
			errors.ErrInvalidGameState,
			"Cannot remove players after game has started",
			string(game.Status),
			string(models.GameStatusSetup),
		)
	}

	// Verify player exists
	_, err = s.GetPlayer(gameID, playerID)
	if err != nil {
		return err
	}

	// Delete player (cascades to scores and side bet data)
	query := `DELETE FROM players WHERE id = ? AND game_id = ?`
	result, err := s.db.Exec(query, playerID, gameID)
	if err != nil {
		return fmt.Errorf("failed to delete player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ResourceNotFoundError("Player", playerID)
	}

	// Reorder remaining players' positions
	if err := s.reorderPlayerPositions(gameID); err != nil {
		log.Warn().Err(err).Msg("Failed to reorder player positions")
	}

	log.Info().
		Str("player_id", playerID).
		Str("game_id", gameID).
		Msg("Player removed successfully")

	return nil
}

// Helper methods

func (s *PlayerService) validateCreatePlayerRequest(req *models.CreatePlayerRequest) error {
	if req.Name == "" {
		return errors.ValidationError("name", "", "is required")
	}

	if len(req.Name) > 100 {
		return errors.ValidationError("name", req.Name, "must be 100 characters or less")
	}

	if req.Handicap < 0 || req.Handicap > 54 {
		return errors.ValidationError("handicap", fmt.Sprintf("%.1f", req.Handicap), "must be between 0 and 54")
	}

	if req.Gender != nil {
		if *req.Gender != models.GenderMale && *req.Gender != models.GenderFemale && *req.Gender != models.GenderOther {
			return errors.ValidationErrorWithAllowedValues(
				"gender",
				string(*req.Gender),
				[]interface{}{models.GenderMale, models.GenderFemale, models.GenderOther},
			)
		}
	}

	return nil
}

func (s *PlayerService) validateUpdatePlayerRequest(req *models.UpdatePlayerRequest) error {
	if req.Name != nil {
		if *req.Name == "" {
			return errors.ValidationError("name", "", "cannot be empty")
		}
		if len(*req.Name) > 100 {
			return errors.ValidationError("name", *req.Name, "must be 100 characters or less")
		}
	}

	if req.Handicap != nil {
		if *req.Handicap < 0 || *req.Handicap > 54 {
			return errors.ValidationError("handicap", fmt.Sprintf("%.1f", *req.Handicap), "must be between 0 and 54")
		}
	}

	return nil
}

func (s *PlayerService) getGameForUpdate(gameID string) (*models.Game, error) {
	var game models.Game
	query := `SELECT id, status FROM games WHERE id = ?`
	err := s.db.QueryRow(query, gameID).Scan(&game.ID, &game.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ResourceNotFoundError("Game", gameID)
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}
	return &game, nil
}

func (s *PlayerService) getPlayerCount(gameID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM players WHERE game_id = ?`
	err := s.db.QueryRow(query, gameID).Scan(&count)
	return count, err
}

func (s *PlayerService) playerNameExists(gameID, name string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM players WHERE game_id = ? AND name = ?`
	err := s.db.QueryRow(query, gameID, name).Scan(&count)
	return count > 0, err
}

func (s *PlayerService) getPlayerStats(playerID string) (*models.PlayerStats, error) {
	// This is a simplified version - would need more complex logic for full stats
	query := `
		SELECT COUNT(*) as holes_completed,
		       COALESCE(SUM(score_to_par), 0) as total_score,
		       COALESCE(SUM(putts), 0) as total_putts
		FROM scores
		WHERE player_id = ?
	`

	var holesCompleted int
	var totalScore int
	var totalPutts int

	err := s.db.QueryRow(query, playerID).Scan(&holesCompleted, &totalScore, &totalPutts)
	if err != nil {
		return nil, err
	}

	// Format current score
	var currentScore string
	if totalScore == 0 {
		currentScore = "E"
	} else if totalScore > 0 {
		currentScore = fmt.Sprintf("+%d", totalScore)
	} else {
		currentScore = fmt.Sprintf("%d", totalScore)
	}

	return &models.PlayerStats{
		HolesCompleted: holesCompleted,
		CurrentScore:   currentScore,
		TotalPutts:     totalPutts,
	}, nil
}

func (s *PlayerService) getPlayerScores(playerID string) ([]models.Score, error) {
	query := `
		SELECT id, player_id, game_id, hole, strokes, putts, par,
		       handicap_stroke, score_to_par, effective_score,
		       created_at, updated_at
		FROM scores
		WHERE player_id = ?
		ORDER BY hole
	`

	rows, err := s.db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []models.Score
	for rows.Next() {
		var score models.Score
		var updatedAt sql.NullTime

		err := rows.Scan(
			&score.ID,
			&score.PlayerID,
			&score.GameID,
			&score.Hole,
			&score.Strokes,
			&score.Putts,
			&score.Par,
			&score.HandicapStroke,
			&score.ScoreToPar,
			&score.EffectiveScore,
			&score.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		if updatedAt.Valid {
			score.UpdatedAt = &updatedAt.Time
		}

		// Format score to par
		score.ScoreToPar = models.FormatScoreToPar(score.Strokes, score.Par)

		scores = append(scores, score)
	}

	return scores, nil
}

func (s *PlayerService) reorderPlayerPositions(gameID string) error {
	// Get all players ordered by current position
	query := `
		SELECT id FROM players
		WHERE game_id = ?
		ORDER BY position, created_at
	`
	rows, err := s.db.Query(query, gameID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var playerIDs []string
	for rows.Next() {
		var playerID string
		if err := rows.Scan(&playerID); err != nil {
			return err
		}
		playerIDs = append(playerIDs, playerID)
	}

	// Update positions
	for i, playerID := range playerIDs {
		position := i + 1
		_, err := s.db.Exec(
			`UPDATE players SET position = ? WHERE id = ?`,
			position, playerID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}