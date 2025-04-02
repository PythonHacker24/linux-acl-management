package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
    "os/user"
	"syscall"

	"backend-server/authentication"
	"backend-server/config"
	"backend-server/ldap"
	"backend-server/models"
	"backend-server/sessionmanager"
    "backend-server/utils"

    "github.com/redis/go-redis/v9"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
    var response models.HealthResponse
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    response.Status = "ok"
    if err := json.NewEncoder(w).Encode(response); err != nil {
        slog.Error("Failed to send health response from the handler", "Error", err.Error())
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return 
    }

    var user models.User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Validate input
    if user.Username == "" || user.Password == "" {
        http.Error(w, "Username and password are required", http.StatusBadRequest)
        return
    }

    authStatus := ldap.AuthenticateUser(user.Username, user.Password, config.BackendConfig.LdapConfig[0].LdapSearchBase)
    if !authStatus {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    sessionmanager.CreateSession(user.Username)

    token, err := authentication.GenerateJWT(user.Username)
    if err != nil {
        slog.Error("Error generating token", "error", err.Error())
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    response := map[string]string{"token": token}
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(response); err != nil {
        slog.Error("Failed to encode response", "error", err.Error())
        http.Error(w, "Failed to encode response", http.StatusInternalServerError)
        return
    }
}

func TransactionHandler(w http.ResponseWriter, r *http.Request) {
	username, err := authentication.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT Token", http.StatusUnauthorized)
		return
	}

	var txnReq models.TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&txnReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

    var txnID string
	txnID, err = sessionmanager.AddTransaction(username, txnReq.Transaction)
    if err != nil {
        fmt.Fprintf(w, "Failed to add transaction")
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, fmt.Sprintf("Transaction added: %s", txnID))
}

func GetCurrentWorkingDir(w http.ResponseWriter, r *http.Request) {
	username, err := authentication.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT Token", http.StatusUnauthorized)
		return
	}

    config.Sessions[username].Mutex.Lock()
    defer config.Sessions[username].Mutex.Unlock()

    currentDir := config.Sessions[username].CurrentWorkingDir

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{'currentDir': %s}", currentDir)
}

func SetCurrentWorkingDir(w http.ResponseWriter, r *http.Request) {
	username, err := authentication.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT Token", http.StatusUnauthorized)
		return
	}
    
    var req models.SetWorkingDirRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON Request", http.StatusBadRequest)
        return
    }

    if req.SetWorkingDir == "" {
        http.Error(w, "Directory name cannot be empty", http.StatusBadRequest)
        return
    }

    if string(req.SetWorkingDir[len(req.SetWorkingDir) - 1]) != "/" {
        http.Error(w, "Need / at the end of the directory name", http.StatusBadRequest)
        return
    }

    config.Sessions[username].Mutex.Lock()
    defer config.Sessions[username].Mutex.Unlock()

    backendConfig, err := utils.LoadConfig("./backend.yaml")
    if err != nil {
        slog.Error("Error loading config", "Error", err.Error())
    }

    // currentDir := config.Sessions[username].CurrentWorkingDir
	newDir := filepath.Join(backendConfig.BasePath, req.SetWorkingDir)
	config.Sessions[username].CurrentWorkingDir = newDir

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "{'currentDir': '%s'}", newDir)
}

func ListFilesInDir(w http.ResponseWriter, r *http.Request) {
	username, err := authentication.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Failed to get username from JWT Token", http.StatusUnauthorized)
		return
	}

    files, err := os.ReadDir(config.Sessions[username].CurrentWorkingDir)
	if err != nil {
		http.Error(w, "Failed to read directory", http.StatusInternalServerError)
		return
	}

	var fileList []models.FileInfo

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		// Get file owner and group
		sys := info.Sys().(*syscall.Stat_t)
		uid := fmt.Sprint(sys.Uid)
		gid := fmt.Sprint(sys.Gid)

		// Convert UID and GID to user and group names
		usr, err := user.LookupId(uid)
		if err != nil {
			usr = &user.User{Username: uid} // Fallback to UID
		}

		grp, err := user.LookupGroupId(gid)
		if err != nil {
			grp = &user.Group{Name: gid} // Fallback to GID
		}

		// Append file info to list
		fileList = append(fileList, models.FileInfo{
			Name:        file.Name(),
			Size:        info.Size(),
			Permissions: info.Mode().String(),
			User:        usr.Username,
			Group:       grp.Name,
			ModTime:     info.ModTime(),
			IsDirectory: file.IsDir(),
		})
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode response as JSON
	if err := json.NewEncoder(w).Encode(fileList); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func GetTransactionResult(w http.ResponseWriter, r *http.Request) {
	username, err := authentication.GetUsernameFromJWT(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	txnID := r.URL.Query().Get("txnID")
	if txnID == "" {
		http.Error(w, "Missing txnID", http.StatusBadRequest)
		return
	}

	storedUsername, err := config.TransactionLogsRedisClient.Get(config.TransactionLogsRedisCtx, "user:"+txnID).Result()
	if err == redis.Nil {
		http.Error(w, "Transaction ID not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching transaction ownership", http.StatusInternalServerError)
		return
	}

	if storedUsername != username {
        http.Error(w, "Forbidden: You don't own the transaction", http.StatusForbidden)
		return
	}

	result, err := config.TransactionLogsRedisClient.Get(config.TransactionLogsRedisCtx, txnID).Result()
	if err == redis.Nil {
		http.Error(w, "Transaction ID not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching transaction result", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"txnID": txnID, "result": result})
}
