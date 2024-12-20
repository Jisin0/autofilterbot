package app

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// BasicMessageValues creates a map with basic values to format message text with
func (app *App) BasicMessageValues(m *gotgbot.Message) map[string]string {
	values := map[string]string{
		"my_name": app.Bot.FirstName,
	}

	if u := m.From; u != nil {
		values["first_name"] = u.FirstName
		values["userid"] = fmt.Sprint(u.Id)
		values["username"] = u.Username

		if u.LastName != "" {
			values["full_name"] = u.FirstName + " " + u.LastName
		}
	}

	if m.Chat.Title != "" {
		values["chat_name"] = m.Chat.Title
	}

	if m.Chat.Username != "" {
		values["chat_username"] = m.Chat.Username
	}

	return values
}
