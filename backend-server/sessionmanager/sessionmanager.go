package sessionmanager

import (
	"container/list"
	"fmt"
	"log/slog"
	"time"
    "strings"

	"backend-server/models"
	"backend-server/config"
    "backend-server/utils"
)

func CreateSession(username string) {

    backendConfig, err := utils.LoadConfig("./backend.yaml")
    if err != nil {
        slog.Error("Error loading config", "Error", err.Error())
    }

    config.SessionMutex.Lock()
    defer config.SessionMutex.Unlock()

    // Check if the session already exists
    if _, exists := config.Sessions[username]; exists {
        return 
    }

    session := &models.Session{
        Username:           username,
        Expiry:             time.Now().Add(config.SessionTimeout),
        Timer:              time.AfterFunc(config.SessionTimeout, func() { ExpireSession(username) }),
        TransactionQueue:   list.New(),
        CurrentWorkingDir:  backendConfig.BasePath,
    }

    config.Sessions[username] = session 
}

func ExpireSession(username string) {
    config.SessionMutex.Lock()
    defer config.SessionMutex.Unlock()

    // Delete session if it exists
    if _, exists := config.Sessions[username]; exists {
        slog.Info("Session Expired", "User", username)
        delete(config.Sessions, username)
    }
}

func AddTransaction(username, txn string) (string, error) {
	config.SessionMutex.Lock()
	session, exists := config.Sessions[username]
	config.SessionMutex.Unlock()

	if !exists {
		return "", fmt.Errorf("session not found")
	}

	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	txnID := utils.GenerateTxnID() // Generate unique transaction ID
	session.TransactionQueue.PushBack(txnID + ":" + txn) 

	// Store txnID and username in Redis
	err := config.TransactionLogsRedisClient.Set(config.TransactionLogsRedisCtx, txnID, "PENDING", 0).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store txnID in Redis")
	}

	err = config.TransactionLogsRedisClient.Set(config.TransactionLogsRedisCtx, "user:"+txnID, username, 0).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store txnID ownership in Redis")
	}

	if session.Timer != nil {
		session.Timer.Stop()
	}

	return txnID, nil
}

// This function starts processing all the transactions in the queue one by one
func ProcessTransactions(username string, session *models.Session) {
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	for session.TransactionQueue.Len() > 0 {
		element := session.TransactionQueue.Front()
		session.TransactionQueue.Remove(element)

		txnData := element.Value.(string)
		parts := strings.SplitN(txnData, ":", 2) // txnID:txn
		txnID, txn := parts[0], parts[1]

		result := dummyFunc(txn)

		// Store result in Redis
		err := config.TransactionLogsRedisClient.Set(config.TransactionLogsRedisCtx, txnID, result, 0).Err()
		if err != nil {
			slog.Error("Failed to store transaction result in Redis", "txnID", txnID)
		}

		slog.Info("Processed transaction", "User", username, "txnID", txnID, "Result", result)
	}

	if session.Timer != nil {
		session.Timer.Stop()
	}
	session.Timer = time.AfterFunc(config.SessionTimeout, func() { ExpireSession(username) })
}

func dummyFunc(txn string) string {
    slog.Info(fmt.Sprintf("Transaction Started: Txn %s", txn))
    time.Sleep(2 * time.Second)
    return fmt.Sprintf("Executed: Txn: %s", txn)
}

func TransactionWorker() {
    slog.Info("Transaction worker started")
    for {
        config.SessionMutex.Lock()
        // Create a local copy of config.Sessions to process
        SessionsToProcess := make(map[string]*models.Session)
        for username, session := range config.Sessions {
            if session.TransactionQueue.Len() > 0 {
                SessionsToProcess[username] = session
            }
        }

        config.SessionMutex.Unlock()

        // Process transactions for each session without holding the global lock
        for username, session := range SessionsToProcess {
            ProcessTransactions(username, session)
        }

        slog.Info("Transaction worker finished all the transactions")
        time.Sleep(2 * time.Second)
    }
}
