package handlers

import (
    "net/http"
    "log/slog"
    "encoding/json"
    "backend-server/models"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    var response models.HealthResponse
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    response.Status = "ok"
    if err := json.NewEncoder(w).Encode(response); err != nil {
        slog.Error("Failed to send health response from the handler", err)
    }
}
