package handlers

import (
	"encoding/json"
	"net/http"

	"golf-gamez/internal/middleware"
	"golf-gamez/internal/models"
	"golf-gamez/internal/services"
	"golf-gamez/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// GameHandler handles game-related HTTP requests
type GameHandler struct {
	gameService      *services.GameService
	websocketService *services.WebSocketService
}

// NewGameHandler creates a new game handler
func NewGameHandler(gameService *services.GameService, websocketService *services.WebSocketService) *GameHandler {
	return &GameHandler{
		gameService:      gameService,
		websocketService: websocketService,
	}
}

// CreateGame handles POST /games
func (h *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var req models.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.New(errors.ErrValidation, "Invalid JSON in request body")
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	game, err := h.gameService.CreateGame(&req)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(game)

	log.Info().
		Str("game_id", game.ID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Game created via API")
}

// GetGame handles GET /games/{gameId}
func (h *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	gameID := chi.URLParam(r, "gameId")

	// Get game auth context to determine which game ID to use
	authCtx, ok := middleware.GetGameAuthFromContext(r.Context())
	if ok {
		gameID = authCtx.GameID
	}

	game, err := h.gameService.GetGame(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

// StartGame handles POST /games/{gameId}/start
func (h *GameHandler) StartGame(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	game, err := h.gameService.StartGame(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Broadcast game started update
	h.websocketService.BroadcastGameUpdate(gameID, "game_started", game)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     game.ID,
		"status": game.Status,
	})

	log.Info().
		Str("game_id", gameID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Game started via API")
}

// CompleteGame handles POST /games/{gameId}/complete
func (h *GameHandler) CompleteGame(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	result, err := h.gameService.CompleteGame(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Broadcast game completed update
	h.websocketService.BroadcastGameUpdate(gameID, "game_completed", result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

	log.Info().
		Str("game_id", gameID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Game completed via API")
}

// DeleteGame handles DELETE /games/{gameId}
func (h *GameHandler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	err := h.gameService.DeleteGame(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Broadcast game deleted update
	h.websocketService.BroadcastGameUpdate(gameID, "game_deleted", map[string]string{
		"game_id": gameID,
	})

	w.WriteHeader(http.StatusNoContent)

	log.Info().
		Str("game_id", gameID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Game deleted via API")
}