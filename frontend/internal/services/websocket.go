//go:build js && wasm
// +build js,wasm

package services

import (
	"encoding/json"
	"fmt"
	"golf-gamez-frontend/internal/models"
	"log"
	"syscall/js"
	"time"
)

// WSMessageHandler is a function type for handling WebSocket messages
type WSMessageHandler func(msgType string, data interface{})

// WSService manages WebSocket connections for real-time updates
type WSService struct {
	ws           js.Value
	connected    bool
	gameID       string
	messageQueue []models.WSMessage
	handlers     map[string][]WSMessageHandler
	reconnectAttempts int
	maxReconnectAttempts int
}

// NewWSService creates a new WebSocket service
func NewWSService() *WSService {
	return &WSService{
		handlers:             make(map[string][]WSMessageHandler),
		maxReconnectAttempts: 5,
	}
}

// Connect establishes a WebSocket connection to the game
func (ws *WSService) Connect(gameID string) error {
	if ws.connected {
		ws.Disconnect()
	}

	ws.gameID = gameID

	// Construct WebSocket URL
	wsURL := fmt.Sprintf("ws://localhost:8080/v1/ws/games/%s", gameID)

	// Create WebSocket connection using JavaScript WebSocket API
	wsConstructor := js.Global().Get("WebSocket")
	if wsConstructor.IsUndefined() {
		return fmt.Errorf("WebSocket not supported in this environment")
	}

	ws.ws = wsConstructor.New(wsURL)

	// Set up event handlers
	ws.setupEventHandlers()

	return nil
}

// setupEventHandlers configures WebSocket event handlers
func (ws *WSService) setupEventHandlers() {
	// OnOpen handler
	openHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.connected = true
		ws.reconnectAttempts = 0
		log.Println("WebSocket connected")

		// Process any queued messages
		ws.processMessageQueue()

		// Notify handlers about connection
		ws.notifyHandlers("connection", map[string]interface{}{
			"status": "connected",
		})

		return nil
	})

	// OnMessage handler
	messageHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			messageEvent := args[0]
			data := messageEvent.Get("data").String()

			var message models.WSMessage
			if err := json.Unmarshal([]byte(data), &message); err != nil {
				log.Printf("Error parsing WebSocket message: %v", err)
				return nil
			}

			// Notify all handlers for this message type
			ws.notifyHandlers(message.Type, message.Data)
		}
		return nil
	})

	// OnError handler
	errorHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		log.Println("WebSocket error occurred")
		ws.connected = false

		ws.notifyHandlers("error", map[string]interface{}{
			"status": "error",
		})

		// Attempt to reconnect
		go ws.attemptReconnect()

		return nil
	})

	// OnClose handler
	closeHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ws.connected = false
		log.Println("WebSocket connection closed")

		ws.notifyHandlers("connection", map[string]interface{}{
			"status": "disconnected",
		})

		// Attempt to reconnect if not intentionally closed
		if ws.reconnectAttempts < ws.maxReconnectAttempts {
			go ws.attemptReconnect()
		}

		return nil
	})

	// Assign event handlers
	ws.ws.Set("onopen", openHandler)
	ws.ws.Set("onmessage", messageHandler)
	ws.ws.Set("onerror", errorHandler)
	ws.ws.Set("onclose", closeHandler)
}

// attemptReconnect tries to reconnect to the WebSocket
func (ws *WSService) attemptReconnect() {
	if ws.reconnectAttempts >= ws.maxReconnectAttempts {
		log.Println("Max reconnection attempts reached")
		return
	}

	ws.reconnectAttempts++

	// Wait before attempting reconnection
	time.Sleep(time.Duration(ws.reconnectAttempts) * time.Second)

	log.Printf("Attempting to reconnect (attempt %d/%d)", ws.reconnectAttempts, ws.maxReconnectAttempts)

	if err := ws.Connect(ws.gameID); err != nil {
		log.Printf("Reconnection attempt failed: %v", err)
	}
}

// AddHandler adds a message handler for a specific message type
func (ws *WSService) AddHandler(messageType string, handler WSMessageHandler) {
	if ws.handlers[messageType] == nil {
		ws.handlers[messageType] = make([]WSMessageHandler, 0)
	}
	ws.handlers[messageType] = append(ws.handlers[messageType], handler)
}

// RemoveHandler removes a message handler (simplified - removes all handlers for the type)
func (ws *WSService) RemoveHandler(messageType string) {
	delete(ws.handlers, messageType)
}

// notifyHandlers calls all registered handlers for a message type
func (ws *WSService) notifyHandlers(messageType string, data interface{}) {
	if handlers, exists := ws.handlers[messageType]; exists {
		for _, handler := range handlers {
			go handler(messageType, data)
		}
	}
}

// SendMessage sends a message through the WebSocket connection
func (ws *WSService) SendMessage(message models.WSMessage) error {
	if !ws.connected {
		// Queue the message for later sending
		ws.messageQueue = append(ws.messageQueue, message)
		return fmt.Errorf("WebSocket not connected, message queued")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ws.ws.Call("send", string(data))
	return nil
}

// processMessageQueue sends all queued messages
func (ws *WSService) processMessageQueue() {
	for _, message := range ws.messageQueue {
		if err := ws.SendMessage(message); err != nil {
			log.Printf("Failed to send queued message: %v", err)
		}
	}
	ws.messageQueue = ws.messageQueue[:0] // Clear the queue
}

// IsConnected returns the current connection status
func (ws *WSService) IsConnected() bool {
	return ws.connected
}

// Disconnect closes the WebSocket connection
func (ws *WSService) Disconnect() {
	if ws.ws.IsUndefined() {
		return
	}

	ws.connected = false
	ws.reconnectAttempts = ws.maxReconnectAttempts // Prevent reconnection attempts
	ws.ws.Call("close")
}

// GetGameID returns the current game ID
func (ws *WSService) GetGameID() string {
	return ws.gameID
}