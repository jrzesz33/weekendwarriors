package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"golf-gamez/pkg/auth"
	"golf-gamez/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// GameAuthContextKey is the context key for game authentication data
type GameAuthContextKey string

const (
	GameAuthKey GameAuthContextKey = "game_auth"
)

// GameAuth represents authentication context for a game
type GameAuthContext struct {
	GameID    string
	Token     string
	TokenType auth.TokenType
	IsShare   bool
	IsSpectator bool
}

// GameAuth middleware validates game access tokens
func GameAuth(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gameID := chi.URLParam(r, "gameId")
			if gameID == "" {
				apiErr := errors.New(errors.ErrValidation, "Game ID is required")
				errors.WriteHTTPError(w, apiErr)
				return
			}

			// Extract token from gameID or Authorization header
			var token string
			var tokenType auth.TokenType

			// Check if gameID is actually a token
			if actualGameID, tType := auth.ExtractGameIDFromToken(gameID); tType != "" {
				token = actualGameID
				tokenType = tType
			} else {
				// Look for token in Authorization header
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					apiErr := errors.New(errors.ErrInvalidToken, "Authorization token required")
					errors.WriteHTTPError(w, apiErr)
					return
				}

				// Extract Bearer token
				parts := strings.Split(authHeader, " ")
				if len(parts) != 2 || parts[0] != "Bearer" {
					apiErr := errors.New(errors.ErrInvalidToken, "Invalid authorization format")
					errors.WriteHTTPError(w, apiErr)
					return
				}

				token = parts[1]
				var err error
				tokenType, err = auth.ValidateTokenFormat(token)
				if err != nil {
					apiErr := errors.New(errors.ErrInvalidToken, "Invalid token format")
					errors.WriteHTTPError(w, apiErr)
					return
				}
			}

			// Validate token against database
			actualGameID, err := validateGameToken(db, token, tokenType)
			if err != nil {
				log.Warn().Str("token", token).Err(err).Msg("Invalid game token")
				apiErr := errors.New(errors.ErrInvalidToken, "Invalid or expired token")
				errors.WriteHTTPError(w, apiErr)
				return
			}

			// Create auth context
			authCtx := &GameAuthContext{
				GameID:      actualGameID,
				Token:       token,
				TokenType:   tokenType,
				IsShare:     tokenType == auth.TokenTypeShare,
				IsSpectator: tokenType == auth.TokenTypeSpectator,
			}

			// Check permissions for request method
			if !hasPermission(authCtx, r.Method) {
				apiErr := errors.New(errors.ErrInsufficientPermissions, "Insufficient permissions for this action")
				errors.WriteHTTPError(w, apiErr)
				return
			}

			// Add to context
			ctx := context.WithValue(r.Context(), GameAuthKey, authCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// validateGameToken validates a token against the database and returns the game ID
func validateGameToken(db *sql.DB, token string, tokenType auth.TokenType) (string, error) {
	var gameID string
	var query string

	switch tokenType {
	case auth.TokenTypeShare:
		query = "SELECT id FROM games WHERE share_token = ? AND status != 'abandoned'"
	case auth.TokenTypeSpectator:
		query = "SELECT id FROM games WHERE spectator_token = ? AND status != 'abandoned'"
	default:
		return "", errors.New(errors.ErrInvalidToken, "Invalid token type")
	}

	err := db.QueryRow(query, token).Scan(&gameID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New(errors.ErrGameNotFound, "Game not found")
		}
		return "", err
	}

	return gameID, nil
}

// hasPermission checks if the auth context has permission for the HTTP method
func hasPermission(authCtx *GameAuthContext, method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		// Both share and spectator tokens can read
		return true
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		// Only share tokens can modify
		return authCtx.IsShare
	default:
		return false
	}
}

// GetGameAuthFromContext extracts game auth context from request context
func GetGameAuthFromContext(ctx context.Context) (*GameAuthContext, bool) {
	authCtx, ok := ctx.Value(GameAuthKey).(*GameAuthContext)
	return authCtx, ok
}

// RequireShareToken middleware that requires a share token (write access)
func RequireShareToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCtx, ok := GetGameAuthFromContext(r.Context())
		if !ok {
			apiErr := errors.New(errors.ErrInvalidToken, "Authentication required")
			errors.WriteHTTPError(w, apiErr)
			return
		}

		if !authCtx.IsShare {
			apiErr := errors.New(errors.ErrInsufficientPermissions, "Share token required for this action")
			errors.WriteHTTPError(w, apiErr)
			return
		}

		next.ServeHTTP(w, r)
	})
}