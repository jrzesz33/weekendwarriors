//go:build js && wasm
// +build js,wasm

package app

import (
	"fmt"
	"syscall/js"

	"golf-gamez-frontend/internal/api"
	"golf-gamez-frontend/internal/models"
	"golf-gamez-frontend/internal/storage"
	"golf-gamez-frontend/internal/ui"
	"golf-gamez-frontend/internal/websocket"
)

// App represents the main application structure
type App struct {
	apiClient    *api.Client
	wsManager    *websocket.Manager
	uiManager    *ui.Manager
	storage      *storage.LocalStorage
	currentGame  *models.Game
	currentRoute string
}

// New creates a new application instance
func New() *App {
	return &App{
		currentRoute: "/",
	}
}

// SetAPIClient sets the API client
func (a *App) SetAPIClient(client *api.Client) {
	a.apiClient = client
}

// SetWebSocketManager sets the WebSocket manager
func (a *App) SetWebSocketManager(manager *websocket.Manager) {
	a.wsManager = manager
}

// SetUIManager sets the UI manager
func (a *App) SetUIManager(manager *ui.Manager) {
	//manager.CreateClicked = a.createGameHandler
	a.uiManager = manager
}

// SetStorage sets the storage manager
func (a *App) SetStorage(storage *storage.LocalStorage) {
	a.storage = storage
}

// Initialize initializes the application
func (a *App) Initialize() error {
	// Set up DOM manipulation functions
	a.setupDOMFunctions()

	// Set up routing
	a.setupRouting()

	// Load initial route
	a.handleRoute()

	return nil
}

// setupDOMFunctions registers JavaScript functions for DOM manipulation
func (a *App) setupDOMFunctions() {
	fmt.Println("Setting up DOM Functions")
	js.Global().Set("golfGamez", js.ValueOf(map[string]interface{}{
		"createGame":    js.FuncOf(a.createGameHandler),
		"addPlayer":     js.FuncOf(a.addPlayerHandler),
		"startGame":     js.FuncOf(a.startGameHandler),
		"recordScore":   js.FuncOf(a.recordScoreHandler),
		"navigateTo":    js.FuncOf(a.navigateToHandler),
		"joinSpectator": js.FuncOf(a.joinSpectatorHandler),
		"updateScore":   js.FuncOf(a.updateScoreHandler),
		"completeGame":  js.FuncOf(a.completeGameHandler),
	}))
}

// setupRouting sets up client-side routing
func (a *App) setupRouting() {
	// Listen for popstate events (back/forward browser navigation)
	popstateHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		a.handleRoute()
		return nil
	})
	js.Global().Get("window").Call("addEventListener", "popstate", popstateHandler)
}

// handleRoute handles the current route
func (a *App) handleRoute() {
	path := js.Global().Get("window").Get("location").Get("pathname").String()
	a.currentRoute = path

	switch {
	case path == "/" || path == "/create":
		a.uiManager.ShowCreateGameView()
	case path == "/game":
		gameID := js.Global().Get("window").Get("location").Get("search").String()
		if gameID != "" {
			a.loadGame(gameID[1:]) // Remove ? from query string
		}
	case path == "/spectate":
		token := js.Global().Get("window").Get("location").Get("search").String()
		if token != "" {
			a.loadSpectatorView(token[1:]) // Remove ? from query string
		}
	default:
		a.uiManager.ShowNotFoundView()
	}
}

// Navigation handlers
func (a *App) navigateToHandler(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	path := args[0].String()
	js.Global().Get("window").Get("history").Call("pushState", nil, "", path)
	a.handleRoute()

	return nil
}

// Game creation handler
func (a *App) createGameHandler(this js.Value, args []js.Value) interface{} {
	fmt.Println("Create Game Clicked")

	go func() {
		request := models.CreateGameRequest{
			Course:          "diamond-run",
			HandicapEnabled: true,
			SideBets:        []string{"best-nine", "putt-putt-poker"},
		}

		game, err := a.apiClient.CreateGame(request)
		if err != nil {
			a.uiManager.ShowError("Failed to create game: " + err.Error())
			return
		}

		a.currentGame = game
		a.storage.SaveGame(game)

		// Navigate to game setup
		js.Global().Call("eval", `window.golfGamez.navigateTo('/game?'+"`+game.ID+`"')`)

	}()
	return nil

}

// Player management handlers
func (a *App) addPlayerHandler(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 || a.currentGame == nil {
		return nil
	}

	go func() {
		request := models.CreatePlayerRequest{
			Name:     args[0].String(),
			Handicap: args[1].Float(),
			Gender:   args[2].String(),
		}

		player, err := a.apiClient.AddPlayer(a.currentGame.ID, request)
		if err != nil {
			a.uiManager.ShowError("Failed to add player: " + err.Error())
			return
		}

		// Add player to current game
		a.currentGame.Players = append(a.currentGame.Players, *player)
		a.storage.SaveGame(a.currentGame)

		// Update UI
		a.uiManager.UpdatePlayerList(a.currentGame.Players)
	}()

	return nil
}

