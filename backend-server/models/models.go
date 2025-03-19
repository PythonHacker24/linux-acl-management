package models

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
