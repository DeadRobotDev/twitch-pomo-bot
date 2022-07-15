package bot

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/DeadRobotDev/twitch-pomo-bot/internal/config"
	"github.com/gempir/go-twitch-irc/v3"
)

type Bot struct {
	config *config.Config
	client *twitch.Client

	tasksMtx sync.RWMutex
	tasks    map[string]*Task
}

const tasksPath = "tasks.txt"

func New(config *config.Config) *Bot {
	return &Bot{
		config: config,
		client: twitch.NewClient(config.BotUsername, config.BotAuthToken),
		tasks:  make(map[string]*Task),
	}
}

func (b *Bot) Start() error {
	b.UpdateTasksFile()

	b.client.OnConnect(b.onConnect)
	b.client.OnReconnectMessage(b.onReconnect)
	b.client.OnPrivateMessage(b.onPrivateMessage)

	return b.client.Connect()
}

func (b *Bot) Reply(parentMessageID, text string, replacer *strings.Replacer) {
	b.client.Reply(b.config.ChannelName, parentMessageID, replacer.Replace(text))
}

func (b *Bot) UpdateTasksFile() {
	file, err := os.Create(tasksPath)
	if err != nil {
		return
	}
	defer file.Close()

	b.tasksMtx.RLock()
	defer b.tasksMtx.RUnlock()

	for _, task := range b.tasks {
		file.WriteString(task.String())
	}
}

func (b *Bot) AddTask(user twitch.User, taskName string) bool {
	b.tasksMtx.Lock()
	defer b.tasksMtx.Unlock()

	if _, ok := b.tasks[user.ID]; ok {
		return false
	}

	b.tasks[user.ID] = newTask(user, taskName)
	return true
}

func (b *Bot) EditTask(userID, taskName string) bool {
	b.tasksMtx.Lock()
	defer b.tasksMtx.Unlock()

	task, ok := b.tasks[userID]
	if !ok {
		return false
	}

	task.Name = taskName
	return true
}

func (b *Bot) RemoveTask(userID string) (*Task, bool) {
	b.tasksMtx.Lock()
	defer b.tasksMtx.Unlock()

	task, ok := b.tasks[userID]
	if !ok {
		return nil, false
	}

	delete(b.tasks, userID)

	return task, true
}

func (b *Bot) onConnect() {
	log.Printf("CONNECTED\n")

	b.client.Join(b.config.ChannelName)
	b.client.Say(b.config.ChannelName, "I have arrived!")
}

func (b *Bot) onReconnect(message twitch.ReconnectMessage) {
	log.Printf("RECONNECTED\n")
}

func (b *Bot) onPrivateMessage(message twitch.PrivateMessage) {
	log.Printf("%s: %s\n", message.User.DisplayName, message.Message)

	if strings.HasPrefix(message.Message, b.config.CommandPrefix) {
		args := strings.Split(message.Message[1:], " ")

		b.onCommandMessage(strings.ToLower(args[0]), args[1:], message)
	}
}

func (b *Bot) onCommandMessage(command string, args []string, message twitch.PrivateMessage) {
	messageVariables := []string{
		"%COMMAND_PREFIX%", b.config.CommandPrefix,
		"%USER_ID%", message.User.ID,
		"%USER_NAME%", message.User.Name,
	}

	switch command {
	case "task":
		// The user has not specified a subcommand.
		if len(args) == 0 {
			b.Reply(message.ID, b.config.TaskHelpMessage, strings.NewReplacer(messageVariables...))
			return
		}

		subCommand := strings.ToLower(args[0])

		switch subCommand {
		case "add":
			// The user has not specified a task name.
			if len(args) < 2 {
				b.Reply(message.ID, b.config.TaskHelpMessage, strings.NewReplacer(messageVariables...))
				return
			}

			taskName := strings.Join(args[1:], " ")
			messageVariables = append(messageVariables, []string{"%TASK_NAME%", taskName}...)

			if b.AddTask(message.User, taskName) {
				log.Printf("ADDED %s: %s\n", message.User.DisplayName, taskName)
				b.Reply(message.ID, b.config.TaskAddedMessage, strings.NewReplacer(messageVariables...))
			} else {
				b.Reply(message.ID, b.config.TaskInProgressMessage, strings.NewReplacer(messageVariables...))
			}
		case "edit":
			// The user has not specified a task name.
			if len(args) < 2 {
				b.Reply(message.ID, b.config.TaskHelpMessage, strings.NewReplacer(messageVariables...))
				return
			}

			taskName := strings.Join(args[1:], " ")
			messageVariables = append(messageVariables, []string{"%TASK_NAME%", taskName}...)

			if b.EditTask(message.User.ID, taskName) {
				log.Printf("EDITED %s: %s\n", message.User.DisplayName, taskName)
				b.Reply(message.ID, b.config.TaskEditedMessage, strings.NewReplacer(messageVariables...))
			} else {
				b.Reply(message.ID, b.config.NoTaskMessage, strings.NewReplacer(messageVariables...))
			}
		case "done":
			fallthrough
		case "complete":
			if task, ok := b.RemoveTask(message.User.ID); ok {
				log.Printf("COMPLETED %s: %s\n", message.User.DisplayName, task.Name)

				messageVariables = append(messageVariables, []string{"%TASK_NAME%", task.Name}...)

				b.Reply(message.ID, b.config.TaskCompletedMessage, strings.NewReplacer(messageVariables...))
			} else {
				b.Reply(message.ID, b.config.NoTaskMessage, strings.NewReplacer(messageVariables...))
			}
		case "delete":
			fallthrough
		case "cancel":
			if task, ok := b.RemoveTask(message.User.ID); ok {
				log.Printf("CANCELLED %s: %s\n", message.User.DisplayName, task.Name)

				messageVariables = append(messageVariables, []string{"%TASK_NAME%", task.Name}...)

				b.Reply(message.ID, b.config.TaskCancelledMessage, strings.NewReplacer(messageVariables...))
			} else {
				b.Reply(message.ID, b.config.NoTaskMessage, strings.NewReplacer(messageVariables...))
			}
		default:
			// The user has not specified a valid subcommand.
			b.Reply(message.ID, b.config.TaskHelpMessage, strings.NewReplacer(messageVariables...))
			return
		}
	}

	b.UpdateTasksFile()
}
