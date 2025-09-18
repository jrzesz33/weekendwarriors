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

// PlayerHandler handles player-related HTTP requests
type PlayerHandler struct {
	playerService    *services.PlayerService
	websocketService *services.WebSocketService
}

// NewPlayerHandler creates a new player handler
func NewPlayerHandler(playerService *services.PlayerService, websocketService *services.WebSocketService) *PlayerHandler {
	return &PlayerHandler{
		playerService:    playerService,
		websocketService: websocketService,
	}
}

// GetPlayers handles GET /games/{gameId}/players
func (h *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	players, err := h.playerService.GetPlayers(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

// AddPlayer handles POST /games/{gameId}/players
func (h *PlayerHandler) AddPlayer(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	var req models.CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.New(errors.ErrValidation, "Invalid JSON in request body")
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	player, err := h.playerService.AddPlayer(gameID, &req)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Broadcast player added update
	h.websocketService.BroadcastGameUpdate(gameID, "player_added", player)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(player)

	log.Info().
		Str("player_id", player.ID).
		Str("game_id", gameID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Player added via API")
}

// GetPlayer handles GET /games/{gameId}/players/{playerId}
func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID
	playerID := chi.URLParam(r, "playerId")

	player, err := h.playerService.GetPlayer(gameID, playerID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)
}

// UpdatePlayer handles PUT /games/{gameId}/players/{playerId}
func (h *PlayerHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID
	playerID := chi.URLParam(r, "playerId")

	var req models.UpdatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiErr := errors.New(errors.ErrValidation, "Invalid JSON in request body")
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	player, err := h.playerService.UpdatePlayer(gameID, playerID, &req)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Broadcast player updated update
	h.websocketService.BroadcastGameUpdate(gameID, "player_updated", player)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(player)

	log.Info().
		Str("player_id", playerID).
		Str("game_id", gameID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Player updated via API")
}

// RemovePlayer handles DELETE /games/{gameId}/players/{playerId}
func (h *PlayerHandler) RemovePlayer(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID
	playerID := chi.URLParam(r, "playerId")

	err := h.playerService.RemovePlayer(gameID, playerID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Broadcast player removed update
	h.websocketService.BroadcastGameUpdate(gameID, "player_removed", map[string]string{
		"player_id": playerID,
	})

	w.WriteHeader(http.StatusNoContent)

	log.Info().
		Str("player_id", playerID).
		Str("game_id", gameID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Player removed via API")
}