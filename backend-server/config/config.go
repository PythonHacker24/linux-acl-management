package config

import (
	"context"
	"sync"
	"time"
    "log/slog"

	"backend-server/models"
    "backend-server/utils"
    
    "github.com/redis/go-redis/v9"
)

var ( 
    Sessions = make(map[string]*models.Session)
    SessionMutex = sync.Mutex{}
    TransactionLogsRedisClient *redis.Client
    TransactionLogsRedisCtx = context.Background()
    BackendConfig *models.Config
)

const SessionTimeout = 5 * time.Minute

func InitTransactionLogsRedis(db int, address, password string) {
    TransactionLogsRedisClient = redis.NewClient(&redis.Options{
        Addr:       address, 
        Password:   password,
        DB:         db,
    })
}

func InitYamlConfig(file string) {
    var err error
    BackendConfig, err = utils.LoadConfig(file)
    if err != nil {
        slog.Error("Error loading config", "Error", err.Error())
    }
}
