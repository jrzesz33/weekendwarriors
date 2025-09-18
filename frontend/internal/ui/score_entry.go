//go:build js && wasm
// +build js,wasm

package ui

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"

	"golf-gamez-frontend/internal/models"
)

// ScoreEntryManager handles score entry interface and validation
type ScoreEntryManager struct {
	currentHole   int
	currentPlayer string
	courseInfo    *models.CourseInfo
}

// NewScoreEntryManager creates a new score entry manager
func NewScoreEntryManager() *ScoreEntryManager {
	return &ScoreEntryManager{
		currentHole: 1,
	}
}

// SetCourseInfo sets the course information for validation
func (sem *ScoreEntryManager) SetCourseInfo(courseInfo *models.CourseInfo) {
	sem.courseInfo = courseInfo
}

// GenerateScoreEntryInterface creates the score entry HTML
func (sem *ScoreEntryManager) GenerateScoreEntryInterface(players []models.Player, currentHole int) string {
	sem.currentHole = currentHole

	html := strings.Builder{}
	html.WriteString(`<div class="score-entry-container">`)

	// Quick hole navigation
	html.WriteString(sem.generateHoleNavigation())

	// Current hole info
	html.WriteString(sem.generateCurrentHoleInfo())

	// Score entry form
	html.WriteString(sem.generateScoreEntryForm(players))

	// Recent scores display
	html.WriteString(sem.generateRecentScores())

	html.WriteString(`</div>`)
	return html.String()
}

// generateHoleNavigation creates hole navigation interface
func (sem *ScoreEntryManager) generateHoleNavigation() string {
	html := strings.Builder{}

	html.WriteString(`<div class="hole-navigation mb-4">`)
	html.WriteString(`<div class="hole-nav-header">Quick Hole Selection</div>`)
	html.WriteString(`<div class="hole-nav-grid">`)

	for hole := 1; hole <= 18; hole++ {
		activeClass := ""
		if hole == sem.currentHole {
			activeClass = " active"
		}

		html.WriteString(fmt.Sprintf(`
			<button class="hole-nav-btn%s" onclick="selectHole(%d)" data-hole="%d">
				<span class="hole-number">%d</span>
				<span class="hole-par">Par %d</span>
			</button>
		`, activeClass, hole, hole, hole, sem.getParForHole(hole)))
	}

	html.WriteString(`</div>`)
	html.WriteString(`</div>`)
	return html.String()
}

// generateCurrentHoleInfo creates current hole information display
func (sem *ScoreEntryManager) generateCurrentHoleInfo() string {
	par := sem.getParForHole(sem.currentHole)
	yardage := sem.getYardageForHole(sem.currentHole)

	return fmt.Sprintf(`
		<div class="current-hole-info mb-6">
			<div class="hole-display">%d</div>
			<div class="hole-details">
				<div class="hole-title">Hole %d</div>
				<div class="hole-specs">
					<span class="par-display">Par %d</span>
					<span class="yardage-display">%d yards</span>
				</div>
			</div>
		</div>
	`, sem.currentHole, sem.currentHole, par, yardage)
}

