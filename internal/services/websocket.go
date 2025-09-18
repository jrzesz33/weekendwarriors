package services

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// WebSocketService handles real-time WebSocket connections
type WebSocketService struct {
	upgrader websocket.Upgrader
	clients  map[string]map[*websocket.Conn]bool // gameID -> connections
	mu       sync.RWMutex
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// TODO: Implement proper origin checking
				return true
			},
		},
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

// HandleWebSocket handles WebSocket connections for a game
func (s *WebSocketService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from URL path
	// This is a simplified implementation
	gameID := "placeholder" // TODO: Extract from URL

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade WebSocket connection")
		return
	}
	defer conn.Close()

	// Add client to game room
	s.addClient(gameID, conn)
	defer s.removeClient(gameID, conn)

	log.Info().
		Str("game_id", gameID).
		Str("remote_addr", r.RemoteAddr).
		Msg("WebSocket client connected")

	// Handle messages from client
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("WebSocket error")
			}
			break
		}

		// Handle incoming message
		s.handleMessage(gameID, conn, message)
	}

	log.Info().
		Str("game_id", gameID).
		Str("remote_addr", r.RemoteAddr).
		Msg("WebSocket client disconnected")
}

// BroadcastGameUpdate broadcasts a game update to all connected clients
func (s *WebSocketService) BroadcastGameUpdate(gameID string, updateType string, data interface{}) {
	message := WebSocketMessage{
		Type:   updateType,
		GameID: gameID,
		Data:   data,
	}

	s.broadcastToGame(gameID, message)
}

// BroadcastScoreUpdate broadcasts a score update
func (s *WebSocketService) BroadcastScoreUpdate(gameID, playerID string, hole int, score interface{}) {
	s.BroadcastGameUpdate(gameID, "score_update", map[string]interface{}{
		"player_id": playerID,
		"hole":      hole,
		"score":     score,
	})
}

// BroadcastLeaderboardUpdate broadcasts a leaderboard update
func (s *WebSocketService) BroadcastLeaderboardUpdate(gameID string, leaderboard interface{}) {
	s.BroadcastGameUpdate(gameID, "leaderboard_update", leaderboard)
}

// BroadcastSideBetUpdate broadcasts a side bet update
func (s *WebSocketService) BroadcastSideBetUpdate(gameID string, sideBetType string, data interface{}) {
	s.BroadcastGameUpdate(gameID, "side_bet_update", map[string]interface{}{
		"bet_type": sideBetType,
		"data":     data,
	})
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type   string      `json:"type"`
	GameID string      `json:"game_id"`
	Data   interface{} `json:"data"`
}

// Helper methods

func (s *WebSocketService) addClient(gameID string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.clients[gameID] == nil {
		s.clients[gameID] = make(map[*websocket.Conn]bool)
	}
	s.clients[gameID][conn] = true
}

func (s *WebSocketService) removeClient(gameID string, conn *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if clients, exists := s.clients[gameID]; exists {
		delete(clients, conn)
		if len(clients) == 0 {
			delete(s.clients, gameID)
		}
	}
}

func (s *WebSocketService) broadcastToGame(gameID string, message WebSocketMessage) {
	s.mu.RLock()
	clients := make(map[*websocket.Conn]bool)
	for conn, active := range s.clients[gameID] {
		if active {
			clients[conn] = true
		}
	}
	s.mu.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal WebSocket message")
		return
	}

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, messageBytes)
		if err != nil {
			log.Error().Err(err).Msg("Failed to write WebSocket message")
			s.removeClient(gameID, conn)
			conn.Close()
		}
	}
}

func (s *WebSocketService) handleMessage(gameID string, conn *websocket.Conn, message []byte) {
	// TODO: Implement message handling
	// This could handle authentication, subscriptions, etc.
	log.Debug().
		Str("game_id", gameID).
		Bytes("message", message).
		Msg("Received WebSocket message")
}