package storage

import (
	"encoding/json"
	"syscall/js"

	"golf-gamez-frontend/internal/models"
)

// LocalStorage handles browser localStorage operations
type LocalStorage struct{}

// NewLocalStorage creates a new localStorage manager
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{}
}

// SaveGame saves game data to localStorage
func (ls *LocalStorage) SaveGame(game *models.Game) error {
	data, err := json.Marshal(game)
	if err != nil {
		return err
	}

	js.Global().Get("localStorage").Call("setItem", "currentGame", string(data))
	return nil
}

// LoadGame loads game data from localStorage
func (ls *LocalStorage) LoadGame() (*models.Game, error) {
	data := js.Global().Get("localStorage").Call("getItem", "currentGame")
	if !data.Truthy() || data.IsNull() {
		return nil, nil
	}

	var game models.Game
	if err := json.Unmarshal([]byte(data.String()), &game); err != nil {
		return nil, err
	}

	return &game, nil
}

// SavePlayerPreferences saves player preferences
func (ls *LocalStorage) SavePlayerPreferences(prefs map[string]interface{}) error {
	data, err := json.Marshal(prefs)
	if err != nil {
		return err
	}

	js.Global().Get("localStorage").Call("setItem", "playerPreferences", string(data))
	return nil
}

// LoadPlayerPreferences loads player preferences
func (ls *LocalStorage) LoadPlayerPreferences() (map[string]interface{}, error) {
	data := js.Global().Get("localStorage").Call("getItem", "playerPreferences")
	if !data.Truthy() || data.IsNull() {
		return make(map[string]interface{}), nil
	}

	var prefs map[string]interface{}
	if err := json.Unmarshal([]byte(data.String()), &prefs); err != nil {
		return nil, err
	}

	return prefs, nil
}

// SaveRecentGames saves recent game IDs
func (ls *LocalStorage) SaveRecentGames(gameIDs []string) error {
	data, err := json.Marshal(gameIDs)
	if err != nil {
		return err
	}

	js.Global().Get("localStorage").Call("setItem", "recentGames", string(data))
	return nil
}

// LoadRecentGames loads recent game IDs
func (ls *LocalStorage) LoadRecentGames() ([]string, error) {
	data := js.Global().Get("localStorage").Call("getItem", "recentGames")
	if !data.Truthy() || data.IsNull() {
		return []string{}, nil
	}

	var gameIDs []string
	if err := json.Unmarshal([]byte(data.String()), &gameIDs); err != nil {
		return nil, err
	}

	return gameIDs, nil
}

// AddRecentGame adds a game to the recent games list
func (ls *LocalStorage) AddRecentGame(gameID string) error {
	recentGames, err := ls.LoadRecentGames()
	if err != nil {
		return err
	}

	// Remove if already exists
	for i, id := range recentGames {
		if id == gameID {
			recentGames = append(recentGames[:i], recentGames[i+1:]...)
			break
		}
	}

	// Add to front
	recentGames = append([]string{gameID}, recentGames...)

	// Keep only last 10
	if len(recentGames) > 10 {
		recentGames = recentGames[:10]
	}

	return ls.SaveRecentGames(recentGames)
}

// ClearGame removes game data from localStorage
func (ls *LocalStorage) ClearGame() {
	js.Global().Get("localStorage").Call("removeItem", "currentGame")
}

// SaveOfflineScores saves scores for offline mode
func (ls *LocalStorage) SaveOfflineScores(scores []models.Score) error {
	data, err := json.Marshal(scores)
	if err != nil {
		return err
	}

	js.Global().Get("localStorage").Call("setItem", "offlineScores", string(data))
	return nil
}

// LoadOfflineScores loads scores from offline storage
func (ls *LocalStorage) LoadOfflineScores() ([]models.Score, error) {
	data := js.Global().Get("localStorage").Call("getItem", "offlineScores")
	if !data.Truthy() || data.IsNull() {
		return []models.Score{}, nil
	}

	var scores []models.Score
	if err := json.Unmarshal([]byte(data.String()), &scores); err != nil {
		return nil, err
	}

	return scores, nil
}

// ClearOfflineScores removes offline scores
func (ls *LocalStorage) ClearOfflineScores() {
	js.Global().Get("localStorage").Call("removeItem", "offlineScores")
}

// SaveConnectionState saves connection state
func (ls *LocalStorage) SaveConnectionState(connected bool) {
	js.Global().Get("localStorage").Call("setItem", "connectionState", connected)
}

// LoadConnectionState loads connection state
func (ls *LocalStorage) LoadConnectionState() bool {
	data := js.Global().Get("localStorage").Call("getItem", "connectionState")
	if !data.Truthy() || data.IsNull() {
		return false
	}

	return data.Bool()
}