// generateScoreEntryForm creates the main score entry form
func (sem *ScoreEntryManager) generateScoreEntryForm(players []models.Player) string {
	html := strings.Builder{}

	html.WriteString(`<div class="card">`)
	html.WriteString(`<div class="card-header">`)
	html.WriteString(`<h3 class="card-title">Record Score</h3>`)
	html.WriteString(`<p class="card-subtitle">Enter strokes and putts for hole ` + strconv.Itoa(sem.currentHole) + `</p>`)
	html.WriteString(`</div>`)

	html.WriteString(`<form id="score-entry-form">`)

	// Player selection
	html.WriteString(`<div class="form-group">`)
	html.WriteString(`<label class="label" for="score-player-select">Player</label>`)
	html.WriteString(`<select id="score-player-select" class="select" required>`)
	html.WriteString(`<option value="">Select player</option>`)
	for _, player := range players {
		html.WriteString(fmt.Sprintf(`<option value="%s">%s</option>`, player.ID, player.Name))
	}
	html.WriteString(`</select>`)
	html.WriteString(`</div>`)

	// Score inputs with enhanced steppers
	html.WriteString(`<div class="score-inputs-grid">`)

	// Strokes input
	html.WriteString(`<div class="score-input-section">`)
	html.WriteString(`<label class="score-input-label">Strokes</label>`)
	html.WriteString(sem.generateEnhancedNumberStepper("strokes", 1, 15, 4))
	html.WriteString(sem.generateScoreHints("strokes"))
	html.WriteString(`</div>`)

	// Putts input
	html.WriteString(`<div class="score-input-section">`)
	html.WriteString(`<label class="score-input-label">Putts</label>`)
	html.WriteString(sem.generateEnhancedNumberStepper("putts", 0, 10, 2))
	html.WriteString(sem.generateScoreHints("putts"))
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)

	// Score to par preview
	html.WriteString(`<div id="score-preview" class="score-preview">`)
	html.WriteString(`<div class="score-preview-label">Score to Par</div>`)
	html.WriteString(`<div id="score-to-par-display" class="score-to-par-display">-</div>`)
	html.WriteString(`</div>`)

	// Submit button
	html.WriteString(`<button type="submit" class="btn btn-primary btn-lg btn-full">`)
	html.WriteString(`Record Score`)
	html.WriteString(`</button>`)

	html.WriteString(`</form>`)
	html.WriteString(`</div>`)

	return html.String()
}

// generateEnhancedNumberStepper creates an enhanced touch-friendly number stepper
func (sem *ScoreEntryManager) generateEnhancedNumberStepper(field string, min, max, defaultValue int) string {
	return fmt.Sprintf(`
		<div class="enhanced-number-stepper" data-field="%s">
			<button type="button" class="stepper-btn stepper-minus" data-action="minus">
				<span class="stepper-icon">−</span>
			</button>
			<input type="number"
				   id="score-%s"
				   class="stepper-input"
				   min="%d"
				   max="%d"
				   value="%d"
				   required
				   inputmode="numeric">
			<button type="button" class="stepper-btn stepper-plus" data-action="plus">
				<span class="stepper-icon">+</span>
			</button>
		</div>
	`, field, field, min, max, defaultValue)
}

// generateScoreHints creates score input hints
func (sem *ScoreEntryManager) generateScoreHints(field string) string {
	switch field {
	case "strokes":
		par := sem.getParForHole(sem.currentHole)
		return fmt.Sprintf(`
			<div class="score-hints">
				<div class="hint-item">Eagle: %d</div>
				<div class="hint-item">Birdie: %d</div>
				<div class="hint-item">Par: %d</div>
				<div class="hint-item">Bogey: %d</div>
			</div>
		`, par-2, par-1, par, par+1)
	case "putts":
		return `
			<div class="score-hints">
				<div class="hint-item">1 putt: +1 card</div>
				<div class="hint-item">Hole-in-one: +2 cards</div>
				<div class="hint-item">3+ putts: $1 penalty</div>
			</div>
		`
	default:
		return ""
	}
}

// generateRecentScores creates recent scores display
func (sem *ScoreEntryManager) generateRecentScores() string {
	return `
		<div class="recent-scores mt-6">
			<div class="recent-scores-header">
				<h4>Recent Scores</h4>
				<button class="btn btn-sm btn-secondary" onclick="toggleRecentScores()">
					View All
				</button>
			</div>
			<div id="recent-scores-list" class="recent-scores-list">
				<!-- Recent scores will be populated dynamically -->
			</div>
		</div>
	`
}

// SetupScoreEntryHandlers sets up enhanced event handlers for score entry
func (sem *ScoreEntryManager) SetupScoreEntryHandlers() {
	// Enhanced number stepper handlers with haptic feedback
	sem.setupEnhancedSteppers()

	// Real-time score preview
	sem.setupScorePreview()

	// Form validation and submission
	sem.setupFormValidation()

	// Auto-advance to next hole
	sem.setupAutoAdvance()
}

