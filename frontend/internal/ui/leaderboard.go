package ui

import (
	"fmt"
	"strings"

	"golf-gamez-frontend/internal/models"
)

// LeaderboardGenerator handles leaderboard HTML generation
type LeaderboardGenerator struct{}

// NewLeaderboardGenerator creates a new leaderboard generator
func NewLeaderboardGenerator() *LeaderboardGenerator {
	return &LeaderboardGenerator{}
}

// GenerateLeaderboard creates the complete leaderboard HTML
func (lg *LeaderboardGenerator) GenerateLeaderboard(leaderboard *models.Leaderboard) string {
	if leaderboard == nil || len(leaderboard.Overall) == 0 {
		return lg.generateEmptyLeaderboard()
	}

	html := strings.Builder{}
	html.WriteString(`<div class="leaderboard">`)
	html.WriteString(`<div class="leaderboard-header">`)
	html.WriteString(`<h3 class="leaderboard-title">Live Leaderboard</h3>`)
	html.WriteString(`</div>`)

	// Generate overall leaderboard entries
	for _, entry := range leaderboard.Overall {
		html.WriteString(lg.generateLeaderboardEntry(entry))
	}

	html.WriteString(`</div>`)
	return html.String()
}

// generateLeaderboardEntry creates a single leaderboard entry
func (lg *LeaderboardGenerator) generateLeaderboardEntry(entry models.LeaderboardEntry) string {
	html := strings.Builder{}

	html.WriteString(`<div class="leaderboard-entry">`)

	// Position indicator
	positionClass := lg.getPositionClass(entry.Position)
	html.WriteString(fmt.Sprintf(`<div class="leaderboard-position %s">%d</div>`, positionClass, entry.Position))

	// Player info
	html.WriteString(`<div class="leaderboard-player">`)
	html.WriteString(fmt.Sprintf(`<div class="leaderboard-player-name">%s</div>`, entry.Player.Name))
	html.WriteString(fmt.Sprintf(`<div class="leaderboard-player-details">%d holes • %d putts • %.1f hdcp</div>`,
		entry.HolesCompleted, entry.TotalPutts, entry.Player.Handicap))
	html.WriteString(`</div>`)

	// Score and trend
	html.WriteString(`<div class="leaderboard-score-container">`)
	html.WriteString(fmt.Sprintf(`<div class="leaderboard-score">%s</div>`, entry.Score))

	// Trend indicator
	if entry.Trend != "" {
		trendIcon := lg.getTrendIcon(entry.Trend)
		trendClass := lg.getTrendClass(entry.Trend)
		html.WriteString(fmt.Sprintf(`<div class="leaderboard-trend %s">%s</div>`, trendClass, trendIcon))
	}
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)
	return html.String()
}

// GenerateBestNineStandings creates Best Nine side bet standings
func (lg *LeaderboardGenerator) GenerateBestNineStandings(standings *models.BestNineStandings) string {
	if standings == nil || len(standings.Standings) == 0 {
		return lg.generateEmptyBestNine()
	}

	html := strings.Builder{}
	html.WriteString(`<div class="side-bet-standings">`)

	for _, result := range standings.Standings {
		html.WriteString(lg.generateBestNineEntry(result))
	}

	html.WriteString(`</div>`)
	return html.String()
}

// generateBestNineEntry creates a Best Nine standing entry
func (lg *LeaderboardGenerator) generateBestNineEntry(result models.BestNineResult) string {
	html := strings.Builder{}

	html.WriteString(`<div class="best-nine-entry">`)
	html.WriteString(fmt.Sprintf(`<div class="best-nine-position">%d</div>`, result.Position))

	html.WriteString(`<div class="best-nine-player">`)
	html.WriteString(fmt.Sprintf(`<div class="player-name">%s</div>`, result.Player.Name))
	html.WriteString(fmt.Sprintf(`<div class="best-nine-details">%d holes completed</div>`, result.HolesCompleted))
	html.WriteString(`</div>`)

	html.WriteString(`<div class="best-nine-scores">`)
	html.WriteString(fmt.Sprintf(`<div class="best-nine-score">%s</div>`, result.BestNineScore))
	html.WriteString(`<div class="best-nine-label">Best 9</div>`)
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)
	return html.String()
}

