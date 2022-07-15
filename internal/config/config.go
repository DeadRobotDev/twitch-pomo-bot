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

func Default() *Config {
	return &Config{
		CommandPrefix:         "!",
		TaskHelpMessage:       "Usage: To add or edit a task, %COMMAND_PREFIX%task <add | edit> <task name>. To complete or cancel a task, %COMMAND_PREFIX%task <done | delete>.",
		TaskInProgressMessage: "You already have a task in progress. To edit the current task, %COMMAND_PREFIX%task edit <new task name>. To complete or cancel the current task, %COMMAND_PREFIX%task <done | delete>.",
		NoTaskMessage:         "You don't have a current task. To start a task, %COMMAND_PREFIX%task add <task name>.",
		TaskAddedMessage:      "Starting work on %TASK_NAME%. Good luck!",
		TaskEditedMessage:     "You have updated your current task.",
		TaskCompletedMessage:  "Congrats! You have completed %TASK_NAME%.",
		TaskCancelledMessage:  "You have cancelled your current task.",
	}
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
