package main

import (
	"fmt"
	"strconv"
	"strings"

	"golf-gamez-frontend/internal/models"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Navigation handlers
func (g *GolfGamezApp) onGoHome(ctx app.Context, e app.Event) {
	g.currentState = StateHome
	g.currentGame = nil
	g.clearMessages()

	// Update URL
	app.Window().Get("history").Call("pushState", nil, "", "/")
}

func (g *GolfGamezApp) onCreateGameStart(ctx app.Context, e app.Event) {
	g.currentState = StateCreateGame
	g.clearMessages()
}

func (g *GolfGamezApp) onJoinGameStart(ctx app.Context, e app.Event) {
	g.currentState = StateJoinGame
	g.clearMessages()
}

// Game creation handlers
func (g *GolfGamezApp) onCreateGameSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// Get form data
	form := e.Get("target")

	// Extract side bets
	var sideBets []models.SideBetType
	checkboxes := form.Call("querySelectorAll", "input[type=checkbox]:checked")
	length := checkboxes.Get("length").Int()

	for i := 0; i < length; i++ {
		value := checkboxes.Index(i).Get("value").String()
		switch value {
		case "best-nine":
			sideBets = append(sideBets, models.SideBetBestNine)
		case "putt-putt-poker":
			sideBets = append(sideBets, models.SideBetPuttPuttPoker)
		}
	}

	// Check handicap enabled
	handicapCheckbox := form.Call("querySelector", "#handicap-enabled")
	handicapEnabled := handicapCheckbox.Get("checked").Bool()

	// Create game request
	req := models.CreateGameRequest{
		Course:          "diamond-run",
		SideBets:        sideBets,
		HandicapEnabled: handicapEnabled,
	}

	g.createGame(ctx, req)
}

func (g *GolfGamezApp) createGame(ctx app.Context, req models.CreateGameRequest) {
	g.setLoading(true)

	ctx.Async(func() {
		game, err := g.apiClient.CreateGame(ctx, req)
		g.currentGame = game
		g.currentState = StateGameSetup
		token := g.apiClient.ExtractTokenFromShareLink(game.ShareLink)

		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				g.setError(fmt.Sprintf("Failed to create game: %v", err))
				fmt.Println("Error Creating Game: ", err)
				return
			}
			g.setLoading(false)
			fmt.Println("Game Created Successfully")
			g.setSuccess("Game created successfully! Share the link with other players.")

			// Update URL
			app.Window().Get("history").Call("pushState", nil, "", fmt.Sprintf("/game/%s", token))
			fmt.Println("Pushed the State...")
		})
	})
}

// Game joining handlers
func (g *GolfGamezApp) onGameTokenChange(ctx app.Context, e app.Event) {
	g.gameToken = e.Get("target").Get("value").String()
	// Extract token if it's a full URL
	if strings.Contains(g.gameToken, "/game/") {
		parts := strings.Split(g.gameToken, "/game/")
		if len(parts) > 1 {
			g.gameToken = strings.TrimSpace(parts[1])
		}
	}
}

func (g *GolfGamezApp) onJoinGameSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	token := strings.TrimSpace(g.gameToken)
	if token == "" {
		g.setError("Please enter a game code")
		return
	}

	g.loadGame(ctx, token)
}

func (g *GolfGamezApp) loadGame(ctx app.Context, token ...string) {
	gameToken := g.gameToken
	if len(token) > 0 {
		gameToken = token[0]
	}

	if gameToken == "" {
		g.setError("No game token provided")
		return
	}

	g.setLoading(true)

	ctx.Async(func() {
		game, err := g.apiClient.GetGame(ctx, gameToken)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				g.setError(fmt.Sprintf("Failed to load game: %v", err))
				return
			}

			g.currentGame = game
			g.gameToken = gameToken
			g.apiClient.SetAuthToken(gameToken)

			// Determine the correct state based on game status
			switch game.Status {
			case models.GameStatusSetup:
				g.currentState = StateGameSetup
			case models.GameStatusInProgress:
				g.currentState = StateGameActive
			case models.GameStatusCompleted:
				g.currentState = StateGameResults
			default:
				g.currentState = StateGameSetup
			}

			g.setLoading(false)

			// Update URL
			app.Window().Get("history").Call("pushState", nil, "", fmt.Sprintf("/game/%s", gameToken))
		})
	})
}

