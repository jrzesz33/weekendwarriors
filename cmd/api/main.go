package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golf-gamez/internal/config"
	"golf-gamez/internal/database"
	"golf-gamez/internal/handlers"
	"golf-gamez/internal/middleware"
	"golf-gamez/internal/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure structured logging
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Run database migrations
	if err := database.Migrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize services
	gameService := services.NewGameService(db)
	playerService := services.NewPlayerService(db)
	scoreService := services.NewScoreService(db)
	sideBetService := services.NewSideBetService(db)
	websocketService := services.NewWebSocketService()

	// Initialize handlers
	gameHandler := handlers.NewGameHandler(gameService, websocketService)
	playerHandler := handlers.NewPlayerHandler(playerService, websocketService)
	scoreHandler := handlers.NewScoreHandler(scoreService, sideBetService, websocketService)
	sideBetHandler := handlers.NewSideBetHandler(sideBetService, websocketService)
	spectatorHandler := handlers.NewSpectatorHandler(gameService)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Rate limiting
	r.Use(middleware.RateLimit(1000, time.Hour))    // Global rate limit
	r.Use(middleware.GameCreationLimit(10, time.Hour)) // Game creation limit

	// API routes
	r.Route("/v1", func(r chi.Router) {
		// Health check
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok"}`))
		})

		// Game management routes
		r.Route("/games", func(r chi.Router) {
			r.Post("/", gameHandler.CreateGame)

			r.Route("/{gameId}", func(r chi.Router) {
				r.Use(middleware.GameAuth(db))
				r.Get("/", gameHandler.GetGame)
				r.Delete("/", gameHandler.DeleteGame)
				r.Post("/start", gameHandler.StartGame)
				r.Post("/complete", gameHandler.CompleteGame)

				// Player management
				r.Route("/players", func(r chi.Router) {
					r.Get("/", playerHandler.GetPlayers)
					r.Post("/", playerHandler.AddPlayer)

					r.Route("/{playerId}", func(r chi.Router) {
						r.Get("/", playerHandler.GetPlayer)
						r.Put("/", playerHandler.UpdatePlayer)
						r.Delete("/", playerHandler.RemovePlayer)

						// Score management
						r.Route("/scores", func(r chi.Router) {
							r.Post("/", scoreHandler.RecordScore)
							r.Put("/{hole}", scoreHandler.UpdateScore)
						})
					})
				})

				// Game data routes
				r.Get("/scorecard", scoreHandler.GetGameScorecard)
				r.Get("/leaderboard", scoreHandler.GetLeaderboard)

				// Side bet routes
				r.Route("/side-bets", func(r chi.Router) {
					r.Get("/best-nine", sideBetHandler.GetBestNineStandings)
					r.Get("/putt-putt-poker", sideBetHandler.GetPuttPuttPokerStatus)
					r.Post("/putt-putt-poker/deal", sideBetHandler.DealPokerCards)
				})
			})
		})

		// Spectator routes
		r.Route("/spectate", func(r chi.Router) {
			r.Get("/{spectatorToken}", spectatorHandler.SpectateGame)
		})

		// WebSocket endpoint
		r.HandleFunc("/ws/games/{gameId}", websocketService.HandleWebSocket)
	})

	// Start server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info().Int("port", cfg.Port).Msg("Starting Golf Gamez API server")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exiting")
}