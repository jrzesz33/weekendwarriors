package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Development server for serving the Golf Gamez PWA
func main() {
	var (
		port = flag.String("port", "8000", "Port to serve on")
		dir  = flag.String("dir", "build", "Directory to serve")
		cors = flag.Bool("cors", true, "Enable CORS headers")
	)
	flag.Parse()

	// Ensure the directory exists
	if _, err := os.Stat(*dir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist. Run 'make build' first.", *dir)
	}

	// Create file server
	fs := http.FileServer(http.Dir(*dir))

	// Wrap with middleware
	handler := loggingMiddleware(corsMiddleware(spaMiddleware(*dir, fs), *cors))

	// Start server
	addr := ":" + *port
	fmt.Printf("🚀 Golf Gamez Frontend Server\n")
	fmt.Printf("📂 Serving: %s\n", *dir)
	fmt.Printf("🌐 Address: http://localhost%s\n", addr)
	fmt.Printf("📱 Mobile: http://[your-ip]%s\n", addr)
	fmt.Printf("\n💡 Tips:\n")
	fmt.Printf("   - Make sure the backend API is running on :8080\n")
	fmt.Printf("   - Open in mobile browser for best PWA experience\n")
	fmt.Printf("   - Use 'make dev' for auto-rebuild during development\n")
	fmt.Printf("\n🛑 Press Ctrl+C to stop\n\n")

	log.Fatal(http.ListenAndServe(addr, handler))
}

// spaMiddleware handles Single Page Application routing
func spaMiddleware(dir string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the path
		path := r.URL.Path

		// Check if it's an API request (proxy to backend)
		if strings.HasPrefix(path, "/v1/") || strings.HasPrefix(path, "/api/") {
			proxyToBackend(w, r)
			return
		}

		// Check if the file exists
		fullPath := filepath.Join(dir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// If it doesn't exist and it's not a static asset, serve index.html
			if !strings.Contains(path, ".") {
				// This is likely a SPA route, serve the main app
				r.URL.Path = "/"
			}
		}

		// Set proper MIME types for WASM and other assets
		if strings.HasSuffix(path, ".wasm") {
			w.Header().Set("Content-Type", "application/wasm")
		} else if strings.HasSuffix(path, ".js") {
			w.Header().Set("Content-Type", "text/javascript")
		} else if strings.HasSuffix(path, ".css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, ".json") {
			w.Header().Set("Content-Type", "application/json")
		} else if strings.HasSuffix(path, ".png") {
			w.Header().Set("Content-Type", "image/png")
		}

		// Serve the file
		h.ServeHTTP(w, r)
	})
}

// corsMiddleware adds CORS headers for development
func corsMiddleware(h http.Handler, enabled bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enabled {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs all requests
func loggingMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Serve the request
		h.ServeHTTP(wrapper, r)

		// Log the request
		fmt.Printf("%s %s %d %s\n",
			r.Method,
			r.URL.Path,
			wrapper.statusCode,
			r.RemoteAddr,
		)
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

// proxyToBackend proxies API requests to the backend server
func proxyToBackend(w http.ResponseWriter, r *http.Request) {
	// This is a simple proxy for development
	// In production, you'd use a proper reverse proxy like nginx

	backendURL := "http://localhost:8080" + r.URL.Path
	if r.URL.RawQuery != "" {
		backendURL += "?" + r.URL.RawQuery
	}

	// Create new request
	req, err := http.NewRequest(r.Method, backendURL, r.Body)
	if err != nil {
		http.Error(w, "Failed to create proxy request", http.StatusInternalServerError)
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Backend server unavailable", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Copy status code
	w.WriteHeader(resp.StatusCode)

	// Copy body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response", http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
