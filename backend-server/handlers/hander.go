package handlers

import (
	"backend-server/authentication"
	"backend-server/config"
	"backend-server/ldap"
	"backend-server/models"
	"encoding/json"
	"log/slog"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    var response models.HealthResponse
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    response.Status = "ok"
    if err := json.NewEncoder(w).Encode(response); err != nil {
        slog.Error("Failed to send health response from the handler", "Error", err.Error())
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return 
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate input
    if user.Username == "" || user.Password == "" {
        http.Error(w, "Username and password are required", http.StatusBadRequest)
        return
    }

    authStatus := ldap.AuthenticateUser(user.Username, user.Password, config.LdapSearchBase)
    if !authStatus {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    token, err := authentication.GenerateJWT(user.Username)
    if err != nil {
        slog.Error("Error generating token", "error", err.Error())
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    response := map[string]string{"token": token}
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        slog.Error("Failed to encode response", "error", err.Error())
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}
