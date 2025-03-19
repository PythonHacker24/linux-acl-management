package utils

import (
	"os"

	"backend-server/models"

	"gopkg.in/yaml.v3"
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
