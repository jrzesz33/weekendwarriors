//go:build js && wasm
// +build js,wasm

package ui

import (
	"fmt"
	"strconv"
	"syscall/js"

	"golf-gamez-frontend/internal/models"
)

// Manager handles UI operations and DOM manipulation
type Manager struct {
	currentView string
	//CreateClicked func(js.Value, []js.Value) interface{}
}

// NewManager creates a new UI manager
func NewManager() *Manager {
	return &Manager{}
}

// ShowCreateGameView displays the game creation interface
func (ui *Manager) ShowCreateGameView() {
	ui.currentView = "create-game"
	ui.setPageContent(createGameHTML)
	ui.setupCreateGameHandlers()
}

// ShowGameSetupView displays the game setup interface for adding players
func (ui *Manager) ShowGameSetupView(game *models.Game) {
	ui.currentView = "game-setup"
	content := ui.generateGameSetupHTML(game)
	ui.setPageContent(content)
	ui.setupGameSetupHandlers(game)
}

// ShowGameView displays the main game interface
func (ui *Manager) ShowGameView(game *models.Game) {
	ui.currentView = "game-view"
	content := ui.generateGameViewHTML(game)
	ui.setPageContent(content)
	ui.setupGameViewHandlers(game)
}

// ShowSpectatorView displays the spectator interface
func (ui *Manager) ShowSpectatorView(view *models.SpectatorView) {
	ui.currentView = "spectator-view"
	content := ui.generateSpectatorViewHTML(view)
	ui.setPageContent(content)
}

// ShowNotFoundView displays a 404 page
func (ui *Manager) ShowNotFoundView() {
	ui.currentView = "not-found"
	ui.setPageContent(notFoundHTML)
}

// ShowFinalResults displays the final game results
func (ui *Manager) ShowFinalResults(results *models.FinalResults) {
	content := ui.generateFinalResultsHTML(results)
	ui.showModal("Game Complete!", content)
}

// ShowError displays an error message
func (ui *Manager) ShowError(message string) {
	ui.showToast("error", "Error", message)
}

// ShowWarning displays a warning message
func (ui *Manager) ShowWarning(message string) {
	ui.showToast("warning", "Warning", message)
}

// ShowSuccess displays a success message
func (ui *Manager) ShowSuccess(message string) {
	ui.showToast("success", "Success", message)
}

// UpdatePlayerList updates the player list display
func (ui *Manager) UpdatePlayerList(players []models.Player) {
	playerListEl := js.Global().Get("document").Call("getElementById", "player-list")
	if !playerListEl.Truthy() {
		return
	}

	html := ""
	for _, player := range players {
		html += ui.generatePlayerCardHTML(player)
	}

	playerListEl.Set("innerHTML", html)
}

// UpdateScore updates a score display
func (ui *Manager) UpdateScore(score *models.Score) {
	// Update scorecard if visible
	ui.updateScorecardScore(score)

	// Update leaderboard
	ui.refreshLeaderboard()

	// Show success message
	ui.ShowSuccess(fmt.Sprintf("Score recorded for hole %d", score.Hole))
}

// UpdateScorecard updates the scorecard display
func (ui *Manager) UpdateScorecard(scorecard *models.GameScorecard) {
	scorecardEl := js.Global().Get("document").Call("getElementById", "scorecard-container")
	if !scorecardEl.Truthy() {
		return
	}

	html := ui.generateScorecardHTML(scorecard)
	scorecardEl.Set("innerHTML", html)
}

// UpdateLeaderboard updates the leaderboard display
func (ui *Manager) UpdateLeaderboard(leaderboard *models.Leaderboard) {
	leaderboardEl := js.Global().Get("document").Call("getElementById", "leaderboard-container")
	if !leaderboardEl.Truthy() {
		return
	}

	html := ui.generateLeaderboardHTML(leaderboard)
	leaderboardEl.Set("innerHTML", html)
}

// setPageContent sets the main page content
func (ui *Manager) setPageContent(html string) {
	appEl := js.Global().Get("document").Call("getElementById", "app")
	if appEl.Truthy() {
		appEl.Set("innerHTML", html)
	}
}

