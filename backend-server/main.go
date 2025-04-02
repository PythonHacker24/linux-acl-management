package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"backend-server/config"
	"backend-server/handlers"
	"backend-server/middleware"
	"backend-server/services"
	"backend-server/sessionmanager"
)

func main() {

    config.InitYamlConfig("./backend.yaml")

    backendConfig := config.BackendConfig

    config.InitTransactionLogsRedis(backendConfig.TransactionLogRedis[0].DB, backendConfig.TransactionLogRedis[0].Address, backendConfig.TransactionLogRedis[0].Password)

    fmt.Printf("Base Path set to: %s \n", backendConfig.BasePath)

    // Connecting to File Server Daemons
    for _, server := range backendConfig.Servers {
        fmt.Printf("Loaded gRPC server: Method %s at %s mounted \n", server.Method, server.Path)
    }

    for _, server := range backendConfig.Servers {
        /*
            ONLY FOR PROTOTYPING! 
            USE SAFE AND CONFIGURABLE GOROUTINES IN THE FINAL VERSION
        */
		go services.ConnectToServer(server)
	}

	// Start Transaction Worker in background
	go sessionmanager.TransactionWorker()

    // Setting up endpoints for interactions
    mux := http.NewServeMux()

    // Health Check Endpoint
    mux.Handle("/health", middleware.LoggingMiddleware(http.HandlerFunc(handlers.HealthHandler)))

    // Frontend Handlers 
    mux.Handle("/login", middleware.LoggingMiddleware(http.HandlerFunc(handlers.LoginHandler)))
    mux.Handle("POST /issue-transaction", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.TransactionHandler)))
    mux.Handle("GET /transaction-result", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.GetTransactionResult)))
    
    // list files handler -> GET METHOD

    mux.Handle("GET /current-working-directory", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.GetCurrentWorkingDir)))
    mux.Handle("POST /current-working-directory", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.SetCurrentWorkingDir)))

    mux.Handle("GET /list-files", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.ListFilesInDir)))

    // Permission Management Endpoints

    // User Settings APIs (For the future) 

    slog.Info("Server Started Listening", "Host", config.BackendConfig.DeploymentConfig[0].Host, "Port", strconv.Itoa(config.BackendConfig.DeploymentConfig[0].Port))
    if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.BackendConfig.DeploymentConfig[0].Host, strconv.Itoa(config.BackendConfig.DeploymentConfig[0].Port)), mux); err != nil {
        slog.Error(fmt.Sprintf("Failed to start server at port %s", strconv.Itoa(config.BackendConfig.DeploymentConfig[0].Port)), "Error", err.Error())
    }
}
