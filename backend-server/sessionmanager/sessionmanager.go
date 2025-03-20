package sessionmanager

import (
	"backend-server/models"
	"container/list"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var sessions = make(map[string]*models.Session)
var sessionMutex = sync.Mutex{}

const sesstionTimeout = 5 * time.Minute

func CreateSession(username string) {
    sessionMutex.Lock()
    defer sessionMutex.Unlock()

    // Check if the session already exists
    if _, exists := sessions[username]; exists {
        return 
    }

    session := &models.Session{
        Username:           username,
        Expiry:             time.Now().Add(sesstionTimeout),
        Timer:              time.AfterFunc(sesstionTimeout, func() { ExpireSession(username) }),
        TransactionQueue:   list.New(),
    }

    sessions[username] = session 
}

func ExpireSession(username string) {
    sessionMutex.Lock()
    defer sessionMutex.Unlock()

    // Delete session if it exists
    if _, exists := sessions[username]; exists {
        slog.Info("Session Expired", "User", username)
        delete(sessions, username)
    }
}

func AddTransaction(username string, txn string) error {
    sessionMutex.Lock()
    defer sessionMutex.Unlock()

    // First check if the session exists
    session, exists := sessions[username]
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
    session.Timer = time.AfterFunc(sesstionTimeout, func() { ExpireSession(username) })
    slog.Info("Session timeout restarted", "User", username)
}

func dummyFunc(counter int) {
    slog.Info(fmt.Sprintf("Transaction Started %d", counter))
    time.Sleep(2 * time.Second)
}

func TransactionWorker() {
    slog.Info("Transaction worker started")
    for {
        sessionMutex.Lock()
        // Create a local copy of sessions to process
        sessionsToProcess := make(map[string]*models.Session)
        for username, session := range sessions {
            if session.TransactionQueue.Len() > 0 {
                sessionsToProcess[username] = session
            }
        }
        sessionMutex.Unlock()

        // Process transactions for each session without holding the global lock
        for username, session := range sessionsToProcess {
            ProcessTransactions(username, session)
        }

        slog.Info("Transaction worker finished all the transactions")
        time.Sleep(2 * time.Second) // Adjust based on processing needs
    }
}
