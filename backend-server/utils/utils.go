package utils

import (
	"os"

	"backend-server/models"

	"gopkg.in/yaml.v3"
    "github.com/google/uuid"
)

func LoadConfig(filename string) (*models.Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    var config models.Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        return nil, err
    }

    return &config, nil
}

func GenerateTxnID() string {
	return uuid.New().String()
}