// setupCreateGameHandlers sets up event handlers for the create game view
func (ui *Manager) setupCreateGameHandlers() {
	// Course selection handler
	courseSelect := js.Global().Get("document").Call("getElementById", "course-select")
	if courseSelect.Truthy() {
		courseSelect.Set("value", "diamond-run")
	}

	// Side bets checkboxes
	bestNineCheck := js.Global().Get("document").Call("getElementById", "best-nine-enabled")
	puttPokerCheck := js.Global().Get("document").Call("getElementById", "putt-poker-enabled")
	handicapCheck := js.Global().Get("document").Call("getElementById", "handicap-enabled")

	if bestNineCheck.Truthy() {
		bestNineCheck.Set("checked", true)
	}
	if puttPokerCheck.Truthy() {
		puttPokerCheck.Set("checked", true)
	}
	if handicapCheck.Truthy() {
		handicapCheck.Set("checked", false)
	}

	/*/register the button handler
	createBtn := js.Global().Get("document").Call("getElementById", "createGmBtn")
	if createBtn.Truthy() {
		fmt.Println("registring the create button")
		if ui.CreateClicked != nil {
			createBtn.Call("addEventListener", "click", js.FuncOf(ui.CreateClicked))
		}
	}//*/
	fmt.Println("Finished Registering Handlers")
}

// setupGameSetupHandlers sets up handlers for the game setup view
func (ui *Manager) setupGameSetupHandlers(game *models.Game) {
	// Add player form handler
	addPlayerForm := js.Global().Get("document").Call("getElementById", "add-player-form")
	if addPlayerForm.Truthy() {
		addPlayerHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			args[0].Call("preventDefault")
			ui.handleAddPlayer()
			return nil
		})
		addPlayerForm.Call("addEventListener", "submit", addPlayerHandler)
	}

	// Handicap help handlers
	ui.setupHandicapHelpers()

	// Share link setup
	ui.setupShareLink(game.ShareLink, game.SpectatorLink)
}

// setupGameViewHandlers sets up handlers for the main game view
func (ui *Manager) setupGameViewHandlers(game *models.Game) {
	// Score entry handlers
	ui.setupScoreEntryHandlers()

	// Navigation handlers
	ui.setupGameNavigation()
}

// handleAddPlayer handles adding a new player
func (ui *Manager) handleAddPlayer() {
	nameInput := js.Global().Get("document").Call("getElementById", "player-name")
	handicapInput := js.Global().Get("document").Call("getElementById", "player-handicap")
	genderSelect := js.Global().Get("document").Call("getElementById", "player-gender")

	if !nameInput.Truthy() || !handicapInput.Truthy() || !genderSelect.Truthy() {
		return
	}

	name := nameInput.Get("value").String()
	handicapStr := handicapInput.Get("value").String()
	gender := genderSelect.Get("value").String()

	if name == "" {
		ui.ShowError("Player name is required")
		return
	}

	handicap, err := strconv.ParseFloat(handicapStr, 64)
	if err != nil {
		ui.ShowError("Invalid handicap value")
		return
	}

	if handicap < 0 || handicap > 54 {
		ui.ShowError("Handicap must be between 0 and 54")
		return
	}

	// Call the global addPlayer function
	js.Global().Call("eval", fmt.Sprintf("window.golfGamez.addPlayer('%s', %f, '%s')", name, handicap, gender))

	// Clear the form
	nameInput.Set("value", "")
	handicapInput.Set("value", "")
	genderSelect.Set("value", "")
}

// setupHandicapHelpers sets up handicap guidance
func (ui *Manager) setupHandicapHelpers() {
	genderSelect := js.Global().Get("document").Call("getElementById", "player-gender")
	handicapGuide := js.Global().Get("document").Call("getElementById", "handicap-guide")

	if !genderSelect.Truthy() || !handicapGuide.Truthy() {
		return
	}

	genderChangeHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		gender := genderSelect.Get("value").String()
		var guideText string

		switch gender {
		case "male":
			guideText = "Typical bogey golfer handicap for men: 18-22"
		case "female":
			guideText = "Typical bogey golfer handicap for women: 22-26"
		default:
			guideText = "Choose your handicap based on your typical scoring"
		}

		handicapGuide.Set("textContent", guideText)
		return nil
	})

	genderSelect.Call("addEventListener", "change", genderChangeHandler)
}

