package main

import (
	"fmt"
	"syscall/js"

	"golf-gamez-frontend/internal/api"
	"golf-gamez-frontend/internal/app"
	"golf-gamez-frontend/internal/storage"
	"golf-gamez-frontend/internal/ui"
	"golf-gamez-frontend/internal/websocket"
)

func main() {
	fmt.Println("App Initializing")
	// Initialize the application
	app := app.New()

	// Set up API client
	apiClient := api.NewClient("http://localhost:8080/v1")
	app.SetAPIClient(apiClient)

	// Set up WebSocket manager
	wsManager := websocket.NewManager()
	app.SetWebSocketManager(wsManager)

	// Set up local storage
	localStorage := storage.NewLocalStorage()
	app.SetStorage(localStorage)

	// Set up UI manager
	uiManager := ui.NewManager()
	app.SetUIManager(uiManager)

	// Initialize the application
	if err := app.Initialize(); err != nil {
		js.Global().Get("console").Call("error", "Failed to initialize app:", err.Error())
		return
	}

	// Keep the application running
	select {}
}
