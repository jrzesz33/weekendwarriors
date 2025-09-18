package websocket

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"

	"golf-gamez-frontend/internal/models"
)

// ConnectionState represents the current WebSocket connection state
type ConnectionState int

const (
	Disconnected ConnectionState = iota
	Connecting
	Connected
	Reconnecting
)

// Manager handles WebSocket connections with resilient retry logic
type Manager struct {
	gameID            string
	websocket         js.Value
	state             ConnectionState
	retryCount        int
	maxRetries        int
	baseRetryDelay    time.Duration
	maxRetryDelay     time.Duration
	heartbeatInterval time.Duration
	lastHeartbeat     time.Time
	onMessage         func(*models.WebSocketMessage)
	onStateChange     func(ConnectionState)
	messageQueue      []*models.WebSocketMessage
	autoReconnect     bool
	reconnectTimer    *time.Timer
	heartbeatTimer    *time.Timer
}

// NewManager creates a new WebSocket manager
func NewManager() *Manager {
	return &Manager{
		state:             Disconnected,
		maxRetries:        10,
		baseRetryDelay:    1 * time.Second,
		maxRetryDelay:     30 * time.Second,
		heartbeatInterval: 30 * time.Second,
		autoReconnect:     true,
		messageQueue:      make([]*models.WebSocketMessage, 0),
	}
}

// SetMessageHandler sets the message handler function
func (m *Manager) SetMessageHandler(handler func(*models.WebSocketMessage)) {
	m.onMessage = handler
}

// SetStateChangeHandler sets the state change handler function
func (m *Manager) SetStateChangeHandler(handler func(ConnectionState)) {
	m.onStateChange = handler
}

// Connect establishes a WebSocket connection to the specified game
func (m *Manager) Connect(gameID string) error {
	if m.state == Connected || m.state == Connecting {
		return fmt.Errorf("connection already active")
	}

	m.gameID = gameID
	m.setState(Connecting)

	return m.connect()
}

// connect performs the actual WebSocket connection
func (m *Manager) connect() error {
	// Construct WebSocket URL
	wsURL := m.buildWebSocketURL()

	// Create WebSocket connection
	ws := js.Global().Get("WebSocket").New(wsURL)
	m.websocket = ws

	// Set up event handlers
	m.setupEventHandlers()

	return nil
}

// buildWebSocketURL constructs the WebSocket URL
func (m *Manager) buildWebSocketURL() string {
	protocol := "ws"
	if js.Global().Get("location").Get("protocol").String() == "https:" {
		protocol = "wss"
	}

	host := js.Global().Get("location").Get("host").String()
	return fmt.Sprintf("%s://%s/ws/games/%s", protocol, host, m.gameID)
}

// setupEventHandlers sets up WebSocket event handlers
func (m *Manager) setupEventHandlers() {
	// Connection opened
	m.websocket.Set("onopen", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m.onOpen()
		return nil
	}))

	// Message received
	m.websocket.Set("onmessage", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			data := args[0].Get("data").String()
			m.onMessageReceived(data)
		}
		return nil
	}))

	// Connection closed
	m.websocket.Set("onclose", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) > 0 {
			code := args[0].Get("code").Int()
			reason := args[0].Get("reason").String()
			m.onClose(code, reason)
		}
		return nil
	}))

	// Connection error
	m.websocket.Set("onerror", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m.onError()
		return nil
	}))
}

// onOpen handles successful connection
func (m *Manager) onOpen() {
	m.setState(Connected)
	m.retryCount = 0
	m.startHeartbeat()
	m.flushMessageQueue()

	// Log successful connection
	js.Global().Get("console").Call("log", "WebSocket connected to game:", m.gameID)
}