// setupEnhancedSteppers sets up the enhanced number steppers
func (sem *ScoreEntryManager) setupEnhancedSteppers() {
	steppers := js.Global().Get("document").Call("querySelectorAll", ".enhanced-number-stepper")

	for i := 0; i < steppers.Get("length").Int(); i++ {
		stepper := steppers.Call("item", i)
		input := stepper.Call("querySelector", ".stepper-input")
		minusBtn := stepper.Call("querySelector", ".stepper-minus")
		plusBtn := stepper.Call("querySelector", ".stepper-plus")

		if !input.Truthy() || !minusBtn.Truthy() || !plusBtn.Truthy() {
			continue
		}

		// Minus button with validation
		minusHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			currentValue, _ := strconv.Atoi(input.Get("value").String())
			minValue, _ := strconv.Atoi(input.Get("min").String())

			if currentValue > minValue {
				newValue := currentValue - 1
				input.Set("value", newValue)
				sem.triggerHapticFeedback()
				sem.updateScorePreview()
				sem.validateInput(input)
			}
			return nil
		})
		minusBtn.Call("addEventListener", "click", minusHandler)

		// Plus button with validation
		plusHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			currentValue, _ := strconv.Atoi(input.Get("value").String())
			maxValue, _ := strconv.Atoi(input.Get("max").String())

			if currentValue < maxValue {
				newValue := currentValue + 1
				input.Set("value", newValue)
				sem.triggerHapticFeedback()
				sem.updateScorePreview()
				sem.validateInput(input)
			}
			return nil
		})
		plusBtn.Call("addEventListener", "click", plusHandler)

		// Input change handler
		changeHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			sem.validateInput(input)
			sem.updateScorePreview()
			return nil
		})
		input.Call("addEventListener", "input", changeHandler)
	}
}

// setupScorePreview sets up real-time score to par preview
func (sem *ScoreEntryManager) setupScorePreview() {
	sem.updateScorePreview()
}

// updateScorePreview updates the score to par display
func (sem *ScoreEntryManager) updateScorePreview() {
	strokesInput := js.Global().Get("document").Call("getElementById", "score-strokes")
	display := js.Global().Get("document").Call("getElementById", "score-to-par-display")

	if !strokesInput.Truthy() || !display.Truthy() {
		return
	}

	strokesStr := strokesInput.Get("value").String()
	if strokesStr == "" {
		display.Set("textContent", "-")
		return
	}

	strokes, err := strconv.Atoi(strokesStr)
	if err != nil {
		display.Set("textContent", "-")
		return
	}

	par := sem.getParForHole(sem.currentHole)
	scoreToPar := strokes - par

	// Format and style the score
	scoreText := sem.formatScoreToPar(scoreToPar)
	scoreClass := sem.getScoreToParClass(scoreToPar)

	// Remove existing classes
	classList := display.Get("classList")
	classList.Call("remove", "eagle", "birdie", "par", "bogey", "double-bogey")
	classList.Call("add", scoreClass)

	display.Set("textContent", scoreText)
}

// setupFormValidation sets up comprehensive form validation
func (sem *ScoreEntryManager) setupFormValidation() {
	form := js.Global().Get("document").Call("getElementById", "score-entry-form")
	if !form.Truthy() {
		return
	}

	submitHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		event := args[0]
		event.Call("preventDefault")

		if sem.validateForm() {
			sem.submitScore()
		}
		return nil
	})
	form.Call("addEventListener", "submit", submitHandler)
}

// validateForm performs comprehensive form validation
func (sem *ScoreEntryManager) validateForm() bool {
	playerSelect := js.Global().Get("document").Call("getElementById", "score-player-select")
	strokesInput := js.Global().Get("document").Call("getElementById", "score-strokes")
	puttsInput := js.Global().Get("document").Call("getElementById", "score-putts")

	if !playerSelect.Truthy() || !strokesInput.Truthy() || !puttsInput.Truthy() {
		return false
	}

	// Validate player selection
	if playerSelect.Get("value").String() == "" {
		sem.showValidationError("Please select a player")
		playerSelect.Call("focus")
		return false
	}

	// Validate strokes
	strokesStr := strokesInput.Get("value").String()
	strokes, err := strconv.Atoi(strokesStr)
	if err != nil || strokes < 1 || strokes > 15 {
		sem.showValidationError("Strokes must be between 1 and 15")
		strokesInput.Call("focus")
		return false
	}

	// Validate putts
	puttsStr := puttsInput.Get("value").String()
	putts, err := strconv.Atoi(puttsStr)
	if err != nil || putts < 0 || putts > 10 {
		sem.showValidationError("Putts must be between 0 and 10")
		puttsInput.Call("focus")
		return false
	}

	// Logical validation - putts shouldn't exceed strokes
	if putts > strokes {
		sem.showValidationError("Putts cannot exceed total strokes")
		puttsInput.Call("focus")
		return false
	}

	return true
}

