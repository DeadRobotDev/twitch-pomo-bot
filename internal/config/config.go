package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BotUsername  string
	BotAuthToken string
	ChannelName  string

	CommandPrefix string

	TaskHelpMessage       string
	TaskInProgressMessage string
	NoTaskMessage         string
	TaskAddedMessage      string
	TaskEditedMessage     string
	TaskCompletedMessage  string
	TaskCancelledMessage  string
}

func FromJSON(path string) (*Config, error) {
	var config Config

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