// setupShareLink sets up share link functionality
func (ui *Manager) setupShareLink(shareLink, spectatorLink string) {
	shareInput := js.Global().Get("document").Call("getElementById", "share-link-input")
	spectatorInput := js.Global().Get("document").Call("getElementById", "spectator-link-input")
	copyShareBtn := js.Global().Get("document").Call("getElementById", "copy-share-btn")
	copySpectatorBtn := js.Global().Get("document").Call("getElementById", "copy-spectator-btn")

	if shareInput.Truthy() {
		shareInput.Set("value", shareLink)
	}
	if spectatorInput.Truthy() {
		spectatorInput.Set("value", spectatorLink)
	}

	if copyShareBtn.Truthy() {
		copyHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ui.copyToClipboard(shareLink)
			ui.ShowSuccess("Share link copied to clipboard!")
			return nil
		})
		copyShareBtn.Call("addEventListener", "click", copyHandler)
	}

	if copySpectatorBtn.Truthy() {
		copyHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ui.copyToClipboard(spectatorLink)
			ui.ShowSuccess("Spectator link copied to clipboard!")
			return nil
		})
		copySpectatorBtn.Call("addEventListener", "click", copyHandler)
	}
}

// setupScoreEntryHandlers sets up score entry functionality
func (ui *Manager) setupScoreEntryHandlers() {
	// Number stepper handlers
	steppers := js.Global().Get("document").Call("querySelectorAll", ".number-stepper")
	for i := 0; i < steppers.Get("length").Int(); i++ {
		stepper := steppers.Call("item", i)
		ui.setupNumberStepper(stepper)
	}

	// Score form submission
	scoreForm := js.Global().Get("document").Call("getElementById", "score-entry-form")
	if scoreForm.Truthy() {
		submitHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			args[0].Call("preventDefault")
			ui.handleScoreSubmission()
			return nil
		})
		scoreForm.Call("addEventListener", "submit", submitHandler)
	}
}

// setupNumberStepper sets up a number stepper component
func (ui *Manager) setupNumberStepper(stepper js.Value) {
	input := stepper.Call("querySelector", "input")
	minusBtn := stepper.Call("querySelector", ".minus-btn")
	plusBtn := stepper.Call("querySelector", ".plus-btn")

	if !input.Truthy() || !minusBtn.Truthy() || !plusBtn.Truthy() {
		return
	}

	// Minus button handler
	minusHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentValue, _ := strconv.Atoi(input.Get("value").String())
		minValue, _ := strconv.Atoi(input.Get("min").String())
		if currentValue > minValue {
			input.Set("value", currentValue-1)
		}
		return nil
	})
	minusBtn.Call("addEventListener", "click", minusHandler)

	// Plus button handler
	plusHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		currentValue, _ := strconv.Atoi(input.Get("value").String())
		maxValue, _ := strconv.Atoi(input.Get("max").String())
		if currentValue < maxValue {
			input.Set("value", currentValue+1)
		}
		return nil
	})
	plusBtn.Call("addEventListener", "click", plusHandler)
}

// handleScoreSubmission handles score form submission
func (ui *Manager) handleScoreSubmission() {
	playerSelect := js.Global().Get("document").Call("getElementById", "score-player-select")
	holeInput := js.Global().Get("document").Call("getElementById", "score-hole")
	strokesInput := js.Global().Get("document").Call("getElementById", "score-strokes")
	puttsInput := js.Global().Get("document").Call("getElementById", "score-putts")

	if !playerSelect.Truthy() || !holeInput.Truthy() || !strokesInput.Truthy() || !puttsInput.Truthy() {
		return
	}

	playerID := playerSelect.Get("value").String()
	hole, _ := strconv.Atoi(holeInput.Get("value").String())
	strokes, _ := strconv.Atoi(strokesInput.Get("value").String())
	putts, _ := strconv.Atoi(puttsInput.Get("value").String())

	if playerID == "" {
		ui.ShowError("Please select a player")
		return
	}

	if hole < 1 || hole > 18 {
		ui.ShowError("Invalid hole number")
		return
	}

	if strokes < 1 || strokes > 15 {
		ui.ShowError("Invalid stroke count")
		return
	}

	if putts < 0 || putts > 10 {
		ui.ShowError("Invalid putt count")
		return
	}

	// Call the global recordScore function
	js.Global().Call("eval", fmt.Sprintf("window.golfGamez.recordScore('%s', %d, %d, %d)", playerID, hole, strokes, putts))
}

// copyToClipboard copies text to clipboard
func (ui *Manager) copyToClipboard(text string) {
	navigator := js.Global().Get("navigator")
	if navigator.Get("clipboard").Truthy() {
		navigator.Get("clipboard").Call("writeText", text)
	} else {
		// Fallback for older browsers
		textArea := js.Global().Get("document").Call("createElement", "textarea")
		textArea.Set("value", text)
		js.Global().Get("document").Get("body").Call("appendChild", textArea)
		textArea.Call("select")
		js.Global().Get("document").Call("execCommand", "copy")
		js.Global().Get("document").Get("body").Call("removeChild", textArea)
	}
}

