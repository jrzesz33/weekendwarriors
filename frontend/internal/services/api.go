package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golf-gamez-frontend/internal/models"
)

const (
	DefaultAPIBaseURL = "http://localhost:8080/v1"
	DefaultTimeout    = 30 * time.Second
)

// APIClient handles all communication with the Golf Gamez backend API
type APIClient struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// NewAPIClient creates a new API client instance
func NewAPIClient(baseURL string) *APIClient {
	if baseURL == "" {
		baseURL = DefaultAPIBaseURL
	}

	return &APIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// SetAuthToken sets the authentication token for API requests
func (c *APIClient) SetAuthToken(token string) {
	c.authToken = token
}

// ExtractTokenFromShareLink extracts the game token from a share link
func (c *APIClient) ExtractTokenFromShareLink(shareLink string) string {
	// Share link format: "/games/gt_xxxxx" or full URL
	parts := strings.Split(shareLink, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "gt_") {
			return part
		}
	}
	return ""
}

// CreateGame creates a new golf game
func (c *APIClient) CreateGame(ctx context.Context, req models.CreateGameRequest) (*models.Game, error) {
	var game models.Game
	err := c.makeRequest(ctx, "POST", "/games", req, &game)
	if err != nil {
		return nil, err
	}

	// Extract and set the auth token from the share link
	token := c.ExtractTokenFromShareLink(game.ShareLink)
	if token != "" {
		c.SetAuthToken(token)
	}

	return &game, nil
}

// GetGame retrieves game information by token
func (c *APIClient) GetGame(ctx context.Context, token string) (*models.Game, error) {
	var game models.Game
	err := c.makeRequest(ctx, "GET", fmt.Sprintf("/games/%s", token), nil, &game)
	return &game, err
}

// AddPlayer adds a new player to a game
func (c *APIClient) AddPlayer(ctx context.Context, token string, req models.CreatePlayerRequest) (*models.Player, error) {
	var player models.Player
	err := c.makeRequest(ctx, "POST", fmt.Sprintf("/games/%s/players", token), req, &player)
	return &player, err
}

// StartGame starts a game session
func (c *APIClient) StartGame(ctx context.Context, token string) (*models.Game, error) {
	var game models.Game
	err := c.makeRequest(ctx, "POST", fmt.Sprintf("/games/%s/start", token), nil, &game)
	return &game, err
}

// RecordScore records a score for a player on a specific hole
func (c *APIClient) RecordScore(ctx context.Context, token, playerID string, req models.CreateScoreRequest) (*models.Score, error) {
	var score models.Score
	err := c.makeRequest(ctx, "POST", fmt.Sprintf("/games/%s/players/%s/scores", token, playerID), req, &score)
	return &score, err
}

// GetLeaderboard retrieves the current leaderboard for a game
func (c *APIClient) GetLeaderboard(ctx context.Context, token string) (*models.FinalResults, error) {
	var leaderboard models.FinalResults
	err := c.makeRequest(ctx, "GET", fmt.Sprintf("/games/%s/leaderboard", token), nil, &leaderboard)
	return &leaderboard, err
}

// GetBestNineStandings retrieves the best nine side bet standings
func (c *APIClient) GetBestNineStandings(ctx context.Context, token string) (*models.BestNineStandings, error) {
	var standings models.BestNineStandings
	err := c.makeRequest(ctx, "GET", fmt.Sprintf("/games/%s/side-bets/best-nine", token), nil, &standings)
	return &standings, err
}

// GetPuttPuttPokerStandings retrieves the putt putt poker standings
func (c *APIClient) GetPuttPuttPokerStandings(ctx context.Context, token string) (*models.PokerFinalResults, error) {
	var standings models.PokerFinalResults
	err := c.makeRequest(ctx, "GET", fmt.Sprintf("/games/%s/side-bets/putt-putt-poker", token), nil, &standings)
	return &standings, err
}

// GetSpectatorView retrieves read-only game data for spectators
func (c *APIClient) GetSpectatorView(ctx context.Context, spectatorToken string) (*models.Game, error) {
	var game models.Game
	err := c.makeRequest(ctx, "GET", fmt.Sprintf("/spectate/%s", spectatorToken), nil, &game)
	return &game, err
}

// makeRequest is a helper method to make HTTP requests to the API
func (c *APIClient) makeRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	url := c.baseURL + endpoint

	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add auth token if available
	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr models.APIError
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}
		return fmt.Errorf("API error: %s - %s", apiErr.Code, apiErr.Message)
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// HealthCheck checks if the API is available
func (c *APIClient) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("health check failed with status %d", resp.StatusCode)
	}

	return nil
}