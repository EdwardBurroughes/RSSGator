package config

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

const gatorConfigFilePathName string = ".gatorconfig.json"

type Config struct {
	DBUrl           string  `json:"db_url"`
	CurrentUserName *string `json:"current_user_name"`
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = &user
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	return writeToFile(configPath, cfg)
}

func Read() (Config, error) {
	// Ok to load full file in memory as small
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, nil
	}
	return readConfigFile(configPath)

}

func writeToFile(filePath string, data *Config) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func readConfigFile(filePath string) (Config, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return Config{}, nil
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, nil
	}

	var config Config
	json.Unmarshal(byteValue, &config)
	return config, nil

}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homeDir, gatorConfigFilePathName), nil
}
