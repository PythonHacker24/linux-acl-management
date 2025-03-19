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
    // Loading Config file 
    backendConfig, err := utils.LoadConfig("./backend.yaml")
    if err != nil {
        slog.Error("Error loading config", "Error", err.Error())
    }

    // Connecting to File Server Daemons
    for _, server := range backendConfig.Servers {
        fmt.Printf("Loaded gRPC server: %s at %s \n", server.Name, server.Address)
    }

    for _, server := range backendConfig.Servers {
		go services.ConnectToServer(server)
	}

    // Setting up endpoints for interactions
    mux := http.NewServeMux()

    // Health Check Endpoint
    mux.Handle("/health", middleware.LoggingMiddleware(http.HandlerFunc(handlers.HealthHandler)))

    // Frontend Handlers 
    // For all the handlers for the frontend exposure, use the authentication middleware
    mux.Handle("/login", middleware.LoggingMiddleware(http.HandlerFunc(handlers.LoginHandler)))
    
    // list files handler -> GET METHOD

    // A FULL CONTENT HANDLER FOR ALL THE CRUD OPERATIONS IN HTTP METHODS
    // get file content handler     -> GET Method
    // delete file content handler  -> DEL Method
    // upload file content handler  -> POST METHOD
    // update file content handler  -> UPDATE METHOD

    slog.Info("Server Started Listening", "Host", config.Host, "Port", config.Port)
    if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host, config.Port), mux); err != nil {
        slog.Error(fmt.Sprintf("Failed to start server at port %s", config.Port), "Error", err.Error())
    }
}
