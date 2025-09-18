package handlers

import (
	"encoding/json"
	"net/http"

	"golf-gamez/internal/middleware"
	"golf-gamez/internal/models"
	"golf-gamez/internal/services"
	"golf-gamez/pkg/auth"
	"golf-gamez/pkg/errors"

	"github.com/go-chi/chi/v5"
)

// SpectatorHandler handles spectator-related HTTP requests
type SpectatorHandler struct {
	gameService *services.GameService
}

// NewSpectatorHandler creates a new spectator handler
func NewSpectatorHandler(gameService *services.GameService) *SpectatorHandler {
	return &SpectatorHandler{
		gameService: gameService,
	}
}

// SpectateGame handles GET /spectate/{spectatorToken}
func (h *SpectatorHandler) SpectateGame(w http.ResponseWriter, r *http.Request) {
	spectatorToken := chi.URLParam(r, "spectatorToken")

	// Validate token format
	tokenType, err := auth.ValidateTokenFormat(spectatorToken)
	if err != nil || tokenType != auth.TokenTypeSpectator {
		apiErr := errors.New(errors.ErrInvalidToken, "Invalid spectator token")
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Get game using spectator token
	game, err := h.gameService.GetGame(spectatorToken)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// TODO: Get leaderboard and spectator count
	// For now, create a basic spectator view
	spectatorView := &models.SpectatorView{
		Game: game,
		Leaderboard: &models.Leaderboard{
			Overall: []models.LeaderboardEntry{},
		},
		LiveUpdates:    true,
		SpectatorCount: 1, // Placeholder
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spectatorView)
}