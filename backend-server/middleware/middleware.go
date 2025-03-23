package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

    "backend-server/authentication"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        slog.Info("Received request", slog.String("method", r.Method), slog.String("path", r.URL.Path))
        next.ServeHTTP(w, r)
        slog.Info("Request completed", slog.String("method", r.Method), slog.String("path", r.URL.Path), slog.Duration("duration", time.Since(start)))
    })
} 

func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Missing or Invalid Token", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        claims, err := authentication.ValidateJWT(tokenString)
        if err != nil {
            http.Error(w, "Invalid Token", http.StatusUnauthorized)
            return
        }

        username, ok := claims["username"].(string)
        if !ok {
            http.Error(w, "Invalid Token Payload", http.StatusUnauthorized)
            return
        }

        r.Header.Set("X-User", username)

        next(w, r)
    }
}
