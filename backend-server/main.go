package main

import (
	"fmt"
	"log/slog"
	"net/http"
    "os"

    "backend-server/config"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Ok")
}

func main() {
    opts := &slog.HandlerOptions{Level: slog.LevelWarn}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

    mux := http.NewServeMux()

    // Frontend Handlers 
    mux.HandleFunc("/health", healthHandler)

    // Backend Handlers

    logger.Info("Server Started Listening", "Host", config.Host, "Port", config.Port)
    if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host, config.Port), mux); err != nil {
        slog.Error("Failed to start server at port 8080", err)
    }
}
