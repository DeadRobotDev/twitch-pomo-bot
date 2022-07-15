package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/DeadRobotDev/twitch-pomo-bot/internal/bot"
	"github.com/DeadRobotDev/twitch-pomo-bot/internal/config"
)

const (
	configPath = "config.json"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	if _, err := os.Stat(configPath); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}

		if err := firstRunSetup(); err != nil {
			return err
		}
	}

	config, err := config.FromJSON(configPath)
	if err != nil {
		return err
	}

	b := bot.New(config)

	if err := b.Start(); err != nil {
		return err
	}

	return nil
}

func firstRunSetup() error {
	config := config.Default()

	for config.BotUsername == "" {
		fmt.Fprintf(os.Stdout, "Enter the bot's Twitch username: ")
		if _, err := fmt.Fscanln(os.Stdin, &config.BotUsername); err != nil {
			continue
		}
	}

	fmt.Fprintf(os.Stdout, "\nTo generate an OAuth Token, go to this link while logged into your bot account: https://twitchapps.com/tmi/\n\n")

	for config.BotAuthToken == "" {
		fmt.Fprintf(os.Stdout, "Enter the bot's OAuth Token, including the 'oauth:' part: ")
		if _, err := fmt.Fscanln(os.Stdin, &config.BotAuthToken); err != nil {
			continue
		}
	}

	fmt.Fprintf(os.Stdout, "\n")

	for config.ChannelName == "" {
		fmt.Fprintf(os.Stdout, "Enter the Twitch channel the bot should join: ")
		if _, err := fmt.Fscanln(os.Stdin, &config.ChannelName); err != nil {
			continue
		}
	}

	fmt.Fprintf(os.Stdout, "\nIf the details are correct, the bot should connect successfully. If not, check the 'config.json' file. If you're still having issues, check GitHub or reach out to me directly on Twitter (@DeadRobotDev), Twitch (DeadRobotDev), or Discord (Fletcher#9914).\n\n")

	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return nil
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}
