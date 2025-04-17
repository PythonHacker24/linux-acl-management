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

type Redis struct {
    Address     string  `yaml:"address"`   
    Password    string  `yaml:"password"`
    DB          int     `yaml:"db"`  
}

type LdapConfig struct {    
    LdapServer      string  `yaml:"ldap-server"`
    LdapBindDN      string  `yaml:"ldap-bind-dn"`
    LdapPassword    string  `yaml:"ldap-password"`
    LdapSearchBase  string  `yaml:"ldap-search-base"` 
}

type DeploymentConfig struct {
    Host    				string      `yaml:"host"` 
    Port    				int         `yaml:"port"`
	MaxTransactionWorkers  	int			`yaml:"max-transaction-workers"`
}

// Contains complete configurations from YAML file
type Config struct {
    DeploymentConfig    []DeploymentConfig  `yaml:"deployment-config"`
    TransactionLogRedis []Redis             `yaml:"transaction-logs-redis"`
    LdapConfig          []LdapConfig        `yaml:"ldap-config"`
	BasePath            string              `yaml:"basePath"`
	Servers             []Server            `yaml:"servers"`
}

// Contains server configurations from YAML file
type Server struct {
	Path   string  `yaml:"path"`
	Method string  `yaml:"method"`
	Remote *Remote `yaml:"remote,omitempty"`
}

// Contains remote configurations from YAML file
type Remote struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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
