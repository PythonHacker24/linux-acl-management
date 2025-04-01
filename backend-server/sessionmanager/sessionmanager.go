package sessionmanager

import (
	"container/list"
	"fmt"
	"log/slog"
	"time"

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

func AddTransaction(username string, txn string) error {
    config.SessionMutex.Lock()
    defer config.SessionMutex.Unlock()

    // First check if the session exists
    session, exists := config.Sessions[username]
    if !exists {
        slog.Info("Session Not Found, Transaction Rejected", "User", username, "Transaction", txn)
        return fmt.Errorf("Session Not Found Transaction Rejected")
    }

    session.Mutex.Lock()
    defer session.Mutex.Unlock()

    session.TransactionQueue.PushBack(txn)

    if session.Timer != nil {
        session.Timer.Stop()
    }

    slog.Info("Transaction added successfully", "User", username, "Transaction", txn)
    return nil
}

// This function starts processing all the transactions in the queue one by one
func ProcessTransactions(username string, session *models.Session) {
    session.Mutex.Lock()
    defer session.Mutex.Unlock()

    counter := 0
    for session.TransactionQueue.Len() > 0 {
        element := session.TransactionQueue.Front()
        session.TransactionQueue.Remove(element)
        
        // TRANSACTION SHOULD BE PROCESSED HERE
        dummyFunc(counter)
        counter++
        
        slog.Info("Processed transaction", "User", username, "Element", element)
    }

    // At this point, all the transactions must be process in the queue
    if session.Timer != nil {
        session.Timer.Stop()
    }
    session.Timer = time.AfterFunc(config.SessionTimeout, func() { ExpireSession(username) })
    slog.Info("Session timeout restarted", "User", username)
}

func dummyFunc(counter int) {
    slog.Info(fmt.Sprintf("Transaction Started %d", counter))
    time.Sleep(2 * time.Second)
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
