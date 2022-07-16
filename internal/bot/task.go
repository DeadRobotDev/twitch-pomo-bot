package bot

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v3"
)

type User struct {
	ID          string
	DisplayName string
}

type Task struct {
	User User
	Name string
}

func newTask(user twitch.User, taskName string) *Task {
	return &Task{
		User: User{
			ID:          user.ID,
			DisplayName: user.DisplayName,
		},
		Name: taskName,
	}
}

func (t *Task) String() string {
	return fmt.Sprintf("%s: %s\n", t.User.DisplayName, t.Name)
}
