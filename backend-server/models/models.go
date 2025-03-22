package models

import (
	"container/list"
	"sync"
	"time"
)

// Health Response Structure
type HealthResponse struct {
    Status string `json:"status"`
}

// ServerConfig represents a single gRPC server in the YAML
type ServerConfig struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

// Config represents the entire YAML structure
type Config struct {
	Servers []ServerConfig `yaml:"servers"`
}

// User struct for login request
type User struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// Session Struct for Session Manager
type Session struct {
    Username            string 
    CurrentWorkingDir   string
    Expiry              time.Time
    Timer               *time.Timer
    TransactionQueue    *list.List
    Mutex               sync.Mutex
}

// Transaction request body
type TransactionRequest struct {
	Transaction string `json:"transaction"`
}
