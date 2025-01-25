package core

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/format"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// BasicMessageValues creates a map with basic values to format message text with
func (core *Core) BasicMessageValues(ctx *ext.Context, extraValues ...map[string]any) map[string]string {
	var (
		m = ctx.EffectiveMessage
		u = ctx.EffectiveUser
	)

	values := map[string]string{
		"my_name": core.Bot.FirstName,
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

// FormatText formats a string in a pythonic syntax using a key value map.
// basic message values are also added by default.
func (core *Core) FormatText(ctx *ext.Context, template string, extraValues map[string]any) string {
	return format.KeyValueFormat(template, core.BasicMessageValues(ctx, extraValues))
}
