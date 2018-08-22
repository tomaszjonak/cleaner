package cleaner

import (
	"encoding/json"
	"io"
	"os"
)

type Configuration struct {
	RootDir              string `json:"root_directory"`
	DefaultRetentionDays int    `json:"default_retention_days"`
}

func MakeConfigurationFromFile(configPath string) (*Configuration, error) {
	reader, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	return MakeConfigurationFromReader(reader)
}

func MakeConfigurationFromReader(reader io.Reader) (*Configuration, error) {
	config := &Configuration{}
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
