package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golf-gamez/internal/middleware"
	"golf-gamez/internal/models"
	"golf-gamez/internal/services"
	"golf-gamez/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// ScoreHandler handles score-related HTTP requests
type ScoreHandler struct {
	scoreService     *services.ScoreService
	sideBetService   *services.SideBetService
	websocketService *services.WebSocketService
}

// NewScoreHandler creates a new score handler
func NewScoreHandler(scoreService *services.ScoreService, sideBetService *services.SideBetService, websocketService *services.WebSocketService) *ScoreHandler {
	return &ScoreHandler{
		scoreService:     scoreService,
		sideBetService:   sideBetService,
		websocketService: websocketService,
	}
}

// RecordScore handles POST /games/{gameId}/players/{playerId}/scores
func (h *ScoreHandler) RecordScore(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID
	playerID := chi.URLParam(r, "playerId")

	var req models.ScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.New(errors.ErrValidation, "Invalid JSON in request body")
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	score, err := h.scoreService.RecordScore(gameID, playerID, &req)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Update side bets
	sideBetUpdates, err := h.sideBetService.UpdateSideBetsForScore(gameID, playerID, score)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to update side bets for score")
	} else {
		score.SideBetUpdates = sideBetUpdates
	}

	// Broadcast score update
	h.websocketService.BroadcastScoreUpdate(gameID, playerID, req.Hole, score)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(score)

	log.Info().
		Str("score_id", score.ID).
		Str("player_id", playerID).
		Str("game_id", gameID).
		Int("hole", req.Hole).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Score recorded via API")
}

// UpdateScore handles PUT /games/{gameId}/players/{playerId}/scores/{hole}
func (h *ScoreHandler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID
	playerID := chi.URLParam(r, "playerId")
	holeStr := chi.URLParam(r, "hole")

	hole, err := strconv.Atoi(holeStr)
	if err != nil {
		apiErr := errors.ValidationError("hole", holeStr, "must be a valid integer")
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	var req models.UpdateScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.New(errors.ErrValidation, "Invalid JSON in request body")
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	score, err := h.scoreService.UpdateScore(gameID, playerID, hole, &req)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Update side bets
	sideBetUpdates, err := h.sideBetService.UpdateSideBetsForScore(gameID, playerID, score)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to update side bets for score")
	} else {
		score.SideBetUpdates = sideBetUpdates
	}

	// Broadcast score update
	h.websocketService.BroadcastScoreUpdate(gameID, playerID, hole, score)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(score)

	log.Info().
		Str("score_id", score.ID).
		Str("player_id", playerID).
		Str("game_id", gameID).
		Int("hole", hole).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Score updated via API")
}

// GetGameScorecard handles GET /games/{gameId}/scorecard
func (h *ScoreHandler) GetGameScorecard(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	scorecard, err := h.scoreService.GetGameScorecard(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scorecard)
}

// GetLeaderboard handles GET /games/{gameId}/leaderboard
func (h *ScoreHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	leaderboard, err := h.scoreService.GetLeaderboard(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}