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

// ScoreService handles score-related business logic
type ScoreService struct {
	db *sql.DB
}

// NewScoreService creates a new score service
func NewScoreService(db *sql.DB) *ScoreService {
	return &ScoreService{db: db}
}

// RecordScore records a score for a player on a specific hole
func (s *ScoreService) RecordScore(gameID, playerID string, req *models.ScoreRequest) (*models.Score, error) {
	// Validate request
	if err := s.validateScoreRequest(req); err != nil {
		return nil, err
	}

	// Validate business rules
	if err := s.validateScoreBusinessRules(gameID, playerID, req); err != nil {
		return nil, err
	}

	// Get hole par from course data
	par, err := s.getHolePar(gameID, req.Hole)
	if err != nil {
		return nil, fmt.Errorf("failed to get hole par: %w", err)
	}

	// Calculate handicap stroke (simplified logic for now)
	handicapStroke := s.calculateHandicapStroke(gameID, playerID, req.Hole)

	// Calculate scores
	scoreToPar := req.Strokes - par
	effectiveScore := models.CalculateEffectiveScore(req.Strokes, par, handicapStroke)

	// Generate score ID
	scoreID, err := auth.GenerateScoreID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate score ID: %w", err)
	}

	// Create score
	score := &models.Score{
		ID:             scoreID,
		PlayerID:       playerID,
		GameID:         gameID,
		Hole:           req.Hole,
		Strokes:        req.Strokes,
		Putts:          req.Putts,
		Par:            par,
		HandicapStroke: handicapStroke,
		EffectiveScore: effectiveScore,
		CreatedAt:      time.Now(),
	}

	// Format score to par
	score.ScoreToPar = models.FormatScoreToPar(req.Strokes, par)

	// Insert into database
	query := `
		INSERT INTO scores (
			id, player_id, game_id, hole, strokes, putts, par,
			handicap_stroke, score_to_par, effective_score, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = s.db.Exec(
		query,
		score.ID,
		score.PlayerID,
		score.GameID,
		score.Hole,
		score.Strokes,
		score.Putts,
		score.Par,
		score.HandicapStroke,
		scoreToPar,
		score.EffectiveScore,
		score.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert score: %w", err)
	}

	log.Info().
		Str("score_id", scoreID).
		Str("player_id", playerID).
		Str("game_id", gameID).
		Int("hole", req.Hole).
		Int("strokes", req.Strokes).
		Msg("Score recorded successfully")

	return score, nil
}

// UpdateScore updates an existing score
func (s *ScoreService) UpdateScore(gameID, playerID string, hole int, req *models.UpdateScoreRequest) (*models.Score, error) {
	// Get existing score
	existingScore, err := s.getScore(gameID, playerID, hole)
	if err != nil {
		return nil, err
	}

	// Validate update request
	if err := s.validateUpdateScoreRequest(req); err != nil {
		return nil, err
	}

	// Build update values
	strokes := existingScore.Strokes
	putts := existingScore.Putts

	if req.Strokes != nil {
		strokes = *req.Strokes
	}
	if req.Putts != nil {
		putts = *req.Putts
	}

	// Validate putts vs strokes
	if putts > strokes {
		return nil, errors.NewWithDetails(
			errors.ErrInvalidPuttCount,
			"Putt count cannot exceed stroke count",
			map[string]interface{}{
				"strokes": strokes,
				"putts":   putts,
			},
		)
	}

	// Recalculate derived values
	scoreToPar := strokes - existingScore.Par
	effectiveScore := models.CalculateEffectiveScore(strokes, existingScore.Par, existingScore.HandicapStroke)

	// Update database
	now := time.Now()
	query := `
		UPDATE scores
		SET strokes = ?, putts = ?, score_to_par = ?, effective_score = ?, updated_at = ?
		WHERE id = ?
	`
	_, err = s.db.Exec(query, strokes, putts, scoreToPar, effectiveScore, now, existingScore.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update score: %w", err)
	}

	// Return updated score
	return s.getScore(gameID, playerID, hole)
}

// GetGameScorecard returns the complete scorecard for a game
func (s *ScoreService) GetGameScorecard(gameID string) (*models.GameScorecard, error) {
	// Get game info
	game, err := s.getGameSummary(gameID)
	if err != nil {
		return nil, err
	}

	// Get course info
	courseInfo, err := s.getCourseHoles(game.Course)
	if err != nil {
		return nil, err
	}

	// Get players with scores
	players, err := s.getPlayersWithScores(gameID)
	if err != nil {
		return nil, err
	}

	return &models.GameScorecard{
		Game:       *game,
		CourseInfo: courseInfo,
		Players:    players,
	}, nil
}

// GetLeaderboard returns the current game leaderboard
func (s *ScoreService) GetLeaderboard(gameID string) (*models.Leaderboard, error) {
	// Get overall leaderboard
	overall, err := s.getOverallLeaderboard(gameID)
	if err != nil {
		return nil, err
	}

	leaderboard := &models.Leaderboard{
		Overall: overall,
	}

	// TODO: Add side bet leaderboards
	// This would be implemented when side bet services are complete

	return leaderboard, nil
}

// Helper methods

func (s *ScoreService) validateScoreRequest(req *models.ScoreRequest) error {
	if req.Hole < 1 || req.Hole > 18 {
		return errors.ValidationError("hole", fmt.Sprintf("%d", req.Hole), "must be between 1 and 18")
	}

	if req.Strokes < 1 || req.Strokes > 20 {
		return errors.ValidationError("strokes", fmt.Sprintf("%d", req.Strokes), "must be between 1 and 20")
	}

	if req.Putts < 0 || req.Putts > 10 {
		return errors.ValidationError("putts", fmt.Sprintf("%d", req.Putts), "must be between 0 and 10")
	}

	if req.Putts > req.Strokes {
		return errors.NewWithDetails(
			errors.ErrInvalidPuttCount,
			"Putt count cannot exceed stroke count",
			map[string]interface{}{
				"strokes": req.Strokes,
				"putts":   req.Putts,
			},
		)
	}

	return nil
}

func (s *ScoreService) validateUpdateScoreRequest(req *models.UpdateScoreRequest) error {
	if req.Strokes != nil {
		if *req.Strokes < 1 || *req.Strokes > 20 {
			return errors.ValidationError("strokes", fmt.Sprintf("%d", *req.Strokes), "must be between 1 and 20")
		}
	}

	if req.Putts != nil {
		if *req.Putts < 0 || *req.Putts > 10 {
			return errors.ValidationError("putts", fmt.Sprintf("%d", *req.Putts), "must be between 0 and 10")
		}
	}

	return nil
}

func (s *ScoreService) validateScoreBusinessRules(gameID, playerID string, req *models.ScoreRequest) error {
	// Check game exists and is in progress
	var status string
	err := s.db.QueryRow("SELECT status FROM games WHERE id = ?", gameID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.ResourceNotFoundError("Game", gameID)
		}
		return err
	}

	if status != string(models.GameStatusInProgress) {
		return errors.BusinessLogicError(
			errors.ErrGameNotStarted,
			"Cannot record scores for game that is not in progress",
			status,
			string(models.GameStatusInProgress),
		)
	}

	// Check player exists in game
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM players WHERE id = ? AND game_id = ?", playerID, gameID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.ResourceNotFoundError("Player", playerID)
	}

	// Check if score already exists
	err = s.db.QueryRow("SELECT COUNT(*) FROM scores WHERE player_id = ? AND hole = ?", playerID, req.Hole).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New(errors.ErrScoreAlreadyExists, "Score for this hole has already been recorded")
	}

	// TODO: Add validation for sequential hole scoring if needed

	return nil
}

func (s *ScoreService) getHolePar(gameID string, hole int) (int, error) {
	// Get course name first
	var course string
	err := s.db.QueryRow("SELECT course FROM games WHERE id = ?", gameID).Scan(&course)
	if err != nil {
		return 0, err
	}

	// Get hole par
	var par int
	err = s.db.QueryRow("SELECT par FROM course_data WHERE course_name = ? AND hole = ?", course, hole).Scan(&par)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.ValidationError("hole", fmt.Sprintf("%d", hole), "invalid hole number for course")
		}
		return 0, err
	}

	return par, nil
}

func (s *ScoreService) calculateHandicapStroke(gameID, playerID string, hole int) bool {
	// Simplified handicap calculation - would need more complex logic in production
	// This is a placeholder that returns false for now
	return false
}

func (s *ScoreService) getScore(gameID, playerID string, hole int) (*models.Score, error) {
	query := `
		SELECT id, player_id, game_id, hole, strokes, putts, par,
		       handicap_stroke, score_to_par, effective_score,
		       created_at, updated_at
		FROM scores
		WHERE game_id = ? AND player_id = ? AND hole = ?
	`

	var score models.Score
	var updatedAt sql.NullTime

	err := s.db.QueryRow(query, gameID, playerID, hole).Scan(
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
		if err == sql.ErrNoRows {
			return nil, errors.ResourceNotFoundError("Score", fmt.Sprintf("hole %d", hole))
		}
		return nil, err
	}

	if updatedAt.Valid {
		score.UpdatedAt = &updatedAt.Time
	}

	// Format score to par
	score.ScoreToPar = models.FormatScoreToPar(score.Strokes, score.Par)

	return &score, nil
}

func (s *ScoreService) getGameSummary(gameID string) (*models.GameSummary, error) {
	var game models.GameSummary
	query := `SELECT id, course, current_hole FROM games WHERE id = ?`
	err := s.db.QueryRow(query, gameID).Scan(&game.ID, &game.Course, &game.CurrentHole)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ResourceNotFoundError("Game", gameID)
		}
		return nil, err
	}
	return &game, nil
}

func (s *ScoreService) getCourseHoles(courseName string) ([]models.HoleInfo, error) {
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
	}

	return holes, nil
}

func (s *ScoreService) getPlayersWithScores(gameID string) ([]models.ScorecardPlayer, error) {
	// Get players
	playersQuery := `
		SELECT id, name, position
		FROM players
		WHERE game_id = ?
		ORDER BY position
	`
	rows, err := s.db.Query(playersQuery, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var players []models.ScorecardPlayer
	for rows.Next() {
		var player models.ScorecardPlayer
		err := rows.Scan(&player.ID, &player.Name, &player.Position)
		if err != nil {
			return nil, err
		}

		// Get scores for this player
		scores, err := s.getPlayerScores(player.ID)
		if err != nil {
			return nil, err
		}
		player.Scores = scores

		// TODO: Calculate totals
		// player.Totals = calculatePlayerTotals(scores)

		players = append(players, player)
	}

	return players, nil
}

func (s *ScoreService) getPlayerScores(playerID string) ([]models.Score, error) {
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

		score.ScoreToPar = models.FormatScoreToPar(score.Strokes, score.Par)
		scores = append(scores, score)
	}

	return scores, nil
}

func (s *ScoreService) getOverallLeaderboard(gameID string) ([]models.LeaderboardEntry, error) {
	query := `
		SELECT p.id, p.name, p.handicap,
		       COUNT(s.id) as holes_completed,
		       COALESCE(SUM(s.score_to_par), 0) as total_score,
		       COALESCE(SUM(s.putts), 0) as total_putts
		FROM players p
		LEFT JOIN scores s ON p.id = s.player_id
		WHERE p.game_id = ?
		GROUP BY p.id, p.name, p.handicap
		ORDER BY total_score ASC, holes_completed DESC
	`

	rows, err := s.db.Query(query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	position := 1

	for rows.Next() {
		var entry models.LeaderboardEntry
		var totalScore int
		var handicap sql.NullFloat64

		err := rows.Scan(
			&entry.Player.ID,
			&entry.Player.Name,
			&handicap,
			&entry.HolesCompleted,
			&totalScore,
			&entry.TotalPutts,
		)
		if err != nil {
			return nil, err
		}

		if handicap.Valid {
			h := handicap.Float64
			entry.Player.Handicap = &h
		}

		entry.Position = position
		entry.Score = models.FormatScoreToPar(totalScore, 0) // Format relative to par

		entries = append(entries, entry)
		position++
	}

	return entries, nil
}