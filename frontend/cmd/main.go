package main

import (
	"fmt"
	"strings"

	"golf-gamez-frontend/internal/models"
	"golf-gamez-frontend/internal/services"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// AppState represents the current application state
type AppState string

const (
	StateHome        AppState = "home"
	StateCreateGame  AppState = "create_game"
	StateJoinGame    AppState = "join_game"
	StateGameSetup   AppState = "game_setup"
	StateGameActive  AppState = "game_active"
	StateGameResults AppState = "game_results"
	StateSpectator   AppState = "spectator"
)

// GolfGamezApp is the main application component
type GolfGamezApp struct {
	app.Compo

	// Application state
	currentState AppState
	currentGame  *models.Game
	apiClient    *services.APIClient

	// UI state
	isLoading      bool
	errorMessage   string
	successMessage string
	menuOpen       bool

	// Scorecard state
	currentHole    int

	// Form state
	gameToken      string
	spectatorToken string
}

// OnMount initializes the component when mounted
func (g *GolfGamezApp) OnMount(ctx app.Context) {
	g.apiClient = services.NewAPIClient("")
	g.currentState = StateHome
	g.currentHole = 1 // Start with hole 1

	// Check URL for direct links
	g.handleInitialRoute(ctx)
}

// handleInitialRoute handles deep linking and initial route
func (g *GolfGamezApp) handleInitialRoute(ctx app.Context) {
	path := app.Window().Get("location").Get("pathname").String()
	search := app.Window().Get("location").Get("search").String()

	switch {
	case strings.HasPrefix(path, "/game/"):
		// Extract game token from path
		parts := strings.Split(path, "/")
		if len(parts) >= 3 {
			g.gameToken = parts[2]
			g.loadGame(ctx)
		}
	case strings.HasPrefix(path, "/spectate/"):
		// Extract spectator token from path
		parts := strings.Split(path, "/")
		if len(parts) >= 3 {
			g.spectatorToken = parts[2]
			g.loadSpectatorView(ctx)
		}
	case strings.Contains(search, "join="):
		// Handle join links with token in query
		params := strings.Split(strings.TrimPrefix(search, "?"), "&")
		for _, param := range params {
			if strings.HasPrefix(param, "join=") {
				g.gameToken = strings.TrimPrefix(param, "join=")
				g.loadGame(ctx)
				break
			}
		}
	}
}

func (g *GolfGamezApp) Render() app.UI {
	fmt.Println("Rendering ", g.currentState)
	return app.Div().
		Class("golf-app").
		Body(
			g.renderHeader(),
			g.renderMain(),
			g.renderFooter(),
			g.renderNotifications(),
		)
}

func (g *GolfGamezApp) renderHeader() app.UI {
	return app.Header().
		Class("header").
		Body(
			app.Div().
				Class("header-content").
				Body(
					app.H1().
						Class("title").
						Text("⛳ Golf Gamez").
						OnClick(g.onGoHome),
					app.Div().
						Class("header-actions").
						Body(
							app.If(g.currentState != StateHome, func() app.UI {
								return app.Button().
									Class("back-btn").
									Text("← Back").
									OnClick(g.onGoHome)
							}),
							g.renderMenuButton(),
						),
				),
		)
}

func (g *GolfGamezApp) renderMenuButton() app.UI {
	return app.Button().
		Class("menu-button").
		OnClick(g.onMenuToggle).
		Body(
			app.Div().
				Class("menu-icon").
				Body(
					app.Span(),
					app.Span(),
					app.Span(),
				),
		)
}

func (g *GolfGamezApp) renderMain() app.UI {
	switch g.currentState {
	case StateHome:
		return g.renderHome()
	case StateCreateGame:
		return g.renderCreateGame()
	case StateJoinGame:
		return g.renderJoinGame()
	case StateGameSetup:
		return g.renderGameSetup()
	case StateGameActive:
		return g.renderGameActive()
	case StateGameResults:
		return g.renderGameResults()
	case StateSpectator:
		return g.renderSpectatorView()
	default:
		return g.renderHome()
	}
}

func (g *GolfGamezApp) renderFooter() app.UI {
	return app.Footer().
		Class("footer").
		Body(
			app.P().Text("Built with Go + WebAssembly • PWA Technology"),
		)
}

func (g *GolfGamezApp) renderNotifications() app.UI {
	//fmt.Println("Rendering Notifications...", g.isLoading, g.errorMessage, g.successMessage)
	return app.Div().
		Class("notifications").
		Body(
			app.If(g.isLoading, func() app.UI {
				return app.Div().
					Class("loading-overlay").
					Body(
						app.Div().
							Class("loading-spinner"),
						app.P().Text("Loading..."),
					)
			}),
			app.If(g.errorMessage != "", func() app.UI {
				return app.Div().
					Class("error-message").
					Body(
						app.P().Text(g.errorMessage),
						app.Button().
							Text("×").
							OnClick(g.onClearError),
					)
			}),
			app.If(g.successMessage != "", func() app.UI {
				return app.Div().
					Class("success-message").
					Body(
						app.P().Text(g.successMessage),
						app.Button().
							Text("×").
							OnClick(g.onClearSuccess),
					)
			}),
		)
}

func main() {
	fmt.Println("Starting Golf Gamez App")

	// Route the main app component for all paths
	app.Route("/", func() app.Composer { return &GolfGamezApp{} })
	app.Route("/game/", func() app.Composer { return &GolfGamezApp{} })
	app.Route("/spectate/", func() app.Composer { return &GolfGamezApp{} })

	fmt.Println("Running Golf Gamez in Browser")
	// Run the app when in browser (WebAssembly)
	app.RunWhenOnBrowser()
}
