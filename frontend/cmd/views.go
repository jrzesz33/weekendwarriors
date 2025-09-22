package main

import (
	"fmt"
	"strconv"

	"golf-gamez-frontend/internal/models"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Home view rendering
func (g *GolfGamezApp) renderHome() app.UI {
	return app.Main().
		Class("main-content home-view").
		Body(
			app.Div().
				Class("welcome-section").
				Body(
					app.H2().Text("Welcome to Golf Gamez"),
					app.P().Text("The ultimate Progressive Web App for weekend warriors who love their side bets!"),
				),

			app.Div().
				Class("features-grid").
				Body(
					g.renderFeatureCard("🎯", "Anonymous Games", "Create games instantly without registration"),
					g.renderFeatureCard("📊", "Live Scoring", "Real-time score tracking and leaderboards"),
					g.renderFeatureCard("🏆", "Best Nine", "Track your best 9 holes for side betting"),
					g.renderFeatureCard("🃏", "Putt Putt Poker", "Card-based putting games with poker hands"),
					g.renderFeatureCard("📱", "Mobile First", "Touch-optimized for outdoor golf course use"),
					g.renderFeatureCard("🔄", "Offline Ready", "PWA technology works without internet"),
				),

			app.Div().
				Class("action-section").
				Body(
					app.Button().
						Class("create-game-btn primary-btn").
						Text("🏌️ Create New Game").
						OnClick(g.onCreateGameStart),
					app.Button().
						Class("join-game-btn secondary-btn").
						Text("🔗 Join Game").
						OnClick(g.onJoinGameStart),
				),
		)
}

// Create game view
func (g *GolfGamezApp) renderCreateGame() app.UI {
	return app.Main().
		Class("main-content create-game-view").
		Body(
			app.H2().Text("🏌️ Create New Game"),
			app.P().Text("Set up a new golf game with your preferred settings"),

			app.Form().
				Class("game-form").
				OnSubmit(g.onCreateGameSubmit).
				Body(
					app.Div().
						Class("form-section").
						Body(
							app.H3().Text("⛳ Course"),
							app.Select().
								Class("course-select").
								Required(true).
								Body(
									app.Option().
										Value("diamond-run").
										Selected(true).
										Text("Diamond Run Golf Course (Par 71)"),
								),
						),

					app.Div().
						Class("form-section").
						Body(
							app.H3().Text("🎯 Side Bets"),
							app.Label().
								Class("checkbox-label").
								Body(
									app.Input().
										Type("checkbox").
										Value("best-nine").
										Checked(true),
									app.Span().Text("Best Nine - Track your best 9 holes"),
								),
							app.Label().
								Class("checkbox-label").
								Body(
									app.Input().
										Type("checkbox").
										Value("putt-putt-poker").
										Checked(true),
									app.Span().Text("Putt Putt Poker - Card-based putting game"),
								),
						),

					app.Div().
						Class("form-section").
						Body(
							app.H3().Text("♿ Handicap System"),
							app.Label().
								Class("checkbox-label").
								Body(
									app.Input().
										Type("checkbox").
										ID("handicap-enabled").
										Checked(true),
									app.Span().Text("Enable handicap scoring"),
								),
							app.P().
								Class("help-text").
								Text("Recommended: 20 for male bogey golfers, 24 for female"),
						),

					app.Div().
						Class("form-actions").
						Body(
							app.Button().
								Type("submit").
								Class("primary-btn").
								Text("🎮 Create Game"),
							app.Button().
								Type("button").
								Class("secondary-btn").
								Text("Cancel").
								OnClick(g.onGoHome),
						),
				),
		)
}

// Join game view
func (g *GolfGamezApp) renderJoinGame() app.UI {
	return app.Main().
		Class("main-content join-game-view").
		Body(
			app.H2().Text("🔗 Join Game"),
			app.P().Text("Enter the game code or share link to join an existing game"),

			app.Form().
				Class("join-form").
				OnSubmit(g.onJoinGameSubmit).
				Body(
					app.Div().
						Class("form-section").
						Body(
							app.Label().
								For("game-token").
								Text("Game Code or Share Link"),
							app.Input().
								Type("text").
								ID("game-token").
								Class("game-token-input").
								Placeholder("gt_xxxxxxxx or full share link").
								Required(true).
								Value(g.gameToken).
								OnInput(g.onGameTokenChange),
							app.P().
								Class("help-text").
								Text("Enter the game code (starts with gt_) or paste the full share link"),
						),

					app.Div().
						Class("form-actions").
						Body(
							app.Button().
								Type("submit").
								Class("primary-btn").
								Text("🎮 Join Game"),
							app.Button().
								Type("button").
								Class("secondary-btn").
								Text("Cancel").
								OnClick(g.onGoHome),
						),
				),
		)
}

// Game setup view (for adding players)
func (g *GolfGamezApp) renderGameSetup() app.UI {
	if g.currentGame == nil {
		return app.Div().Text("Loading game...")
	}
	fmt.Println("Course Info: ", g.currentGame.CourseInfo)
	return app.Main().
		Class("main-content game-setup-view").
		Body(
			app.H2().Text("🎯 Game Setup"),
			app.P().Text(fmt.Sprintf("Course: %s", g.currentGame.CourseInfo.Name)),

			g.renderGameInfo(),
			g.renderPlayersList(),
			g.renderAddPlayerForm(),
			g.renderGameActions(),
		)
}

// Game info display
func (g *GolfGamezApp) renderGameInfo() app.UI {
	if g.currentGame == nil {
		return app.Div()
	}

	return app.Div().
		Class("game-info").
		Body(
			app.H3().Text("📋 Game Information"),
			app.Div().
				Class("info-grid").
				Body(
					app.Div().
						Class("info-item").
						Body(
							app.Strong().Text("Course: "),
							app.Span().Text(g.currentGame.CourseInfo.Name),
						),
					app.Div().
						Class("info-item").
						Body(
							app.Strong().Text("Par: "),
							app.Span().Text(fmt.Sprintf("%d", g.currentGame.CourseInfo.TotalPar)),
						),
					app.Div().
						Class("info-item").
						Body(
							app.Strong().Text("Side Bets: "),
							app.Span().Text(fmt.Sprintf("%v", g.currentGame.SideBets)),
						),
					app.Div().
						Class("info-item").
						Body(
							app.Strong().Text("Share Code: "),
							app.Code().
								Class("share-code").
								Text(g.apiClient.ExtractTokenFromShareLink(g.currentGame.ShareLink)).
								OnClick(g.onCopyShareCode),
						),
				),
		)
}

// Players list display
func (g *GolfGamezApp) renderPlayersList() app.UI {
	if g.currentGame == nil {
		return app.Div()
	}

	return app.Div().
		Class("players-section").
		Body(
			app.H3().Text(fmt.Sprintf("👥 Players (%d/4)", len(g.currentGame.Players))),
			app.Div().
				Class("players-list").
				Body(
					app.Range(g.currentGame.Players).Slice(func(i int) app.UI {
						player := g.currentGame.Players[i]
						return app.Div().
							Class("player-card").
							Body(
								app.Div().
									Class("player-info").
									Body(
										app.Strong().Text(player.Name),
										app.Span().
											Class("player-handicap").
											Text(fmt.Sprintf("Handicap: %.1f", player.Handicap)),
									),
								app.Div().
									Class("player-stats").
									Body(
										app.Span().
											Class("player-gender").
											Text(string(player.Gender)),
									),
							)
					}),
				),
			app.If(len(g.currentGame.Players) == 0, func() app.UI {
				return app.P().
					Class("no-players").
					Text("No players added yet. Add at least one player to start the game.")
			}),
		)
}

// Add player form
func (g *GolfGamezApp) renderAddPlayerForm() app.UI {
	if g.currentGame == nil || len(g.currentGame.Players) >= 4 {
		return app.Div()
	}

	return app.Div().
		Class("add-player-section").
		Body(
			app.H3().Text("➕ Add Player"),
			app.Form().
				Class("player-form").
				OnSubmit(g.onAddPlayerSubmit).
				Body(
					app.Div().
						Class("form-row").
						Body(
							app.Input().
								Type("text").
								ID("player-name").
								Placeholder("Player Name").
								Required(true).
								Class("player-name-input"),
							app.Input().
								Type("number").
								ID("player-handicap").
								Placeholder("Handicap").
								Step(0.1).
								Min("0").
								Max("54").
								Required(true).
								Class("player-handicap-input"),
						),
					app.Div().
						Class("form-row").
						Body(
							app.Select().
								ID("player-gender").
								Required(true).
								Class("player-gender-select").
								Body(
									app.Option().Value("").Text("Select Gender"),
									app.Option().Value("male").Text("Male"),
									app.Option().Value("female").Text("Female"),
									app.Option().Value("other").Text("Other"),
								),
							app.Button().
								Type("submit").
								Class("add-player-btn").
								Text("➕ Add Player"),
						),
				),
		)
}

// Game actions (start game, share links)
func (g *GolfGamezApp) renderGameActions() app.UI {
	if g.currentGame == nil {
		return app.Div()
	}

	canStart := len(g.currentGame.Players) > 0 && g.currentGame.Status == models.GameStatusSetup

	return app.Div().
		Class("game-actions").
		Body(
			app.If(canStart, func() app.UI {
				return app.Button().
					Class("start-game-btn primary-btn").
					Text("🚀 Start Game").
					OnClick(g.onStartGame)
			}),
			app.Div().
				Class("share-section").
				Body(
					app.H3().Text("🔗 Share Game"),
					app.Div().
						Class("share-links").
						Body(
							app.Div().
								Class("share-item").
								Body(
									app.Label().Text("Player Link:"),
									app.Input().
										Type("text").
										ReadOnly(true).
										Value(g.getShareURL()).
										Class("share-input").
										OnClick(g.onSelectText),
									app.Button().
										Class("copy-btn").
										Text("📋").
										OnClick(g.onCopyShareURL),
								),
							app.Div().
								Class("share-item").
								Body(
									app.Label().Text("Spectator Link:"),
									app.Input().
										Type("text").
										ReadOnly(true).
										Value(g.getSpectatorURL()).
										Class("share-input").
										OnClick(g.onSelectText),
									app.Button().
										Class("copy-btn").
										Text("📋").
										OnClick(g.onCopySpectatorURL),
								),
						),
				),
		)
}

// Helper method to render feature cards
func (g *GolfGamezApp) renderFeatureCard(icon, title, description string) app.UI {
	return app.Div().
		Class("feature-card").
		Body(
			app.Div().
				Class("feature-icon").
				Text(icon),
			app.H3().
				Class("feature-title").
				Text(title),
			app.P().
				Class("feature-description").
				Text(description),
		)
}

// Active game view (scorecard)
func (g *GolfGamezApp) renderGameActive() app.UI {
	if g.currentGame == nil {
		return app.Div().Text("Loading game...")
	}

	return app.Main().
		Class("main-content game-active-view").
		Body(
			g.renderHoleHeader(),
			g.renderHoleScoreEntry(),
			g.renderHoleNavigation(),
			g.renderCurrentHoleScores(),
			g.renderLeaderboard(),
		)
}

// Scorecard rendering
func (g *GolfGamezApp) renderScorecard() app.UI {

	fmt.Println("Rendering Scorecard... ", g.currentGame, g.currentGame.CourseInfo)
	if g.currentGame == nil || g.currentGame.CourseInfo == nil {
		return app.Div()
	}

	return app.Div().
		Class("scorecard-section").
		Body(
			app.H3().Text("📊 Scorecard"),
			app.Div().
				Class("scorecard").
				Body(
					g.renderScorecardHeader(),
					app.Range(g.currentGame.Players).Slice(func(i int) app.UI {
						return g.renderPlayerScoreRow(g.currentGame.Players[i])
					}),
				),
		)
}

// Scorecard header with hole numbers
func (g *GolfGamezApp) renderScorecardHeader() app.UI {
	if g.currentGame.CourseInfo == nil {
		return app.Div()
	}

	return app.Div().
		Class("scorecard-header").
		Body(
			app.Div().Class("player-cell").Text("Player"),
			app.Range(g.currentGame.CourseInfo.Holes).Slice(func(i int) app.UI {
				hole := g.currentGame.CourseInfo.Holes[i]
				return app.Div().
					Class("hole-cell").
					Body(
						app.Div().Class("hole-number").Text(strconv.Itoa(hole.Hole)),
						app.Div().Class("hole-par").Text(fmt.Sprintf("Par %d", hole.Par)),
					)
			}),
			app.Div().Class("total-cell").Text("Total"),
		)
}

// Player score row
func (g *GolfGamezApp) renderPlayerScoreRow(player models.Player) app.UI {
	return app.Div().
		Class("scorecard-row").
		Body(
			app.Div().
				Class("player-cell").
				Text(player.Name),
			app.Range(g.currentGame.CourseInfo.Holes).Slice(func(i int) app.UI {
				hole := g.currentGame.CourseInfo.Holes[i]
				return g.renderScoreCell(player, hole.Hole)
			}),
			app.Div().
				Class("total-cell").
				Text(g.getPlayerTotal(player)),
		)
}

// Individual score cell
func (g *GolfGamezApp) renderScoreCell(player models.Player, holeNumber int) app.UI {
	// This would typically get the score from player stats or a separate scores array
	// For now, we'll show placeholder
	return app.Div().
		Class("score-cell").
		Body(
			app.Input().
				Type("number").
				Min("1").
				Max("15").
				Class("score-input").
				Placeholder("-").
				OnChange(func(ctx app.Context, e app.Event) {
					g.onScoreChange(ctx, player.ID, holeNumber, e.Get("target").Get("value").String())
				}),
		)
}

// Simple leaderboard
func (g *GolfGamezApp) renderLeaderboard() app.UI {
	return app.Div().
		Class("leaderboard-section").
		Body(
			app.H3().Text("🏆 Leaderboard"),
			app.Div().
				Class("leaderboard").
				Body(
					app.Range(g.currentGame.Players).Slice(func(i int) app.UI {
						player := g.currentGame.Players[i]
						return app.Div().
							Class("leaderboard-entry").
							Body(
								app.Div().Class("position").Text(fmt.Sprintf("%d", i+1)),
								app.Div().Class("player-name").Text(player.Name),
								app.Div().Class("score").Text(g.getPlayerTotal(player)),
							)
					}),
				),
		)
}

// Game results view
func (g *GolfGamezApp) renderGameResults() app.UI {
	return app.Main().
		Class("main-content game-results-view").
		Body(
			app.H2().Text("🏆 Game Results"),
			app.P().Text("Final results and side bet winners"),
			// TODO: Implement results display
		)
}

// Spectator view
func (g *GolfGamezApp) renderSpectatorView() app.UI {
	return app.Main().
		Class("main-content spectator-view").
		Body(
			app.H2().Text("👀 Spectator View"),
			app.P().Text("Live game viewing"),
			// TODO: Implement spectator-specific view
		)
}

// Helper methods
func (g *GolfGamezApp) getShareURL() string {
	if g.currentGame == nil {
		return ""
	}
	baseURL := app.Window().Get("location").Get("origin").String()
	token := g.apiClient.ExtractTokenFromShareLink(g.currentGame.ShareLink)
	return fmt.Sprintf("%s/game/%s", baseURL, token)
}

func (g *GolfGamezApp) getSpectatorURL() string {
	if g.currentGame == nil {
		return ""
	}
	baseURL := app.Window().Get("location").Get("origin").String()
	token := g.apiClient.ExtractTokenFromShareLink(g.currentGame.SpectatorLink)
	return fmt.Sprintf("%s/spectate/%s", baseURL, token)
}

func (g *GolfGamezApp) getPlayerTotal(player models.Player) string {
	// Placeholder for now - would calculate from actual scores
	if player.Stats != nil {
		return player.Stats.CurrentScore
	}
	return "E"
}

// Hole-by-hole scorecard render methods
func (g *GolfGamezApp) renderHoleHeader() app.UI {
	if g.currentGame == nil || g.currentGame.CourseInfo == nil {
		return app.Div()
	}

	maxHoles := len(g.currentGame.CourseInfo.Holes)
	if maxHoles == 0 {
		maxHoles = 18 // Default to 18 holes
	}

	// Get current hole info
	var holeInfo *models.HoleInfo
	if g.currentHole <= len(g.currentGame.CourseInfo.Holes) {
		holeInfo = &g.currentGame.CourseInfo.Holes[g.currentHole-1]
	}

	return app.Div().
		Class("hole-header").
		Body(
			app.Div().
				Class("hole-info-card").
				Body(
					app.H2().
						Class("hole-title").
						Text(fmt.Sprintf("Hole %d", g.currentHole)),
					app.Div().
						Class("hole-details").
						Body(
							app.If(holeInfo != nil, func() app.UI {
								return app.Div().
									Class("hole-stats").
									Body(
										app.Span().
											Class("hole-par").
											Text(fmt.Sprintf("Par %d", holeInfo.Par)),
										app.If(holeInfo.Yardage != nil, func() app.UI {
											return app.Span().
												Class("hole-yardage").
												Text(fmt.Sprintf("%d yards", *holeInfo.Yardage))
										}),
										app.If(holeInfo.HandicapRanking != nil && *holeInfo.HandicapRanking > 0, func() app.UI {
											return app.Span().
												Class("hole-handicap").
												Text(fmt.Sprintf("HCP %d", *holeInfo.HandicapRanking))
										}),
									)
							}),
							app.Div().
								Class("hole-progress").
								Body(
									app.Span().
										Class("hole-counter").
										Text(fmt.Sprintf("%d of %d", g.currentHole, maxHoles)),
									app.Div().
										Class("progress-bar").
										Body(
											app.Div().
												Class("progress-fill").
												Style("width", fmt.Sprintf("%.1f%%", float64(g.currentHole)/float64(maxHoles)*100)),
										),
								),
						),
				),
		)
}

func (g *GolfGamezApp) renderHoleScoreEntry() app.UI {
	if g.currentGame == nil || len(g.currentGame.Players) == 0 {
		return app.Div().
			Class("score-entry-placeholder").
			Body(
				app.P().Text("No players in game"),
			)
	}

	return app.Div().
		Class("score-entry-section").
		Body(
			app.H3().
				Class("score-entry-title").
				Text("Enter Scores"),
			app.Form().
				Class("score-entry-form").
				OnSubmit(g.onScoreSubmit).
				Body(
					app.Div().
						Class("form-group").
						Body(
							app.Label().
								For("player-select").
								Text("Player"),
							app.Select().
								ID("player-select").
								Class("form-control player-select").
								Required(true).
								Body(
									app.Option().
										Value("").
										Text("Select player..."),
									app.Range(g.currentGame.Players).Slice(func(i int) app.UI {
										player := g.currentGame.Players[i]
										return app.Option().
											Value(player.ID).
											Text(player.Name)
									}),
								),
						),
					app.Div().
						Class("score-inputs-row").
						Body(
							app.Div().
								Class("form-group").
								Body(
									app.Label().
										For("strokes").
										Text("Strokes"),
									app.Input().
										Type("number").
										ID("strokes").
										Class("form-control score-input").
										Min("1").
										Max("15").
										Required(true).
										Placeholder("1-15"),
								),
							app.Div().
								Class("form-group").
								Body(
									app.Label().
										For("putts").
										Text("Putts"),
									app.Input().
										Type("number").
										ID("putts").
										Class("form-control score-input").
										Min("0").
										Max("10").
										Required(true).
										Placeholder("0-10"),
								),
						),
					app.Div().
						Class("form-actions").
						Body(
							app.Button().
								Type("submit").
								Class("submit-score-btn primary-btn").
								Text("Record Score"),
						),
				),
		)
}

func (g *GolfGamezApp) renderHoleNavigation() app.UI {
	maxHoles := 18
	if g.currentGame != nil && g.currentGame.CourseInfo != nil {
		maxHoles = len(g.currentGame.CourseInfo.Holes)
	}

	return app.Div().
		Class("hole-navigation").
		Body(
			app.Button().
				Class("nav-btn prev-hole-btn").
				Disabled(g.currentHole <= 1).
				OnClick(g.onPreviousHole).
				Body(
					app.Span().Class("nav-icon").Text("←"),
					app.Span().Class("nav-text").Text("Previous"),
				),
			app.Div().
				Class("hole-selector").
				Body(
					app.Span().Text("Hole"),
					app.Select().
						Class("hole-select").
						OnChange(func(ctx app.Context, e app.Event) {
							holeStr := e.Get("target").Get("value").String()
							if hole, err := strconv.Atoi(holeStr); err == nil && hole >= 1 && hole <= maxHoles {
								g.currentHole = hole
							}
						}).
						Body(
							app.Range(make([]int, maxHoles)).Slice(func(i int) app.UI {
								holeNum := i + 1
								return app.Option().
									Value(strconv.Itoa(holeNum)).
									Text(strconv.Itoa(holeNum)).
									Selected(holeNum == g.currentHole)
							}),
						),
				),
			app.Button().
				Class("nav-btn next-hole-btn").
				Disabled(g.currentHole >= maxHoles).
				OnClick(g.onNextHole).
				Body(
					app.Span().Class("nav-text").Text("Next"),
					app.Span().Class("nav-icon").Text("→"),
				),
		)
}

func (g *GolfGamezApp) renderCurrentHoleScores() app.UI {
	if g.currentGame == nil || len(g.currentGame.Players) == 0 {
		return app.Div()
	}

	return app.Div().
		Class("current-hole-scores").
		Body(
			app.H3().
				Class("scores-title").
				Text(fmt.Sprintf("Hole %d Scores", g.currentHole)),
			app.Div().
				Class("scores-grid").
				Body(
					app.Range(g.currentGame.Players).Slice(func(i int) app.UI {
						player := g.currentGame.Players[i]
						return g.renderPlayerHoleScore(player)
					}),
				),
		)
}

func (g *GolfGamezApp) renderPlayerHoleScore(player models.Player) app.UI {
	// Find score for current hole - for now this is a placeholder
	// In the full implementation, scores would be fetched from API or stored in app state
	var holeScore *models.Score
	// TODO: Implement score retrieval from app state or API
	// For now, we'll show placeholder data

	return app.Div().
		Class("player-hole-score").
		Body(
			app.Div().
				Class("player-info").
				Body(
					app.Span().
						Class("player-name").
						Text(player.Name),
					app.If(player.Handicap > 0, func() app.UI {
						return app.Span().
							Class("player-handicap").
							Text(fmt.Sprintf("(%.1f)", player.Handicap))
					}),
				),
			app.Div().
				Class("score-display").
				Body(
					app.If(holeScore != nil, func() app.UI {
						return app.Div().
							Class("recorded-score").
							Body(
								app.Span().
									Class("strokes").
									Text(fmt.Sprintf("%d", holeScore.Strokes)),
								app.Span().
									Class("putts").
									Text(fmt.Sprintf("(%dp)", holeScore.Putts)),
							)
					}).Else(func() app.UI {
						return app.Span().
							Class("no-score").
							Text("--")
					}),
				),
		)
}
