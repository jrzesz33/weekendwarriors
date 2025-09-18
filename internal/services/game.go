package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"golf-gamez/internal/models"
	"golf-gamez/pkg/auth"
	"golf-gamez/pkg/errors"

	"github.com/rs/zerolog/log"
)

// GameService handles game-related business logic
type GameService struct {
	db *sql.DB
}

// NewGameService creates a new game service
func NewGameService(db *sql.DB) *GameService {
	return &GameService{db: db}
}

// CreateGame creates a new golf game
func (s *GameService) CreateGame(req *models.CreateGameRequest) (*models.Game, error) {
	// Generate IDs and tokens
	gameID, err := auth.GenerateGameID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate game ID: %w", err)
	}

	tokens, err := auth.GenerateTokenPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Validate course
	if req.Course != "diamond-run" {
		return nil, errors.ValidationError("course", req.Course, "must be diamond-run")
	}

	// Validate side bets
	for _, sideBet := range req.SideBets {
		if sideBet != models.SideBetBestNine && sideBet != models.SideBetPuttPuttPoker {
			return nil, errors.ValidationErrorWithAllowedValues(
				"side_bets",
				string(sideBet),
				[]interface{}{models.SideBetBestNine, models.SideBetPuttPuttPoker},
			)
		}
	}

	// Create game object
	game := &models.Game{
		ID:              gameID,
		Course:          req.Course,
		Status:          models.GameStatusSetup,
		HandicapEnabled: req.HandicapEnabled,
		SideBets:        req.SideBets,
		ShareToken:      tokens.ShareToken,
		SpectatorToken:  tokens.SpectatorToken,
		CreatedAt:       time.Now(),
	}

	// Marshal side bets
	sideBetsJSON, err := game.MarshalSideBets()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal side bets: %w", err)
	}

	// Insert into database
	query := `
		INSERT INTO games (
			id, course, status, handicap_enabled, side_bets,
			share_token, spectator_token, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.Exec(
		query,
		game.ID,
		game.Course,
		game.Status,
		game.HandicapEnabled,
		sideBetsJSON,
		game.ShareToken,
		game.SpectatorToken,
		game.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert game: %w", err)
	}

	// Set share and spectator links
	game.ShareLink = fmt.Sprintf("/games/%s", tokens.ShareToken)
	game.SpectatorLink = fmt.Sprintf("/spectate/%s", tokens.SpectatorToken)

	// Load course info
	courseInfo, err := s.getCourseInfo(req.Course)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load course info")
	} else {
		game.CourseInfo = courseInfo
	}

	log.Info().Str("game_id", gameID).Msg("Game created successfully")
	return game, nil
}

// GetGame retrieves a game by ID or token
func (s *GameService) GetGame(gameIDOrToken string) (*models.Game, error) {
	// Try to determine if it's a token or game ID
	var query string
	var param string

	if gameIDOrToken[:3] == "gt_" {
		query = `
			SELECT id, course, status, handicap_enabled, side_bets,
			       share_token, spectator_token, current_hole,
			       created_at, started_at, completed_at, final_results
			FROM games WHERE share_token = ?
		`
		param = gameIDOrToken
	} else if gameIDOrToken[:3] == "st_" {
		query = `
			SELECT id, course, status, handicap_enabled, side_bets,
			       share_token, spectator_token, current_hole,
			       created_at, started_at, completed_at, final_results
			FROM games WHERE spectator_token = ?
		`
		param = gameIDOrToken
	} else {
		query = `
			SELECT id, course, status, handicap_enabled, side_bets,
			       share_token, spectator_token, current_hole,
			       created_at, started_at, completed_at, final_results
			FROM games WHERE id = ?
		`
		param = gameIDOrToken
	}

	var game models.Game
	var sideBetsJSON string
	var finalResultsJSON sql.NullString

	err := s.db.QueryRow(query, param).Scan(
		&game.ID,
		&game.Course,
		&game.Status,
		&game.HandicapEnabled,
		&sideBetsJSON,
		&game.ShareToken,
		&game.SpectatorToken,
		&game.CurrentHole,
		&game.CreatedAt,
		&game.StartedAt,
		&game.CompletedAt,
		&finalResultsJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ResourceNotFoundError("Game", gameIDOrToken)
		}
		return nil, fmt.Errorf("failed to get game: %w", err)
	}

	// Unmarshal side bets
	if err := game.UnmarshalSideBets(sideBetsJSON); err != nil {
		return nil, fmt.Errorf("failed to unmarshal side bets: %w", err)
	}

	// Unmarshal final results
	if finalResultsJSON.Valid {
		if err := game.UnmarshalFinalResults(finalResultsJSON.String); err != nil {
			return nil, fmt.Errorf("failed to unmarshal final results: %w", err)
		}
	}

	// Set links
	game.ShareLink = fmt.Sprintf("/games/%s", game.ShareToken)
	game.SpectatorLink = fmt.Sprintf("/spectate/%s", game.SpectatorToken)

	// Load players
	players, err := s.getGamePlayers(game.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load players: %w", err)
	}
	game.Players = players

	// Load course info
	courseInfo, err := s.getCourseInfo(game.Course)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load course info")
	} else {
		game.CourseInfo = courseInfo
	}

	return &game, nil
}

// StartGame starts a golf game
func (s *GameService) StartGame(gameID string) (*models.Game, error) {
	// Check current game state
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != models.GameStatusSetup {
		return nil, errors.BusinessLogicError(
			errors.ErrInvalidGameState,
			"Cannot start game",
			string(game.Status),
			string(models.GameStatusSetup),
		)
	}

	// Check if game has players
	if len(game.Players) == 0 {
		return nil, errors.New(errors.ErrInvalidGameState, "Cannot start game without players")
	}

	// Update game status
	now := time.Now()
	query := `
		UPDATE games
		SET status = ?, started_at = ?, current_hole = 1
		WHERE id = ?
	`
	_, err = s.db.Exec(query, models.GameStatusInProgress, now, gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}

	// Initialize side bet calculations for each player
	if err := s.initializeSideBets(gameID, game.Players, game.SideBets); err != nil {
		log.Warn().Err(err).Msg("Failed to initialize side bets")
	}

	// Return updated game
	return s.GetGame(gameID)
}

// CompleteGame marks a game as completed and calculates final results
func (s *GameService) CompleteGame(gameID string) (*models.GameCompletionResult, error) {
	game, err := s.GetGame(gameID)
	if err != nil {
		return nil, err
	}

	if game.Status != models.GameStatusInProgress {
		return nil, errors.BusinessLogicError(
			errors.ErrInvalidGameState,
			"Cannot complete game",
			string(game.Status),
			string(models.GameStatusInProgress),
		)
	}

	// Calculate final results
	finalResults, err := s.calculateFinalResults(gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate final results: %w", err)
	}

	// Marshal final results
	finalResultsJSON, err := json.Marshal(finalResults)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal final results: %w", err)
	}

	// Update game
	now := time.Now()
	query := `
		UPDATE games
		SET status = ?, completed_at = ?, final_results = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(query, models.GameStatusCompleted, now, string(finalResultsJSON), gameID)
	if err != nil {
		return nil, fmt.Errorf("failed to complete game: %w", err)
	}

	return &models.GameCompletionResult{
		ID:           gameID,
		Status:       models.GameStatusCompleted,
		CompletedAt:  now,
		FinalResults: finalResults,
	}, nil
}