// onMessageReceived handles incoming messages
func (m *Manager) onMessageReceived(data string) {
	// Update last heartbeat time
	m.lastHeartbeat = time.Now()

	// Parse the message
	var message models.WebSocketMessage
	if err := json.Unmarshal([]byte(data), &message); err != nil {
		js.Global().Get("console").Call("error", "Failed to parse WebSocket message:", err.Error())
		return
	}

	// Handle heartbeat/ping messages
	if message.Type == "ping" {
		m.sendPong()
		return
	}

	// Call message handler if set
	if m.onMessage != nil {
		m.onMessage(&message)
	}
}

// onClose handles connection closure
func (m *Manager) onClose(code int, reason string) {
	m.setState(Disconnected)
	m.stopHeartbeat()

	js.Global().Get("console").Call("log", fmt.Sprintf("WebSocket closed (code: %d, reason: %s)", code, reason))

	// Attempt reconnection if auto-reconnect is enabled and not manually closed
	if m.autoReconnect && code != 1000 {
		m.scheduleReconnection()
	}
}

// onError handles connection errors
func (m *Manager) onError() {
	js.Global().Get("console").Call("error", "WebSocket error occurred")

	if m.state == Connecting {
		m.setState(Disconnected)
		m.scheduleReconnection()
	}
}

// scheduleReconnection schedules a reconnection attempt with exponential backoff
func (m *Manager) scheduleReconnection() {
	if m.retryCount >= m.maxRetries {
		js.Global().Get("console").Call("error", "Max reconnection attempts reached")
		return
	}

	m.setState(Reconnecting)
	m.retryCount++

	// Calculate delay with exponential backoff and jitter
	delay := m.calculateRetryDelay()

	js.Global().Get("console").Call("log", fmt.Sprintf("Scheduling reconnection attempt %d in %s", m.retryCount, delay))

	m.reconnectTimer = time.AfterFunc(delay, func() {
		m.setState(Connecting)
		if err := m.connect(); err != nil {
			js.Global().Get("console").Call("error", "Reconnection failed:", err.Error())
		}
	})
}

// calculateRetryDelay calculates the delay for the next retry attempt
func (m *Manager) calculateRetryDelay() time.Duration {
	// Exponential backoff: baseDelay * 2^retryCount
	delay := time.Duration(float64(m.baseRetryDelay) * math.Pow(2, float64(m.retryCount-1)))

	// Add jitter (±25%)
	jitter := time.Duration(float64(delay) * 0.25 * (rand.Float64()*2 - 1))
	delay += jitter

	// Cap at max delay
	if delay > m.maxRetryDelay {
		delay = m.maxRetryDelay
	}

	return delay
}

// startHeartbeat starts the heartbeat mechanism
func (m *Manager) startHeartbeat() {
	m.lastHeartbeat = time.Now()

	m.heartbeatTimer = time.AfterFunc(m.heartbeatInterval, func() {
		m.checkHeartbeat()
	})
}

// stopHeartbeat stops the heartbeat mechanism
func (m *Manager) stopHeartbeat() {
	if m.heartbeatTimer != nil {
		m.heartbeatTimer.Stop()
		m.heartbeatTimer = nil
	}
}

// checkHeartbeat checks if the connection is still alive
func (m *Manager) checkHeartbeat() {
	if m.state != Connected {
		return
	}

	// Check if we've received a message recently
	if time.Since(m.lastHeartbeat) > m.heartbeatInterval*2 {
		js.Global().Get("console").Call("warn", "Heartbeat timeout, connection may be stale")
		m.Close()
		return
	}

	// Send ping
	m.sendPing()

	// Schedule next heartbeat check
	m.heartbeatTimer = time.AfterFunc(m.heartbeatInterval, func() {
		m.checkHeartbeat()
	})
}

// sendPing sends a ping message
func (m *Manager) sendPing() {
	message := &models.WebSocketMessage{
		Type:   "ping",
		GameID: m.gameID,
	}
	m.SendMessage(message)
}

