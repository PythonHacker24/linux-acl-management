package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        slog.Info("Received request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
        next.ServeHTTP(w, r)
        slog.Info("Request completed", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.Duration("duration", time.Since(start)))
    })
} 
