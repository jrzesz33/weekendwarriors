package handlers

import (
	"encoding/json"
	"net/http"

	"golf-gamez/internal/middleware"
	"golf-gamez/internal/services"
	"golf-gamez/pkg/errors"

	"github.com/rs/zerolog/log"
)

// SideBetHandler handles side bet related HTTP requests
type SideBetHandler struct {
	sideBetService   *services.SideBetService
	websocketService *services.WebSocketService
}

// NewSideBetHandler creates a new side bet handler
func NewSideBetHandler(sideBetService *services.SideBetService, websocketService *services.WebSocketService) *SideBetHandler {
	return &SideBetHandler{
		sideBetService:   sideBetService,
		websocketService: websocketService,
	}
}

// GetBestNineStandings handles GET /games/{gameId}/side-bets/best-nine
func (h *SideBetHandler) GetBestNineStandings(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	standings, err := h.sideBetService.GetBestNineStandings(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(standings)
}

// GetPuttPuttPokerStatus handles GET /games/{gameId}/side-bets/putt-putt-poker
func (h *SideBetHandler) GetPuttPuttPokerStatus(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	status, err := h.sideBetService.GetPuttPuttPokerStatus(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// DealPokerCards handles POST /games/{gameId}/side-bets/putt-putt-poker/deal
func (h *SideBetHandler) DealPokerCards(w http.ResponseWriter, r *http.Request) {
	authCtx, _ := middleware.GetGameAuthFromContext(r.Context())
	gameID := authCtx.GameID

	result, err := h.sideBetService.DealPokerCards(gameID)
	if err != nil {
		apiErr := errors.FromError(err)
		apiErr.RequestID = middleware.GetRequestID(r.Context())
		errors.WriteHTTPError(w, apiErr)
		return
	}

	// Broadcast poker cards dealt update
	h.websocketService.BroadcastSideBetUpdate(gameID, "putt_putt_poker", result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)

	log.Info().
		Str("game_id", gameID).
		Str("request_id", middleware.GetRequestID(r.Context())).
		Msg("Poker cards dealt via API")
}