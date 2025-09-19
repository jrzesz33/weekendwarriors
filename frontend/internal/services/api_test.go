package services

import (
	"context"
	"testing"

	"golf-gamez-frontend/internal/models"
	"golf-gamez-frontend/internal/utils"
)

func TestAPIClient_CreateGame(t *testing.T) {
	// Setup test helper
	th := utils.NewTestHelper(t)
	defer th.Cleanup()

	// Create mock server
	server := th.MockAPIServer()

	// Create API client
	client := NewAPIClient(server.URL)

	// Test game creation
	req := models.CreateGameRequest{
		Course:          "diamond-run",
		SideBets:        []models.SideBetType{models.SideBetBestNine},
		HandicapEnabled: true,
	}

	ctx := utils.CreateTestContext()
	game, err := client.CreateGame(ctx, req)

	if err != nil {
		t.Fatalf("CreateGame failed: %v", err)
	}

	if game == nil {
		t.Fatal("Expected game to be created, got nil")
	}

	if game.Course != "diamond-run" {
		t.Errorf("Expected course 'diamond-run', got '%s'", game.Course)
	}

	if game.Status != models.GameStatusSetup {
		t.Errorf("Expected status 'setup', got '%s'", game.Status)
	}
}

func TestAPIClient_ExtractTokenFromShareLink(t *testing.T) {
	client := NewAPIClient("")

	tests := []struct {
		name      string
		shareLink string
		expected  string
	}{
		{
			name:      "Simple share link",
			shareLink: "/games/gt_abc123",
			expected:  "gt_abc123",
		},
		{
			name:      "Full URL share link",
			shareLink: "https://example.com/games/gt_xyz789",
			expected:  "gt_xyz789",
		},
		{
			name:      "No token",
			shareLink: "/games/invalid",
			expected:  "",
		},
		{
			name:      "Empty string",
			shareLink: "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.ExtractTokenFromShareLink(tt.shareLink)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestAPIClient_HealthCheck(t *testing.T) {
	// Setup test helper
	th := utils.NewTestHelper(t)
	defer th.Cleanup()

	// Create mock server
	server := th.MockAPIServer()

	// Create API client
	client := NewAPIClient(server.URL)

	// Test health check
	ctx := utils.CreateTestContext()
	err := client.HealthCheck(ctx)

	if err != nil {
		t.Fatalf("HealthCheck failed: %v", err)
	}
}

func TestAPIClient_SetAuthToken(t *testing.T) {
	client := NewAPIClient("")

	token := "gt_test123"
	client.SetAuthToken(token)

	if client.authToken != token {
		t.Errorf("Expected auth token '%s', got '%s'", token, client.authToken)
	}
}

// Benchmark tests
func BenchmarkAPIClient_ExtractTokenFromShareLink(b *testing.B) {
	client := NewAPIClient("")
	shareLink := "/games/gt_abc123def456"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.ExtractTokenFromShareLink(shareLink)
	}
}

// Example test showing usage
func ExampleAPIClient_CreateGame() {
	// Create API client
	client := NewAPIClient("http://localhost:8080/v1")

	// Create a new game
	req := models.CreateGameRequest{
		Course:          "diamond-run",
		SideBets:        []models.SideBetType{models.SideBetBestNine},
		HandicapEnabled: true,
	}

	ctx := context.Background()
	game, err := client.CreateGame(ctx, req)
	if err != nil {
		// Handle error
		return
	}

	// Extract and use the game token
	token := client.ExtractTokenFromShareLink(game.ShareLink)
	client.SetAuthToken(token)

	// Now you can use the client for subsequent API calls
	_ = game
}