// Game control handlers
func (a *App) startGameHandler(this js.Value, args []js.Value) interface{} {
	if a.currentGame == nil {
		return nil
	}

	go func() {
		err := a.apiClient.StartGame(a.currentGame.ID)
		if err != nil {
			a.uiManager.ShowError("Failed to start game: " + err.Error())
			return
		}

		a.currentGame.Status = "in_progress"
		a.storage.SaveGame(a.currentGame)

		// Connect to WebSocket for real-time updates
		err = a.wsManager.Connect(a.currentGame.ID)
		if err != nil {
			a.uiManager.ShowWarning("Live updates may not work: " + err.Error())
		}

		// Show game interface
		a.uiManager.ShowGameView(a.currentGame)
	}()

	return nil
}

// Score recording handlers
func (a *App) recordScoreHandler(this js.Value, args []js.Value) interface{} {
	if len(args) < 4 || a.currentGame == nil {
		return nil
	}

	go func() {
		playerID := args[0].String()
		hole := int(args[1].Int())
		strokes := int(args[2].Int())
		putts := int(args[3].Int())

		request := models.ScoreRequest{
			Hole:    hole,
			Strokes: strokes,
			Putts:   putts,
		}

		score, err := a.apiClient.RecordScore(a.currentGame.ID, playerID, request)
		if err != nil {
			a.uiManager.ShowError("Failed to record score: " + err.Error())
			return
		}

		// Update UI with new score
		a.uiManager.UpdateScore(score)
	}()

	return nil
}

func (a *App) updateScoreHandler(this js.Value, args []js.Value) interface{} {
	if len(args) < 4 || a.currentGame == nil {
		return nil
	}

	go func() {
		playerID := args[0].String()
		hole := int(args[1].Int())
		strokes := int(args[2].Int())
		putts := int(args[3].Int())

		err := a.apiClient.UpdateScore(a.currentGame.ID, playerID, hole, strokes, putts)
		if err != nil {
			a.uiManager.ShowError("Failed to update score: " + err.Error())
			return
		}

		// Refresh scorecard
		a.refreshScorecard()
	}()

	return nil
}

func (a *App) completeGameHandler(this js.Value, args []js.Value) interface{} {
	if a.currentGame == nil {
		return nil
	}

	go func() {
		result, err := a.apiClient.CompleteGame(a.currentGame.ID)
		if err != nil {
			a.uiManager.ShowError("Failed to complete game: " + err.Error())
			return
		}

		a.currentGame.Status = "completed"
		a.currentGame.FinalResults = result.FinalResults
		a.storage.SaveGame(a.currentGame)

		// Show final results
		a.uiManager.ShowFinalResults(result.FinalResults)
	}()

	return nil
}

// Spectator mode handler
func (a *App) joinSpectatorHandler(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return nil
	}

	token := args[0].String()
	a.loadSpectatorView(token)

	return nil
}

// Utility functions
func (a *App) loadGame(gameID string) {
	go func() {
		game, err := a.apiClient.GetGame(gameID)
		if err != nil {
			a.uiManager.ShowError("Failed to load game: " + err.Error())
			return
		}

		a.currentGame = game
		a.storage.SaveGame(game)

		if game.Status == "in_progress" {
			// Connect to WebSocket
			err = a.wsManager.Connect(game.ID)
			if err != nil {
				a.uiManager.ShowWarning("Live updates may not work: " + err.Error())
			}
			a.uiManager.ShowGameView(game)
		} else {
			a.uiManager.ShowGameSetupView(game)
		}
	}()
}

func (a *App) loadSpectatorView(token string) {
	go func() {
		spectatorView, err := a.apiClient.GetSpectatorView(token)
		if err != nil {
			a.uiManager.ShowError("Failed to load spectator view: " + err.Error())
			return
		}

		// Connect to WebSocket for live updates
		err = a.wsManager.Connect(spectatorView.Game.ID)
		if err != nil {
			a.uiManager.ShowWarning("Live updates may not work: " + err.Error())
		}

		a.uiManager.ShowSpectatorView(spectatorView)
	}()
}

func (a *App) refreshScorecard() {
	if a.currentGame == nil {
		return
	}

	go func() {
		scorecard, err := a.apiClient.GetScorecard(a.currentGame.ID)
		if err != nil {
			a.uiManager.ShowError("Failed to refresh scorecard: " + err.Error())
			return
		}

		a.uiManager.UpdateScorecard(scorecard)
	}()
}
