package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// GolfGamezApp is the main application component
type GolfGamezApp struct {
	app.Compo
}

func (g *GolfGamezApp) Render() app.UI {
	fmt.Println("Rendering Golf Gamez App")

	return app.Div().
		Class("golf-app").
		Body(
			app.Header().
				Class("header").
				Body(
					app.H1().
						Class("title").
						Text("⛳ Golf Gamez"),
					app.P().
						Class("subtitle").
						Text("Track scores • Side bets • Real-time updates"),
				),

			app.Main().
				Class("main-content").
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
								Class("create-game-btn").
								Text("Create New Game").
								OnClick(g.onCreateGame),
							app.Button().
								Class("join-game-btn").
								Text("Join Game").
								OnClick(g.onJoinGame),
						),
				),

			app.Footer().
				Class("footer").
				Body(
					app.P().Text("Built with Go + WebAssembly • PWA Technology"),
				),
		)
}

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

func (g *GolfGamezApp) onCreateGame(ctx app.Context, e app.Event) {
	fmt.Println("Create game clicked")
	// TODO: Implement game creation
	app.Window().Get("alert").Invoke("Game creation coming soon!")
}

func (g *GolfGamezApp) onJoinGame(ctx app.Context, e app.Event) {
	fmt.Println("Join game clicked")
	// TODO: Implement game joining
	app.Window().Get("alert").Invoke("Game joining coming soon!")
}

func main() {
	fmt.Println("Starting the Go-App Composer")
	// Route the main app component
	app.Route("/", func() app.Composer { return &GolfGamezApp{} })
	fmt.Println("Running when on Browser")
	// Run the app when in browser (WebAssembly)
	app.RunWhenOnBrowser()
}
