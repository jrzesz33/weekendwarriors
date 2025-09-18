//go:build js && wasm
// +build js,wasm

package ui

import (
	"fmt"
	"strconv"
	"strings"

	"golf-gamez-frontend/internal/models"
)

// ScorecardGenerator handles scorecard HTML generation
type ScorecardGenerator struct{}

// NewScorecardGenerator creates a new scorecard generator
func NewScorecardGenerator() *ScorecardGenerator {
	return &ScorecardGenerator{}
}

// GenerateScorecard creates the complete scorecard HTML
func (sg *ScorecardGenerator) GenerateScorecard(scorecard *models.GameScorecard) string {
	if scorecard == nil || len(scorecard.Players) == 0 {
		return sg.generateEmptyScorecard()
	}

	html := strings.Builder{}
	html.WriteString(`<div class="scorecard">`)
	html.WriteString(`<div class="scorecard-header">`)
	html.WriteString(`<h3 class="text-lg font-semibold">Live Scorecard</h3>`)
	html.WriteString(`</div>`)

	// Generate mobile-friendly scorecard
	html.WriteString(sg.generateMobileScorecard(scorecard))

	html.WriteString(`</div>`)
	return html.String()
}

// generateMobileScorecard creates a mobile-optimized scorecard layout
func (sg *ScorecardGenerator) generateMobileScorecard(scorecard *models.GameScorecard) string {
	html := strings.Builder{}

	// Front Nine Header
	html.WriteString(`<div class="scorecard-section">`)
	html.WriteString(`<div class="scorecard-section-header">Front Nine</div>`)
	html.WriteString(sg.generateNineHoles(scorecard, 1, 9))
	html.WriteString(`</div>`)

	// Back Nine Header
	html.WriteString(`<div class="scorecard-section">`)
	html.WriteString(`<div class="scorecard-section-header">Back Nine</div>`)
	html.WriteString(sg.generateNineHoles(scorecard, 10, 18))
	html.WriteString(`</div>`)

	// Totals Section
	html.WriteString(`<div class="scorecard-section">`)
	html.WriteString(`<div class="scorecard-section-header">Totals</div>`)
	html.WriteString(sg.generateTotalsSection(scorecard))
	html.WriteString(`</div>`)

	return html.String()
}

// generateNineHoles creates HTML for 9 holes
func (sg *ScorecardGenerator) generateNineHoles(scorecard *models.GameScorecard, startHole, endHole int) string {
	html := strings.Builder{}

	// Hole numbers row
	html.WriteString(`<div class="scorecard-row scorecard-holes">`)
	html.WriteString(`<div class="scorecard-cell scorecard-label">Hole</div>`)
	for hole := startHole; hole <= endHole; hole++ {
		html.WriteString(fmt.Sprintf(`<div class="scorecard-cell scorecard-hole">%d</div>`, hole))
	}
	html.WriteString(`<div class="scorecard-cell scorecard-total">Out/In</div>`)
	html.WriteString(`</div>`)

	// Par row
	html.WriteString(`<div class="scorecard-row scorecard-par">`)
	html.WriteString(`<div class="scorecard-cell scorecard-label">Par</div>`)
	subtotal := 0
	for hole := startHole; hole <= endHole; hole++ {
		par := sg.getParForHole(hole)
		subtotal += par
		html.WriteString(fmt.Sprintf(`<div class="scorecard-cell">%d</div>`, par))
	}
	html.WriteString(fmt.Sprintf(`<div class="scorecard-cell scorecard-total">%d</div>`, subtotal))
	html.WriteString(`</div>`)

	// Player rows
	for _, player := range scorecard.Players {
		html.WriteString(sg.generatePlayerRow(player, startHole, endHole))
	}

	return html.String()
}

