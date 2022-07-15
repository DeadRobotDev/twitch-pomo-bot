package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/DeadRobotDev/twitch-pomo-bot/internal/bot"
	"github.com/DeadRobotDev/twitch-pomo-bot/internal/config"
)

const (
	configPath        = "config.json"
	defaultConfigPath = "default.config.json"
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

		if err := copyFileContents(configPath, defaultConfigPath); err != nil {
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

func copyFileContents(dstPath, srcPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
