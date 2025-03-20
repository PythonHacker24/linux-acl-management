package sessionmanager

import (
	"backend-server/models"
	"container/list"
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

func AddTransaction(username string, txn string) {
    sessionMutex.Lock()
    defer sessionMutex.Unlock()

    // First check if the session exists
    session, exists := sessions[username]
    if !exists {
        slog.Info("Session Not Found, Transaction Rejected", "User", username, "Transaction", txn)
        return
    }

    session.Mutex.Lock()
    defer session.Mutex.Unlock()

    err := session.TransactionQueue.PushBack(txn)
    if err != nil {
        slog.Info("Failed to add transaction to the queue", "User", username, "Transaction", txn)
        return
    }

    if session.Timer != nil {
        session.Timer.Stop()
    }

    slog.Info("Transaction added successfully", "User", username, "Transaction", txn)
}

// This function starts processing all the transactions in the queue one by one
func ProcessTransactions(username string) {
    sessionMutex.Lock()
    defer sessionMutex.Unlock()

    session, exists := sessions[username]
    if !exists {
        return 
    }
    
    session.Mutex.Lock()
    defer session.Mutex.Unlock()

    for session.TransactionQueue.Len() > 0 {
        element := session.TransactionQueue.Front()
        session.TransactionQueue.Remove(element)
        
        // TRANSACTION SHOULD BE PROCESSED HERE
        
        slog.Info("Processed transaction", "User", username, "Element", element)
    }

    // At this point, all the transactions must be process in the queue

    session.Timer = time.AfterFunc(sesstionTimeout, func() { ExpireSession(username) })
    slog.Info("Session timeout restarted", "User", username)
}
