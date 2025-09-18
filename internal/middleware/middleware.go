package middleware

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"golf-gamez/pkg/errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

// RequestIDKey is the context key for request IDs
type RequestIDKey string

const RequestIDCtxKey RequestIDKey = "request_id"

// RequestID middleware generates a unique request ID for each request
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDCtxKey, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger middleware logs HTTP requests with structured logging
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Get request ID from context
		requestID, _ := r.Context().Value(RequestIDCtxKey).(string)

		// Log request start
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", getRemoteAddr(r)).
			Str("user_agent", r.UserAgent()).
			Str("request_id", requestID).
			Msg("HTTP request started")

		// Process request
		next.ServeHTTP(ww, r)

		// Log request completion
		duration := time.Since(start)
		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.statusCode).
			Dur("duration", duration).
			Str("request_id", requestID).
			Msg("HTTP request completed")
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// getRemoteAddr extracts the real client IP address
func getRemoteAddr(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the chain
		if ip := net.ParseIP(xff); ip != nil {
			return ip.String()
		}
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip := net.ParseIP(xri); ip != nil {
			return ip.String()
		}
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// Recoverer middleware recovers from panics and logs them
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := r.Context().Value(RequestIDCtxKey).(string)

				log.Error().
					Interface("panic", err).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("request_id", requestID).
					Msg("HTTP request panic recovered")

				apiErr := errors.New(errors.ErrInternalServer, "An unexpected error occurred")
				apiErr.RequestID = requestID
				errors.WriteHTTPError(w, apiErr)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// RateLimiter holds rate limiting state
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	limit    rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		limit:    rate.Every(window / time.Duration(requests)),
		burst:    requests,
	}
}

// getLimiter returns the rate limiter for a given key (IP address)
func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		// Double-check pattern
		limiter, exists = rl.limiters[key]
		if !exists {
			limiter = rate.NewLimiter(rl.limit, rl.burst)
			rl.limiters[key] = limiter
		}
		rl.mu.Unlock()
	}

	return limiter
}

// RateLimit middleware implements per-IP rate limiting
func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(requests, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getRemoteAddr(r)
			rl := limiter.getLimiter(ip)

			if !rl.Allow() {
				apiErr := errors.NewWithDetails(
					errors.ErrRateLimitExceeded,
					"Rate limit exceeded",
					map[string]interface{}{
						"limit":  requests,
						"window": window.String(),
					},
				)
				errors.WriteHTTPError(w, apiErr)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GameCreationLimiter tracks game creation per IP
type GameCreationLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	limit    rate.Limit
	burst    int
}

// NewGameCreationLimiter creates a limiter specifically for game creation
func NewGameCreationLimiter(games int, window time.Duration) *GameCreationLimiter {
	return &GameCreationLimiter{
		limiters: make(map[string]*rate.Limiter),
		limit:    rate.Every(window / time.Duration(games)),
		burst:    games,
	}
}

// GameCreationLimit middleware limits game creation per IP
func GameCreationLimit(games int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewGameCreationLimiter(games, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only apply to POST requests on game creation endpoint
			if r.Method == http.MethodPost && r.URL.Path == "/v1/games" {
				ip := getRemoteAddr(r)

				limiter.mu.RLock()
				rl, exists := limiter.limiters[ip]
				limiter.mu.RUnlock()

				if !exists {
					limiter.mu.Lock()
					rl, exists = limiter.limiters[ip]
					if !exists {
						rl = rate.NewLimiter(limiter.limit, limiter.burst)
						limiter.limiters[ip] = rl
					}
					limiter.mu.Unlock()
				}

				if !rl.Allow() {
					apiErr := errors.NewWithDetails(
						errors.ErrRateLimitExceeded,
						"Game creation rate limit exceeded",
						map[string]interface{}{
							"limit":  games,
							"window": window.String(),
						},
					)
					errors.WriteHTTPError(w, apiErr)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORS headers middleware
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin {
					allowed = true
					break
				}
			}

			if allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
			w.Header().Set("Access-Control-Max-Age", "300")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDCtxKey).(string); ok {
		return requestID
	}
	return ""
}