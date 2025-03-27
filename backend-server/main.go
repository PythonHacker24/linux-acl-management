package main

import (
	"fmt"
	"log/slog"
	"net/http"

	"backend-server/config"
	"backend-server/handlers"
	"backend-server/middleware"
	"backend-server/services"
	"backend-server/sessionmanager"
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
        fmt.Printf("Loaded gRPC server: %s at %s mounted \n", server.Name, server.Address)
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
    
    // list files handler -> GET METHOD

    mux.Handle("GET /current-working-directory", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.GetCurrentWorkingDir)))
    mux.Handle("POST /current-working-directory", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.SetCurrentWorkingDir)))

    mux.Handle("GET /list-files", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.ListFilesInDir)))

    // This works on files only and not directories. These hanlders allow you to get download the file, delete the file from the servers, upload a new file or update the file. These files must be in the current directory - which is tracked in the session itself.

    /* NOT THE PART OF PROTOTYPE - WOULD BE DONE AT THE END */ 
    if backendConfig.Operations[0].Read {
        mux.Handle("GET /file", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.GetFile)))
    }

    if backendConfig.Operations[0].Delete {
        mux.Handle("DELETE /file", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.DeleteFile)))
    }

    if backendConfig.Operations[0].Write {
        mux.Handle("POST /file", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.UploadFile)))
    }

    if backendConfig.Operations[0].Update {
        mux.Handle("UPDATE /file", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.UpdateFile)))
    }

    // Permission Management Endpoints
    
    // User Settings APIs (For the future) 

    slog.Info("Server Started Listening", "Host", config.Host, "Port", config.Port)
    if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Host, config.Port), mux); err != nil {
        slog.Error(fmt.Sprintf("Failed to start server at port %s", config.Port), "Error", err.Error())
    }
}
