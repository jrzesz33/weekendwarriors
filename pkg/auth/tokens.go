package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	// Token prefixes for different token types
	ShareTokenPrefix     = "gt_"
	SpectatorTokenPrefix = "st_"

	// Token lengths (excluding prefix)
	TokenLength = 20
)

// TokenType represents the type of access token
type TokenType string

const (
	TokenTypeShare     TokenType = "share"
	TokenTypeSpectator TokenType = "spectator"
)

// TokenPair represents a pair of tokens for a game
type TokenPair struct {
	ShareToken     string
	SpectatorToken string
}

// GenerateTokenPair generates a new pair of share and spectator tokens
func GenerateTokenPair() (*TokenPair, error) {
	shareToken, err := generateToken(ShareTokenPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to generate share token: %w", err)
	}

	spectatorToken, err := generateToken(SpectatorTokenPrefix)
	if err != nil {
		return nil, fmt.Errorf("failed to generate spectator token: %w", err)
	}

	return &TokenPair{
		ShareToken:     shareToken,
		SpectatorToken: spectatorToken,
	}, nil
}

// generateToken generates a cryptographically secure random token with the given prefix
func generateToken(prefix string) (string, error) {
	// Generate random bytes
	bytes := make([]byte, TokenLength/2) // Each byte becomes 2 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Convert to hex string and add prefix
	token := prefix + hex.EncodeToString(bytes)
	return token, nil
}

// ValidateTokenFormat validates that a token has the correct format
func ValidateTokenFormat(token string) (TokenType, error) {
	if strings.HasPrefix(token, ShareTokenPrefix) {
		if len(token) != len(ShareTokenPrefix)+TokenLength {
			return "", fmt.Errorf("invalid share token length")
		}
		return TokenTypeShare, nil
	}

	if strings.HasPrefix(token, SpectatorTokenPrefix) {
		if len(token) != len(SpectatorTokenPrefix)+TokenLength {
			return "", fmt.Errorf("invalid spectator token length")
		}
		return TokenTypeSpectator, nil
	}

	return "", fmt.Errorf("invalid token prefix")
}

// ExtractGameIDFromToken attempts to extract a game ID from a token
// This is used when tokens are passed as game IDs in URL paths
func ExtractGameIDFromToken(tokenOrGameID string) (string, TokenType) {
	// Check if it's a share token
	if strings.HasPrefix(tokenOrGameID, ShareTokenPrefix) {
		if tokenType, err := ValidateTokenFormat(tokenOrGameID); err == nil {
			return tokenOrGameID, tokenType
		}
	}

	// Check if it's a spectator token
	if strings.HasPrefix(tokenOrGameID, SpectatorTokenPrefix) {
		if tokenType, err := ValidateTokenFormat(tokenOrGameID); err == nil {
			return tokenOrGameID, tokenType
		}
	}

	// If not a token, assume it's a game ID
	return tokenOrGameID, ""
}

// GenerateGameID generates a unique game ID
func GenerateGameID() (string, error) {
	return generateToken("game_")
}

// GeneratePlayerID generates a unique player ID
func GeneratePlayerID() (string, error) {
	return generateToken("player_")
}

// GenerateScoreID generates a unique score ID
func GenerateScoreID() (string, error) {
	return generateToken("score_")
}

// GenerateSideBetID generates a unique side bet calculation ID
func GenerateSideBetID() (string, error) {
	return generateToken("sidebet_")
}

// GeneratePokerCardID generates a unique poker card event ID
func GeneratePokerCardID() (string, error) {
	return generateToken("card_")
}

// GeneratePokerHandID generates a unique poker hand ID
func GeneratePokerHandID() (string, error) {
	return generateToken("hand_")
}