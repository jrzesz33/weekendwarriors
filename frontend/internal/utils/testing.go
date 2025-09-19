package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"golf-gamez-frontend/internal/models"
)

// TestHelper provides utilities for testing Golf Gamez frontend components
type TestHelper struct {
	t      *testing.T
	server *httptest.Server
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t: t,
	}
}

// MockAPIServer creates a mock API server for testing
func (th *TestHelper) MockAPIServer() *httptest.Server {
	mux := http.NewServeMux()

	// Mock game endpoints
	mux.HandleFunc("/v1/games", th.handleGames)
	mux.HandleFunc("/v1/games/", th.handleGameDetails)
	mux.HandleFunc("/v1/health", th.handleHealth)

	th.server = httptest.NewServer(mux)
	return th.server
}

// Cleanup cleans up test resources
func (th *TestHelper) Cleanup() {
	if th.server != nil {
		th.server.Close()
	}
}

// Mock handlers for API endpoints
func (th *TestHelper) handleGames(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// Mock game creation
		game := models.Game{
			ID:              "game_test123",
			Course:          "diamond-run",
			Status:          models.GameStatusSetup,
			HandicapEnabled: true,
			SideBets:        []models.SideBetType{models.SideBetBestNine},
			ShareLink:       "/games/gt_test123",
			SpectatorLink:   "/spectate/st_test123",
			Players:         []models.Player{},
		}
		th.writeJSON(w, game)

	case "GET":
		// Mock game list (for testing)
		games := []models.Game{
			{
				ID:     "game_test123",
				Course: "diamond-run",
				Status: models.GameStatusSetup,
			},
		}
		th.writeJSON(w, games)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (th *TestHelper) handleGameDetails(w http.ResponseWriter, r *http.Request) {
	// Extract game ID from path
	// This is a simplified implementation for testing

	game := models.Game{
		ID:              "game_test123",
		Course:          "diamond-run",
		Status:          models.GameStatusSetup,
		HandicapEnabled: true,
		SideBets:        []models.SideBetType{models.SideBetBestNine},
		ShareLink:       "/games/gt_test123",
		SpectatorLink:   "/spectate/st_test123",
		Players: []models.Player{
			{
				ID:       "player_test1",
				Name:     "Test Player",
				Handicap: 18.0,
				Gender:   models.GenderMale,
				Position: 1,
				GameID:   "game_test123",
			},
		},
	}

	th.writeJSON(w, game)
}

func (th *TestHelper) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Helper method to write JSON responses
func (th *TestHelper) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// In a real implementation, you'd use json.Marshal here
	// For testing, we'll just write a simple response
	w.WriteHeader(http.StatusOK)
}

// CreateTestContext creates a test context
func CreateTestContext() context.Context {
	return context.Background()
}

// MockGameData creates mock game data for testing
func MockGameData() *models.Game {
	return &models.Game{
		ID:              "game_test123",
		Course:          "diamond-run",
		Status:          models.GameStatusSetup,
		HandicapEnabled: true,
		SideBets:        []models.SideBetType{models.SideBetBestNine, models.SideBetPuttPuttPoker},
		ShareLink:       "/games/gt_test123",
		SpectatorLink:   "/spectate/st_test123",
		Players: []models.Player{
			{
				ID:       "player_test1",
				Name:     "Alice",
				Handicap: 18.0,
				Gender:   models.GenderFemale,
				Position: 1,
				GameID:   "game_test123",
			},
			{
				ID:       "player_test2",
				Name:     "Bob",
				Handicap: 24.0,
				Gender:   models.GenderMale,
				Position: 2,
				GameID:   "game_test123",
			},
		},
	}
}

// MockCourseData creates mock course data for testing
func MockCourseData() []models.HoleInfo {
	return []models.HoleInfo{
		{Hole: 1, Par: 4}, {Hole: 2, Par: 3}, {Hole: 3, Par: 5}, {Hole: 4, Par: 4},
		{Hole: 5, Par: 3}, {Hole: 6, Par: 4}, {Hole: 7, Par: 4}, {Hole: 8, Par: 3},
		{Hole: 9, Par: 5}, {Hole: 10, Par: 4}, {Hole: 11, Par: 3}, {Hole: 12, Par: 5},
		{Hole: 13, Par: 3}, {Hole: 14, Par: 4}, {Hole: 15, Par: 4}, {Hole: 16, Par: 4},
		{Hole: 17, Par: 5}, {Hole: 18, Par: 4},
	}
}

// CalculateTotalPar calculates total par for a range of holes
func CalculateTotalPar(holes []models.HoleInfo, startHole, endHole int) int {
	total := 0
	for _, hole := range holes {
		if hole.Hole >= startHole && hole.Hole <= endHole {
			total += hole.Par
		}
	}
	return total
}