func (g *GolfGamezApp) loadSpectatorView(ctx app.Context, token ...string) {
	_ = g.spectatorToken // Mark as used
	if len(token) > 0 {
		_ = token[0] // Mark as used
	}
	ctx.Dispatch(func(ctx app.Context) {
		g.setLoading(true)
		// TODO: Implement spectator view loading
		g.setError("Spectator view not yet implemented")
	})

}

// Player management handlers
func (g *GolfGamezApp) onAddPlayerSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	form := e.Get("target")
	name := form.Call("querySelector", "#player-name").Get("value").String()
	handicapStr := form.Call("querySelector", "#player-handicap").Get("value").String()
	gender := form.Call("querySelector", "#player-gender").Get("value").String()

	if name == "" {
		g.setError("Player name is required")
		return
	}

	handicap, err := strconv.ParseFloat(handicapStr, 64)
	if err != nil {
		g.setError("Invalid handicap value")
		return
	}

	if gender == "" {
		g.setError("Please select a gender")
		return
	}

	req := models.CreatePlayerRequest{
		Name:     name,
		Handicap: handicap,
		Gender:   models.Gender(gender),
	}

	g.addPlayer(ctx, req)
}

func (g *GolfGamezApp) addPlayer(ctx app.Context, req models.CreatePlayerRequest) {
	if g.currentGame == nil {
		g.setError("No active game")
		return
	}

	g.setLoading(true)

	ctx.Async(func() {
		token := g.apiClient.ExtractTokenFromShareLink(g.currentGame.ShareLink)
		player, err := g.apiClient.AddPlayer(ctx, token, req)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				g.setError(fmt.Sprintf("Failed to add player: %v", err))
				return
			}

			g.currentGame.Players = append(g.currentGame.Players, *player)
			g.setSuccess(fmt.Sprintf("Player %s added successfully!", player.Name))
			g.setLoading(false)
			g.clearPlayerForm()
		})
	})
}

func (g *GolfGamezApp) clearPlayerForm() {
	// Clear the form fields
	app.Window().Get("document").Call("querySelector", "#player-name").Set("value", "")
	app.Window().Get("document").Call("querySelector", "#player-handicap").Set("value", "")
	app.Window().Get("document").Call("querySelector", "#player-gender").Set("value", "")
}

func (g *GolfGamezApp) clearScoreForm() {
	// Clear the score entry form fields
	app.Window().Get("document").Call("querySelector", "#strokes").Set("value", "")
	app.Window().Get("document").Call("querySelector", "#putts").Set("value", "")
}

// Game control handlers
func (g *GolfGamezApp) onStartGame(ctx app.Context, e app.Event) {
	if g.currentGame == nil {
		g.setError("No active game")
		return
	}

	if len(g.currentGame.Players) == 0 {
		g.setError("Add at least one player before starting the game")
		return
	}

	g.setLoading(true)

	ctx.Async(func() {
		token := g.apiClient.ExtractTokenFromShareLink(g.currentGame.ShareLink)
		game, err := g.apiClient.StartGame(ctx, token)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				g.setError(fmt.Sprintf("Failed to start game: %v", err))
				return
			}
			g.currentGame.Status = game.Status
			g.currentState = StateGameActive
			g.setSuccess("Game started! Let's play golf!")
			g.setLoading(false)
		})
	})
}

// Score management handlers
func (g *GolfGamezApp) onScoreSubmit(ctx app.Context, e app.Event) {
	e.PreventDefault()

	form := e.Get("target")
	playerID := form.Call("querySelector", "#player-select").Get("value").String()
	strokesStr := form.Call("querySelector", "#strokes").Get("value").String()
	puttsStr := form.Call("querySelector", "#putts").Get("value").String()

	if playerID == "" {
		g.setError("Please select a player")
		return
	}

	strokes, err := strconv.Atoi(strokesStr)
	if err != nil || strokes < 1 || strokes > 15 {
		g.setError("Please enter valid strokes (1-15)")
		return
	}

	putts, err := strconv.Atoi(puttsStr)
	if err != nil || putts < 0 || putts > 10 {
		g.setError("Please enter valid putts (0-10)")
		return
	}

	req := models.CreateScoreRequest{
		Hole:    g.currentHole,
		Strokes: strokes,
		Putts:   putts,
	}

	g.recordScore(ctx, playerID, req)
	g.clearScoreForm()
}

