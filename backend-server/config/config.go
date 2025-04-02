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

const(
    Host = "localhost"
    Port = "9999"
)

var Sessions = make(map[string]*models.Session)
var SessionMutex = sync.Mutex{}

const SessionTimeout = 5 * time.Minute

var TransactionLogsRedisClient *redis.Client
var TransactionLogsRedisCtx = context.Background()

var BackendConfig *models.Config

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
