//go:build js && wasm
// +build js,wasm

package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"syscall/js"
	"time"

	"golf-gamez-frontend/internal/models"
)

// Client handles API communication with the Golf Gamez backend
type Client struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAuthToken sets the authentication token for API requests
func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

// CreateGame creates a new golf game
func (c *Client) CreateGame(request models.CreateGameRequest) (*models.Game, error) {
	url := c.baseURL + "/games"

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.makeRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	var game models.Game
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Store the share token for future requests
	if game.ShareLink != "" {
		// Extract token from share link
		// This is a simplified extraction - in reality you'd parse the URL properly
		c.authToken = game.ID
	}

	return &game, nil
}

// GetGame retrieves game information
func (c *Client) GetGame(gameID string) (*models.Game, error) {
	url := c.baseURL + "/games/" + gameID

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var game models.Game
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &game, nil
}

// StartGame starts a golf game
func (c *Client) StartGame(gameID string) error {
	url := c.baseURL + "/games/" + gameID + "/start"

	resp, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// CompleteGame completes a golf game
func (c *Client) CompleteGame(gameID string) (*models.GameCompletionResult, error) {
	url := c.baseURL + "/games/" + gameID + "/complete"

	resp, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var result models.GameCompletionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// AddPlayer adds a player to a game
func (c *Client) AddPlayer(gameID string, request models.CreatePlayerRequest) (*models.Player, error) {
	url := c.baseURL + "/games/" + gameID + "/players"

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.makeRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	var player models.Player
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &player, nil
}

// UpdatePlayer updates player information
func (c *Client) UpdatePlayer(gameID, playerID string, name string, handicap float64) (*models.Player, error) {
	url := c.baseURL + "/games/" + gameID + "/players/" + playerID

	request := map[string]interface{}{
		"name":     name,
		"handicap": handicap,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.makeRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var player models.Player
	if err := json.NewDecoder(resp.Body).Decode(&player); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &player, nil
}

// RemovePlayer removes a player from a game
func (c *Client) RemovePlayer(gameID, playerID string) error {
	url := c.baseURL + "/games/" + gameID + "/players/" + playerID

	resp, err := c.makeRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// RecordScore records a score for a player
func (c *Client) RecordScore(gameID, playerID string, request models.ScoreRequest) (*models.Score, error) {
	url := c.baseURL + "/games/" + gameID + "/players/" + playerID + "/scores"

	data, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.makeRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	var score models.Score
	if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &score, nil
}

// UpdateScore updates an existing score
func (c *Client) UpdateScore(gameID, playerID string, hole, strokes, putts int) error {
	url := fmt.Sprintf("%s/games/%s/players/%s/scores/%d", c.baseURL, gameID, playerID, hole)

	request := map[string]interface{}{
		"strokes": strokes,
		"putts":   putts,
	}

	data, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.makeRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp)
	}

	return nil
}

// GetScorecard retrieves the game scorecard
func (c *Client) GetScorecard(gameID string) (*models.GameScorecard, error) {
	url := c.baseURL + "/games/" + gameID + "/scorecard"

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var scorecard models.GameScorecard
	if err := json.NewDecoder(resp.Body).Decode(&scorecard); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &scorecard, nil
}

// GetLeaderboard retrieves the current leaderboard
func (c *Client) GetLeaderboard(gameID string) (*models.Leaderboard, error) {
	url := c.baseURL + "/games/" + gameID + "/leaderboard"

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var leaderboard models.Leaderboard
	if err := json.NewDecoder(resp.Body).Decode(&leaderboard); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &leaderboard, nil
}

// GetBestNineStandings retrieves Best Nine side bet standings
func (c *Client) GetBestNineStandings(gameID string) (*models.BestNineStandings, error) {
	url := c.baseURL + "/games/" + gameID + "/side-bets/best-nine"

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var standings models.BestNineStandings
	if err := json.NewDecoder(resp.Body).Decode(&standings); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &standings, nil
}

// GetPuttPuttPokerStatus retrieves Putt Putt Poker side bet status
func (c *Client) GetPuttPuttPokerStatus(gameID string) (*models.PuttPuttPokerStatus, error) {
	url := c.baseURL + "/games/" + gameID + "/side-bets/putt-putt-poker"

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var status models.PuttPuttPokerStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &status, nil
}

// DealPokerCards deals final poker cards
func (c *Client) DealPokerCards(gameID string) (*models.PokerDealResult, error) {
	url := c.baseURL + "/games/" + gameID + "/side-bets/putt-putt-poker/deal"

	resp, err := c.makeRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, c.handleErrorResponse(resp)
	}

	var result models.PokerDealResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetSpectatorView retrieves spectator view
func (c *Client) GetSpectatorView(token string) (*models.SpectatorView, error) {
	url := c.baseURL + "/spectate/" + token

	resp, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var view models.SpectatorView
	if err := json.NewDecoder(resp.Body).Decode(&view); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &view, nil
}

// makeRequest makes an HTTP request with proper headers and authentication
func (c *Client) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
	// For WebAssembly, we need to use the Fetch API instead of Go's http.Client
	return c.makeFetchRequest(method, url, body)
}

// makeFetchRequest uses the JavaScript Fetch API for HTTP requests
func (c *Client) makeFetchRequest(method, url string, body io.Reader) (*http.Response, error) {
	// Convert body to string if present
	var bodyStr string
	if body != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(body)
		bodyStr = buf.String()
	}

	// Create fetch options
	fetchOptions := map[string]interface{}{
		"method": method,
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
	}

	// Add authorization header if token is set
	if c.authToken != "" {
		headers := fetchOptions["headers"].(map[string]interface{})
		headers["Authorization"] = "Bearer " + c.authToken
	}

	// Add body if present
	if bodyStr != "" {
		fetchOptions["body"] = bodyStr
	}

	// Make the fetch request
	promise := js.Global().Call("fetch", url, fetchOptions)

	// Convert Promise to synchronous operation (for simplicity in this example)
	// In a real implementation, you'd want to handle this asynchronously
	return c.waitForFetchResponse(promise)
}