// sendPong sends a pong response
func (m *Manager) sendPong() {
	message := &models.WebSocketMessage{
		Type:   "pong",
		GameID: m.gameID,
	}
	m.SendMessage(message)
}

// SendMessage sends a message over the WebSocket connection
func (m *Manager) SendMessage(message *models.WebSocketMessage) error {
	if m.state != Connected {
		// Queue the message for later delivery
		m.queueMessage(message)
		return fmt.Errorf("not connected, message queued")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Check if WebSocket is ready
	if m.websocket.Get("readyState").Int() != 1 {
		m.queueMessage(message)
		return fmt.Errorf("websocket not ready, message queued")
	}

	m.websocket.Call("send", string(data))
	return nil
}

// queueMessage adds a message to the queue for later delivery
func (m *Manager) queueMessage(message *models.WebSocketMessage) {
	// Prevent queue from growing too large
	if len(m.messageQueue) >= 100 {
		// Remove oldest message
		m.messageQueue = m.messageQueue[1:]
	}

	m.messageQueue = append(m.messageQueue, message)
}

// flushMessageQueue sends all queued messages
func (m *Manager) flushMessageQueue() {
	for _, message := range m.messageQueue {
		if err := m.SendMessage(message); err != nil {
			js.Global().Get("console").Call("error", "Failed to send queued message:", err.Error())
		}
	}
	m.messageQueue = m.messageQueue[:0] // Clear the queue
}

// Close closes the WebSocket connection
func (m *Manager) Close() {
	m.autoReconnect = false

	if m.reconnectTimer != nil {
		m.reconnectTimer.Stop()
		m.reconnectTimer = nil
	}

	m.stopHeartbeat()

	if m.websocket.Truthy() && m.websocket.Get("readyState").Int() < 2 {
		m.websocket.Call("close")
	}

	m.setState(Disconnected)
}

// GetState returns the current connection state
func (m *Manager) GetState() ConnectionState {
	return m.state
}

// IsConnected returns true if the WebSocket is connected
func (m *Manager) IsConnected() bool {
	return m.state == Connected
}

// SetAutoReconnect enables or disables automatic reconnection
func (m *Manager) SetAutoReconnect(enabled bool) {
	m.autoReconnect = enabled
}

// GetRetryCount returns the current retry count
func (m *Manager) GetRetryCount() int {
	return m.retryCount
}

// setState updates the connection state and notifies handlers
func (m *Manager) setState(state ConnectionState) {
	if m.state == state {
		return
	}

	m.state = state

	if m.onStateChange != nil {
		m.onStateChange(state)
	}

	// Update UI connection indicator
	m.updateConnectionIndicator()
}

// updateConnectionIndicator updates the connection status indicator in the UI
func (m *Manager) updateConnectionIndicator() {
	indicator := js.Global().Get("document").Call("getElementById", "connection-status")
	if !indicator.Truthy() {
		return
	}

	// Remove existing status classes
	classList := indicator.Get("classList")
	classList.Call("remove", "connected", "connecting", "disconnected", "reconnecting")

	var statusText string
	var statusClass string
	fmt.Println("Connection Status Update: ", m.state)
	switch m.state {
	case Connected:
		statusText = "Connected"
		statusClass = "connected"
	case Connecting:
		statusText = "Connecting..."
		statusClass = "connecting"
	case Reconnecting:
		statusText = fmt.Sprintf("Reconnecting... (%d/%d)", m.retryCount, m.maxRetries)
		statusClass = "connecting"
	case Disconnected:
		statusText = "Disconnected"
		statusClass = "disconnected"
	}

	classList.Call("add", statusClass)
	indicator.Set("textContent", statusText)
}

// StateString returns a string representation of the connection state
func (state ConnectionState) String() string {
	switch state {
	case Disconnected:
		return "disconnected"
	case Connecting:
		return "connecting"
	case Connected:
		return "connected"
	case Reconnecting:
		return "reconnecting"
	default:
		return "unknown"
	}
}
