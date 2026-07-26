package checker

import (
	"encoding/json"
	"os"
)

type Config struct {
    URLs           []string `json:"urls"`
    TimeoutSeconds int      `json:"timeout_seconds"`
    MaxConcurrent  int      `json:"max_concurrent"`
}

func LoadConfig(path string) (*Config, error){
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}