// waitForFetchResponse waits for the fetch promise to resolve
func (c *Client) waitForFetchResponse(promise js.Value) (*http.Response, error) {
	fmt.Println("Waiting for Fetch Response")
	// Create a channel to wait for the promise to resolve
	done := make(chan struct{})
	var response *http.Response
	var err error
	fmt.Println("Setting up Promise Handlers")
	// Set up promise handlers
	successHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer close(done)

		if len(args) == 0 {
			err = fmt.Errorf("no response received")
			return nil
		}

		fetchResponse := args[0]

		// Get response status
		status := fetchResponse.Get("status").Int()

		// Get response headers
		headers := make(http.Header)

		// Get response body as text
		textPromise := fetchResponse.Call("text")
		textDone := make(chan string, 1)

		textHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) > 0 {
				textDone <- args[0].String()
			} else {
				textDone <- ""
			}
			return nil
		})

		textPromise.Call("then", textHandler)
		bodyText := <-textDone

		response = &http.Response{
			StatusCode: status,
			Header:     headers,
			Body:       io.NopCloser(bytes.NewReader([]byte(bodyText))),
		}

		return nil
	})

	errorHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer close(done)
		if len(args) > 0 {
			err = fmt.Errorf("fetch error: %s", args[0].String())
		} else {
			err = fmt.Errorf("unknown fetch error")
		}
		return nil
	})

	// Set up promise handlers
	promise.Call("then", successHandler, errorHandler)
	fmt.Println("Waiting for the Handler to Resolve")
	// Wait for the promise to resolve
	<-done
	fmt.Println("Response Sent")
	return response, err
}

// handleErrorResponse handles API error responses
func (c *Client) handleErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP %d: failed to read error response", resp.StatusCode)
	}

	var errorResp models.ErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return fmt.Errorf("API error (%s): %s", errorResp.Error.Code, errorResp.Error.Message)
}

// GetDiamondRunCourseInfo returns the course information for Diamond Run
func (c *Client) GetDiamondRunCourseInfo() *models.CourseInfo {
	holes := make([]models.HoleInfo, 18)
	pars := []int{4, 3, 5, 4, 3, 4, 4, 3, 5, 4, 3, 5, 3, 4, 4, 4, 5, 4}

	for i := 0; i < 18; i++ {
		holes[i] = models.HoleInfo{
			Hole:            i + 1,
			Par:             pars[i],
			HandicapRanking: i + 1, // Simplified ranking
			Yardage:         350,   // Default yardage
			Description:     fmt.Sprintf("Hole %d", i+1),
		}
	}

	return &models.CourseInfo{
		Name:     "Diamond Run",
		Holes:    holes,
		TotalPar: 72,
	}
}
