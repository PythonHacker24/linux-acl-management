package routes

import (
	"net/http"

	"backend-server/middleware"
	"backend-server/handlers"
)

func SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

    mux.Handle("/health", middleware.LoggingMiddleware(http.HandlerFunc(handlers.HealthHandler)))

    mux.Handle("/login", middleware.LoggingMiddleware(http.HandlerFunc(handlers.LoginHandler)))
    mux.Handle("POST /issue-transaction", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.TransactionHandler)))
    mux.Handle("GET /transaction-result", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.GetTransactionResult)))
    
    mux.Handle("GET /current-working-directory", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.GetCurrentWorkingDir)))
    mux.Handle("POST /current-working-directory", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.SetCurrentWorkingDir)))

    mux.Handle("GET /list-files", middleware.LoggingMiddleware(middleware.AuthenticationMiddleware(handlers.ListFilesInDir)))

    // Permission Management Endpoints

    // User Settings APIs (For the future) 

	return mux
}