// submitScore submits the validated score
func (sem *ScoreEntryManager) submitScore() {
	playerSelect := js.Global().Get("document").Call("getElementById", "score-player-select")
	strokesInput := js.Global().Get("document").Call("getElementById", "score-strokes")
	puttsInput := js.Global().Get("document").Call("getElementById", "score-putts")

	playerID := playerSelect.Get("value").String()
	strokes, _ := strconv.Atoi(strokesInput.Get("value").String())
	putts, _ := strconv.Atoi(puttsInput.Get("value").String())

	// Call the global recordScore function
	js.Global().Call("eval", fmt.Sprintf(
		"window.golfGamez.recordScore('%s', %d, %d, %d)",
		playerID, sem.currentHole, strokes, putts))
}

// setupAutoAdvance sets up auto-advance to next hole
func (sem *ScoreEntryManager) setupAutoAdvance() {
	// Auto-advance could be implemented here
	// For now, we'll keep it manual for better UX control
}

// Validation and feedback methods

func (sem *ScoreEntryManager) validateInput(input js.Value) {
	value, _ := strconv.Atoi(input.Get("value").String())
	min, _ := strconv.Atoi(input.Get("min").String())
	max, _ := strconv.Atoi(input.Get("max").String())

	classList := input.Get("classList")
	classList.Call("remove", "input-valid", "input-invalid")

	if value >= min && value <= max {
		classList.Call("add", "input-valid")
	} else {
		classList.Call("add", "input-invalid")
	}
}

func (sem *ScoreEntryManager) showValidationError(message string) {
	// Create or update validation error display
	errorDiv := js.Global().Get("document").Call("getElementById", "validation-error")
	if !errorDiv.Truthy() {
		errorDiv = js.Global().Get("document").Call("createElement", "div")
		errorDiv.Set("id", "validation-error")
		errorDiv.Set("className", "alert alert-error mt-4")

		form := js.Global().Get("document").Call("getElementById", "score-entry-form")
		if form.Truthy() {
			form.Call("appendChild", errorDiv)
		}
	}

	errorDiv.Set("textContent", message)
	errorDiv.Set("style", "display: block")

	// Auto-hide after 3 seconds
	js.Global().Call("setTimeout", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		errorDiv.Set("style", "display: none")
		return nil
	}), 3000)
}

func (sem *ScoreEntryManager) triggerHapticFeedback() {
	// Trigger haptic feedback on supported devices
	navigator := js.Global().Get("navigator")
	if navigator.Get("vibrate").Truthy() {
		navigator.Call("vibrate", 50) // 50ms vibration
	}
}

// Helper methods
func (sem *ScoreEntryManager) getParForHole(hole int) int {
	pars := []int{4, 3, 5, 4, 3, 4, 4, 3, 5, 4, 3, 5, 3, 4, 4, 4, 5, 4}
	if hole >= 1 && hole <= 18 {
		return pars[hole-1]
	}
	return 4
}

func (sem *ScoreEntryManager) getYardageForHole(hole int) int {
	// Simplified yardage - in a real app this would come from course data
	return 350
}

func (sem *ScoreEntryManager) formatScoreToPar(scoreToPar int) string {
	switch {
	case scoreToPar == 0:
		return "E"
	case scoreToPar > 0:
		return "+" + strconv.Itoa(scoreToPar)
	default:
		return strconv.Itoa(scoreToPar)
	}
}

func (sem *ScoreEntryManager) getScoreToParClass(scoreToPar int) string {
	switch {
	case scoreToPar <= -2:
		return "eagle"
	case scoreToPar == -1:
		return "birdie"
	case scoreToPar == 0:
		return "par"
	case scoreToPar == 1:
		return "bogey"
	default:
		return "double-bogey"
	}
}
