package http

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/cors"
)

// JSONer is an interface for writing JSON responses.
type JSONer interface {
	JSON(w http.ResponseWriter, r *http.Request, code int, data interface{})
}

// Logger is an interface for logging.
type Logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

// RequestLogger logs requests.
func RequestLogger(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{w, http.StatusOK}
			defer func() {
				logger.Info("request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.status,
					"duration", time.Since(start),
				)
			}()
			next.ServeHTTP(ww, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// PanicRecoverer recovers from panics.
func PanicRecoverer(logger Logger, renderer JSONer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Error("panic recovered", "error", err, "stack", string(debug.Stack()))
					renderer.JSON(w, r, http.StatusInternalServerError, map[string]string{
						"error":   "internal_error",
						"message": "An unexpected error occurred",
					})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Cors sets up CORS.
func Cors(allowedOrigins []string) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	return c.Handler
}