func (g *GolfGamezApp) onScoreChange(ctx app.Context, playerID string, hole int, scoreStr string) {
	score, err := strconv.Atoi(scoreStr)
	if err != nil || score < 1 || score > 15 {
		return // Invalid score, ignore
	}

	// For now, just log the score change
	// In a full implementation, this would save to the backend
	fmt.Printf("Score change: Player %s, Hole %d, Score %d\n", playerID, hole, score)

	// TODO: Implement actual score saving
	req := models.CreateScoreRequest{
		Hole:    hole,
		Strokes: score,
		Putts:   2, // Default putts for now
	}

	g.recordScore(ctx, playerID, req)
}

func (g *GolfGamezApp) recordScore(ctx app.Context, playerID string, req models.CreateScoreRequest) {
	if g.currentGame == nil {
		return
	}

	ctx.Async(func() {
		token := g.apiClient.ExtractTokenFromShareLink(g.currentGame.ShareLink)
		_, err := g.apiClient.RecordScore(ctx, token, playerID, req)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				g.setError(fmt.Sprintf("Failed to record score: %v", err))
				return
			}

			g.setSuccess("Score recorded!")
		})
		// TODO: Update local game state
	})
}

// Utility handlers
func (g *GolfGamezApp) onClearError(ctx app.Context, e app.Event) {
	g.errorMessage = ""
}

func (g *GolfGamezApp) onClearSuccess(ctx app.Context, e app.Event) {
	g.successMessage = ""
}

func (g *GolfGamezApp) onMenuToggle(ctx app.Context, e app.Event) {
	g.menuOpen = !g.menuOpen
}

// Hole navigation handlers
func (g *GolfGamezApp) onPreviousHole(ctx app.Context, e app.Event) {
	if g.currentHole > 1 {
		g.currentHole--
	}
}

func (g *GolfGamezApp) onNextHole(ctx app.Context, e app.Event) {
	maxHoles := 18
	if g.currentGame != nil && g.currentGame.CourseInfo != nil {
		maxHoles = len(g.currentGame.CourseInfo.Holes)
	}
	if g.currentHole < maxHoles {
		g.currentHole++
	}
}

func (g *GolfGamezApp) onCopyShareCode(ctx app.Context, e app.Event) {
	if g.currentGame == nil {
		return
	}

	token := g.apiClient.ExtractTokenFromShareLink(g.currentGame.ShareLink)
	g.copyToClipboard(token)
	g.setSuccess("Share code copied to clipboard!")
}

func (g *GolfGamezApp) onCopyShareURL(ctx app.Context, e app.Event) {
	url := g.getShareURL()
	g.copyToClipboard(url)
	g.setSuccess("Share link copied to clipboard!")
}

func (g *GolfGamezApp) onCopySpectatorURL(ctx app.Context, e app.Event) {
	url := g.getSpectatorURL()
	g.copyToClipboard(url)
	g.setSuccess("Spectator link copied to clipboard!")
}

func (g *GolfGamezApp) onSelectText(ctx app.Context, e app.Event) {
	e.Get("target").Call("select")
}

func (g *GolfGamezApp) copyToClipboard(text string) {
	navigator := app.Window().Get("navigator")
	if !navigator.Get("clipboard").IsUndefined() {
		navigator.Get("clipboard").Call("writeText", text)
	} else {
		// Fallback for older browsers
		textarea := app.Window().Get("document").Call("createElement", "textarea")
		textarea.Set("value", text)
		app.Window().Get("document").Get("body").Call("appendChild", textarea)
		textarea.Call("select")
		app.Window().Get("document").Call("execCommand", "copy")
		app.Window().Get("document").Get("body").Call("removeChild", textarea)
	}
}

// State management helpers
func (g *GolfGamezApp) setLoading(loading bool) {
	g.isLoading = loading
}

func (g *GolfGamezApp) setError(message string) {
	g.errorMessage = message
	g.isLoading = false
}

func (g *GolfGamezApp) setSuccess(message string) {
	g.successMessage = message
	g.errorMessage = ""
}

func (g *GolfGamezApp) clearMessages() {
	g.errorMessage = ""
	g.successMessage = ""
}