// generatePlayerRow creates a player's score row
func (sg *ScorecardGenerator) generatePlayerRow(player interface{}, startHole, endHole int) string {
	html := strings.Builder{}

	// Type assertion to get player data
	// This would need to be properly implemented based on the actual player structure
	playerName := "Player" // Placeholder

	html.WriteString(`<div class="scorecard-row scorecard-player">`)
	html.WriteString(fmt.Sprintf(`<div class="scorecard-cell scorecard-label">%s</div>`, playerName))

	subtotal := 0
	for hole := startHole; hole <= endHole; hole++ {
		score := sg.getPlayerScoreForHole(player, hole)
		if score > 0 {
			subtotal += score
			par := sg.getParForHole(hole)
			scoreToPar := score - par
			cellClass := sg.getScoreCellClass(scoreToPar)
			html.WriteString(fmt.Sprintf(`<div class="scorecard-cell %s">%d</div>`, cellClass, score))
		} else {
			html.WriteString(`<div class="scorecard-cell scorecard-empty">-</div>`)
		}
	}

	if subtotal > 0 {
		html.WriteString(fmt.Sprintf(`<div class="scorecard-cell scorecard-total">%d</div>`, subtotal))
	} else {
		html.WriteString(`<div class="scorecard-cell scorecard-total">-</div>`)
	}

	html.WriteString(`</div>`)
	return html.String()
}

// generateTotalsSection creates the totals section
func (sg *ScorecardGenerator) generateTotalsSection(scorecard *models.GameScorecard) string {
	html := strings.Builder{}

	html.WriteString(`<div class="scorecard-totals">`)
	for _, player := range scorecard.Players {
		html.WriteString(sg.generatePlayerTotalRow(player))
	}
	html.WriteString(`</div>`)

	return html.String()
}

// generatePlayerTotalRow creates a player's total row
func (sg *ScorecardGenerator) generatePlayerTotalRow(player interface{}) string {
	html := strings.Builder{}

	playerName := "Player" // Placeholder
	totalScore := 0        // Calculate actual total
	totalToPar := 0        // Calculate vs par

	html.WriteString(`<div class="scorecard-total-row">`)
	html.WriteString(fmt.Sprintf(`<div class="player-name">%s</div>`, playerName))
	html.WriteString(fmt.Sprintf(`<div class="total-score">%d</div>`, totalScore))

	scoreIndicator := sg.formatScoreToPar(totalToPar)
	indicatorClass := sg.getScoreIndicatorClass(totalToPar)
	html.WriteString(fmt.Sprintf(`<div class="score-indicator %s">%s</div>`, indicatorClass, scoreIndicator))
	html.WriteString(`</div>`)

	return html.String()
}

// generateEmptyScorecard creates placeholder scorecard
func (sg *ScorecardGenerator) generateEmptyScorecard() string {
	return `
		<div class="scorecard">
			<div class="scorecard-header">
				<h3 class="text-lg font-semibold">Scorecard</h3>
			</div>
			<div class="p-8 text-center text-gray-500">
				<p>Scorecard will appear here as players enter scores</p>
			</div>
		</div>
	`
}

// Helper methods

func (sg *ScorecardGenerator) getParForHole(hole int) int {
	// Diamond Run pars
	pars := []int{4, 3, 5, 4, 3, 4, 4, 3, 5, 4, 3, 5, 3, 4, 4, 4, 5, 4}
	if hole >= 1 && hole <= 18 {
		return pars[hole-1]
	}
	return 4 // Default
}

func (sg *ScorecardGenerator) getPlayerScoreForHole(player interface{}, hole int) int {
	// This would need proper implementation based on player structure
	return 0 // Placeholder
}

func (sg *ScorecardGenerator) getScoreCellClass(scoreToPar int) string {
	switch {
	case scoreToPar <= -2:
		return "scorecard-eagle"
	case scoreToPar == -1:
		return "scorecard-birdie"
	case scoreToPar == 0:
		return "scorecard-par"
	case scoreToPar == 1:
		return "scorecard-bogey"
	default:
		return "scorecard-double-bogey"
	}
}

func (sg *ScorecardGenerator) formatScoreToPar(scoreToPar int) string {
	switch {
	case scoreToPar == 0:
		return "E"
	case scoreToPar > 0:
		return "+" + strconv.Itoa(scoreToPar)
	default:
		return strconv.Itoa(scoreToPar)
	}
}

func (sg *ScorecardGenerator) getScoreIndicatorClass(scoreToPar int) string {
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
