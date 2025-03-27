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

type Mount struct {
    MountPoint string `yaml:"mount-point"`
}

type Operation struct {
    Read    bool    `yaml:"read"` 
    Write   bool    `yaml:"write"`
    Delete  bool    `yaml:"delete"`
    Update  bool    `yaml:"update"`
}

// Config represents the entire YAML structure
type Config struct {
    MountPath   []Mount         `yaml:"mount"` 
	Servers     []ServerConfig  `yaml:"servers"`
    Operations  []Operation     `yaml:"operations"`
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

// Request Payload Structure
type SetWorkingDirRequest struct {
	SetWorkingDir string `json:"setWorkingDir"`
}

// File Information Structure
type FileInfo struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	Permissions string    `json:"permissions"`
	User        string    `json:"user"`
	Group       string    `json:"group"`
	ModTime     time.Time `json:"mod_time"`
	IsDirectory bool      `json:"is_directory"`
}
