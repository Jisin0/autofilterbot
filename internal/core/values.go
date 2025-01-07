package core

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// BasicMessageValues creates a map with basic values to format message text with
func (app *App) BasicMessageValues(ctx *ext.Context, extraValues ...map[string]any) map[string]string {
	m := ctx.EffectiveMessage
	u := ctx.EffectiveUser

	values := map[string]string{
		"my_name": app.Bot.FirstName,
	}

	if u != nil {
		values["first_name"] = u.FirstName
		values["user_id"] = fmt.Sprint(u.Id)

		fullName := u.FirstName

		if u.LastName != "" {
			fullName = fullName + " " + u.LastName
		}
		values["full_name"] = fullName

		var mention string
		if u.Username != "" {
			values["username"] = u.Username
			mention = "@" + u.Username
		} else {
			mention = fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>", u.Id, fullName)
		}

		values["mention"] = mention

	}

	if m.Chat.Title != "" {
		values["chat_name"] = m.Chat.Title
	}

	if m.Chat.Username != "" {
		values["chat_username"] = m.Chat.Username
	}

	if len(extraValues) != 0 {
		for key, val := range extraValues[0] {
			values[key] = fmt.Sprint(val)
		}
	}

	return values
}
