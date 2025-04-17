package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	"backend-server/config"
	"backend-server/routes"
	"backend-server/services"
	"backend-server/sessionmanager"
)

func init() {
    config.InitYamlConfig("./backend.yaml")

    backendConfig := config.BackendConfig
    basePath := backendConfig.BasePath

    slog.Info("Base Path set to", "BasePath", basePath)
    if basePath == "" || basePath == "/" {
        slog.Info("Base Path not set, considering mount paths as absolute paths")
        backendConfig.BasePath = ""
    }
	if len(backendConfig.TransactionLogRedis) == 0 {
		slog.Error("No Redis config found")
		return
	}

    config.InitTransactionLogsRedis(backendConfig.TransactionLogRedis[0].DB, backendConfig.TransactionLogRedis[0].Address, backendConfig.TransactionLogRedis[0].Password)

    for _, server := range backendConfig.Servers {
        /*
            ONLY FOR PROTOTYPING! 
            USE SAFE AND CONFIGURABLE GOROUTINES IN THE FINAL VERSION
        */
        slog.Info("Loaded gRPC server", "Method", server.Method, "Mount", server.Path)
		go services.ConnectToServer(server)
	}
}

func main() {

	officialCPUCount := runtime.NumCPU()
	maxCPUs := config.BackendConfig.DeploymentConfig[0].MaxTransactionWorkers // MODIFY THIS TO CONFIG FILE
	
	numWorkers := maxCPUs
	if maxCPUs == 0 {
		numWorkers = 1
	} else if maxCPUs < officialCPUCount {
		numWorkers = maxCPUs
	}

	slog.Info("Spawning GoRoutines: ", numWorkers)

	var TransactionWorkersWg sync.WaitGroup
	TransactionWorkersWg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer TransactionWorkersWg.Done()
			sessionmanager.TransactionWorker(workerID)
		}(i)
	}

    mux := routes.SetupRoutes() 

    slog.Info("Server Started Listening", "Host", config.BackendConfig.DeploymentConfig[0].Host, "Port", strconv.Itoa(config.BackendConfig.DeploymentConfig[0].Port))
    if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.BackendConfig.DeploymentConfig[0].Host, strconv.Itoa(config.BackendConfig.DeploymentConfig[0].Port)), mux); err != nil {
        slog.Error(fmt.Sprintf("Failed to start server at port %s", strconv.Itoa(config.BackendConfig.DeploymentConfig[0].Port)), "Error", err.Error())
    }

	TransactionWorkersWg.Wait()
}