// GeneratePuttPuttPokerStatus creates Putt Putt Poker status display
func (lg *LeaderboardGenerator) GeneratePuttPuttPokerStatus(status *models.PuttPuttPokerStatus) string {
	if status == nil || len(status.Players) == 0 {
		return lg.generateEmptyPuttPoker()
	}

	html := strings.Builder{}
	html.WriteString(`<div class="putt-poker-standings">`)

	for _, result := range status.Players {
		html.WriteString(lg.generatePuttPokerEntry(result))
	}

	// Pot information
	if status.PotInfo.TotalPot > 0 {
		html.WriteString(lg.generatePotInfo(status.PotInfo))
	}

	html.WriteString(`</div>`)
	return html.String()
}

// generatePuttPokerEntry creates a Putt Putt Poker entry
func (lg *LeaderboardGenerator) generatePuttPokerEntry(result models.PuttPuttPokerResult) string {
	html := strings.Builder{}

	html.WriteString(`<div class="putt-poker-entry">`)

	// Player info
	html.WriteString(`<div class="putt-poker-player">`)
	html.WriteString(fmt.Sprintf(`<div class="player-name">%s</div>`, result.Player.Name))
	html.WriteString(`</div>`)

	// Card display
	html.WriteString(`<div class="poker-cards">`)
	for i := 0; i < result.TotalCards; i++ {
		cardClass := "poker-card"
		if i < result.StartingCards {
			cardClass += " starting"
		} else {
			cardClass += " earned"
		}
		html.WriteString(fmt.Sprintf(`<div class="%s">♠</div>`, cardClass))
	}
	html.WriteString(`</div>`)

	// Stats
	html.WriteString(`<div class="putt-poker-stats">`)
	html.WriteString(fmt.Sprintf(`<div class="stat">%d cards</div>`, result.TotalCards))
	if result.Penalties > 0 {
		html.WriteString(fmt.Sprintf(`<div class="stat penalty">%d penalties</div>`, result.Penalties))
	}
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)
	return html.String()
}

// generatePotInfo creates pot information display
func (lg *LeaderboardGenerator) generatePotInfo(potInfo interface{}) string {
	// This would need proper implementation based on pot info structure
	return `
		<div class="pot-info">
			<div class="pot-title">Current Pot</div>
			<div class="pot-amount">$0.00</div>
		</div>
	`
}

// generateEmptyLeaderboard creates placeholder leaderboard
func (lg *LeaderboardGenerator) generateEmptyLeaderboard() string {
	return `
		<div class="leaderboard">
			<div class="leaderboard-header">
				<h3 class="leaderboard-title">Leaderboard</h3>
			</div>
			<div class="p-8 text-center text-gray-500">
				<p>Leaderboard will update as scores are entered</p>
			</div>
		</div>
	`
}

// generateEmptyBestNine creates placeholder Best Nine display
func (lg *LeaderboardGenerator) generateEmptyBestNine() string {
	return `
		<div class="side-bet-standings">
			<div class="p-4 text-center text-gray-500">
				<p>Best Nine standings will appear after more scores are entered</p>
			</div>
		</div>
	`
}

// generateEmptyPuttPoker creates placeholder Putt Putt Poker display
func (lg *LeaderboardGenerator) generateEmptyPuttPoker() string {
	return `
		<div class="putt-poker-standings">
			<div class="p-4 text-center text-gray-500">
				<p>Card counts will update as putts are recorded</p>
			</div>
		</div>
	`
}

// Helper methods

func (lg *LeaderboardGenerator) getPositionClass(position int) string {
	switch position {
	case 1:
		return "first"
	case 2:
		return "second"
	case 3:
		return "third"
	default:
		return ""
	}
}

func (lg *LeaderboardGenerator) getTrendIcon(trend string) string {
	switch trend {
	case "up":
		return "↗"
	case "down":
		return "↘"
	case "same":
		return "→"
	default:
		return ""
	}
}

func (lg *LeaderboardGenerator) getTrendClass(trend string) string {
	switch trend {
	case "up":
		return "trend-up"
	case "down":
		return "trend-down"
	case "same":
		return "trend-same"
	default:
		return ""
	}
}