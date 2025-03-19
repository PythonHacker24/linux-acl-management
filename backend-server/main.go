package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"backend-server/config"
	"backend-server/handlers"
	"backend-server/middleware"
	"backend-server/services"
	"backend-server/utils"
)

func main() {
    backendConfig, err := utils.LoadConfig("./backend.yaml")
    if err != nil {
        slog.Error("Error loading config", "Error", err.Error())
    }

    for _, server := range backendConfig.Servers {
        fmt.Printf("Loaded gRPC server: %s at %s \n", server.Name, server.Address)
    }

    for _, server := range backendConfig.Servers {
		go services.ConnectToServer(server)
	}

    mux := http.NewServeMux()

    // Health Check Endpoint
    mux.Handle("/health", middleware.LoggingMiddleware(http.HandlerFunc(handlers.HealthHandler)))

    // Frontend Handlers 

    // Daemons Handlers

    slog.Info("Server Started Listening", "Host", config.Host, "Port", config.Port)
    if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host, config.Port), mux); err != nil {
        slog.Error("Failed to start server at port 8080", "Error", err.Error())
    }
}