// showToast displays a toast notification
func (ui *Manager) showToast(type_, title, message string) {
	toastHTML := fmt.Sprintf(`
		<div class="toast toast-%s" id="toast-notification">
			<div class="toast-header">
				<strong>%s</strong>
				<button class="toast-close" onclick="this.parentElement.parentElement.remove()">×</button>
			</div>
			<div class="toast-body">%s</div>
		</div>
	`, type_, title, message)

	// Create or get toast container
	toastContainer := js.Global().Get("document").Call("getElementById", "toast-container")
	if !toastContainer.Truthy() {
		toastContainer = js.Global().Get("document").Call("createElement", "div")
		toastContainer.Set("id", "toast-container")
		toastContainer.Set("className", "toast-container")
		js.Global().Get("document").Get("body").Call("appendChild", toastContainer)
	}

	// Add toast
	toastElement := js.Global().Get("document").Call("createElement", "div")
	toastElement.Set("innerHTML", toastHTML)
	toastContainer.Call("appendChild", toastElement.Get("firstElementChild"))

	// Auto-remove after 5 seconds
	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		toast := js.Global().Get("document").Call("getElementById", "toast-notification")
		if toast.Truthy() {
			toast.Call("remove")
		}
		return nil
	}), 5000)
}

// showModal displays a modal dialog
func (ui *Manager) showModal(title, content string) {
	modalHTML := fmt.Sprintf(`
		<div class="modal-backdrop" id="modal-backdrop" onclick="this.remove()">
			<div class="modal" onclick="event.stopPropagation()">
				<div class="modal-header">
					<h3>%s</h3>
					<button class="modal-close" onclick="document.getElementById('modal-backdrop').remove()">×</button>
				</div>
				<div class="modal-body">%s</div>
			</div>
		</div>
	`, title, content)

	modalElement := js.Global().Get("document").Call("createElement", "div")
	modalElement.Set("innerHTML", modalHTML)
	js.Global().Get("document").Get("body").Call("appendChild", modalElement.Get("firstElementChild"))
}

// Helper methods for generating HTML content
func (ui *Manager) generatePlayerCardHTML(player models.Player) string {
	return fmt.Sprintf(`
		<div class="player-card" data-player-id="%s">
			<div class="player-name">%s</div>
			<div class="player-handicap">Handicap: %.1f</div>
			<div class="player-stats">
				<div class="player-stat">
					<span class="stat-value">%s</span>
					<span class="stat-label">Score</span>
				</div>
				<div class="player-stat">
					<span class="stat-value">%d</span>
					<span class="stat-label">Holes</span>
				</div>
				<div class="player-stat">
					<span class="stat-value">%d</span>
					<span class="stat-label">Putts</span>
				</div>
			</div>
		</div>
	`, player.ID, player.Name, player.Handicap, player.Stats.CurrentScore,
		player.Stats.HolesCompleted, player.Stats.TotalPutts)
}

// Additional helper methods would be implemented here for generating
// scorecard HTML, leaderboard HTML, etc.

// Placeholder implementations for remaining UI generation methods
func (ui *Manager) generateGameSetupHTML(game *models.Game) string {
	return gameSetupHTML // This would be defined as a constant
}

func (ui *Manager) generateGameViewHTML(game *models.Game) string {
	return gameViewHTML // This would be defined as a constant
}

func (ui *Manager) generateSpectatorViewHTML(view *models.SpectatorView) string {
	return spectatorViewHTML // This would be defined as a constant
}

func (ui *Manager) generateFinalResultsHTML(results *models.FinalResults) string {
	return "Final results content" // Would be implemented with actual HTML
}

func (ui *Manager) generateScorecardHTML(scorecard *models.GameScorecard) string {
	return "Scorecard content" // Would be implemented with actual HTML
}

func (ui *Manager) generateLeaderboardHTML(leaderboard *models.Leaderboard) string {
	return "Leaderboard content" // Would be implemented with actual HTML
}

func (ui *Manager) updateScorecardScore(score *models.Score) {
	// Implementation for updating specific scorecard cell
}

func (ui *Manager) refreshLeaderboard() {
	// Implementation for refreshing leaderboard data
}

func (ui *Manager) setupGameNavigation() {
	// Implementation for game navigation setup
}