// DeleteGame deletes a game and all associated data
func (s *GameService) DeleteGame(gameID string) error {
	// Verify game exists
	_, err := s.GetGame(gameID)
	if err != nil {
		return err
	}

	// Delete game (cascades to all related tables)
	query := `DELETE FROM games WHERE id = ?`
	result, err := s.db.Exec(query, gameID)
	if err != nil {
		return fmt.Errorf("failed to delete game: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.ResourceNotFoundError("Game", gameID)
	}

	log.Info().Str("game_id", gameID).Msg("Game deleted successfully")
	return nil
}

// getGamePlayers loads players for a game
func (s *GameService) getGamePlayers(gameID string) ([]models.Player, error) {
	query := `
		SELECT id, game_id, name, handicap, gender, position, created_at
		FROM players
		WHERE game_id = ?
		ORDER BY position
	`
	rows, err := s.db.Query(query, gameID)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		if gender.Valid {
			g := models.Gender(gender.String)
			player.Gender = &g
		}

		players = append(players, player)
	}

	return players, nil
}

// getCourseInfo loads course information
func (s *GameService) getCourseInfo(courseName string) (*models.CourseInfo, error) {
	query := `
		SELECT hole, par, handicap_ranking, yardage, description
		FROM course_data
		WHERE course_name = ?
		ORDER BY hole
	`
	rows, err := s.db.Query(query, courseName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var holes []models.HoleInfo
	totalPar := 0

	for rows.Next() {
		var hole models.HoleInfo
		var yardage sql.NullInt64
		var description sql.NullString

		err := rows.Scan(
			&hole.Hole,
			&hole.Par,
			&hole.HandicapRanking,
			&yardage,
			&description,
		)
		if err != nil {
			return nil, err
		}

		if yardage.Valid {
			y := int(yardage.Int64)
			hole.Yardage = &y
		}

		if description.Valid {
			hole.Description = description.String
		}

		holes = append(holes, hole)
		totalPar += hole.Par
	}

	return &models.CourseInfo{
		Name:     "Diamond Run",
		Holes:    holes,
		TotalPar: totalPar,
	}, nil
}

// initializeSideBets creates initial side bet calculations for players
func (s *GameService) initializeSideBets(gameID string, players []models.Player, sideBets []models.SideBetType) error {
	for _, player := range players {
		for _, sideBet := range sideBets {
			sideBetID, err := auth.GenerateSideBetID()
			if err != nil {
				return err
			}

			var calculationData []byte
			switch sideBet {
			case models.SideBetBestNine:
				data := models.BestNineCalculationData{
					BestHoles:          []int{},
					WorstHoles:         []int{},
					RawScore:           0,
					HandicapAdjustment: 0,
					FinalScore:         0,
					UsedScores:         []models.HoleScore{},
				}
				calculationData, err = json.Marshal(data)
			case models.SideBetPuttPuttPoker:
				data := models.PuttPuttPokerCalculationData{
					TotalCards:   3, // Starting cards
					CardsEarned:  0,
					Penalties:    0,
					PuttingStats: models.PuttingStats{},
					CardHistory:  []models.CardEvent{},
				}
				calculationData, err = json.Marshal(data)

				// Also create initial poker card record
				if err == nil {
					cardID, _ := auth.GeneratePokerCardID()
					_, err = s.db.Exec(`
						INSERT INTO putt_putt_poker_cards
						(id, player_id, game_id, action, cards_change, total_cards)
						VALUES (?, ?, ?, 'starting', 3, 3)
					`, cardID, player.ID, gameID)
				}
			}

			if err != nil {
				return err
			}

			// Insert side bet calculation
			_, err = s.db.Exec(`
				INSERT INTO side_bet_calculations
				(id, game_id, player_id, bet_type, calculation_data)
				VALUES (?, ?, ?, ?, ?)
			`, sideBetID, gameID, player.ID, sideBet, string(calculationData))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// calculateFinalResults calculates final game results
func (s *GameService) calculateFinalResults(gameID string) (*models.FinalResults, error) {
	// TODO: Implement final results calculation
	// This would calculate overall winner, best nine winner, and poker winner
	return &models.FinalResults{}, nil
}