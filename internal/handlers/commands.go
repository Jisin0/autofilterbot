/*
Basic static commands that don't require additional helpers
*/

package handlers

import (
	"strings"

	"github.com/Jisin0/autofilterbot/internal/app"
	"github.com/Jisin0/autofilterbot/pkg/configpanel/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// StaticCommands handles all static text commands like about, help, privacy etc.
// Also handles callback queries in the format cmd:<command_name>
func StaticCommands(app *app.App, ctx *ext.Context, bot *gotgbot.Bot) error {
	var (
		commandName string
		isMedia     bool
	)

	if c := ctx.CallbackQuery; c != nil {
		data := callbackdata.FromString(c.Data)
		if len(data.Path) > 1 {
			commandName = strings.ToLower(data.Path[1])
		}

		switch m := c.Message.(type) {
		case gotgbot.Message:
			isMedia = hasMedia(&m)
		}
	} else {
		m := ctx.EffectiveMessage

		commandName = strings.ToLower(strings.Split(strings.ToLower(strings.Fields(ctx.EffectiveMessage.GetText())[0]), "@")[0][1:])
		isMedia = hasMedia(m)
	}

	var (
		content string
		markup  [][]gotgbot.InlineKeyboardButton //TODO: implement
		err     error
	)

	switch commandName {
	case "about":
		content = app.Config.GetAboutText()
	case "help":
		content = app.Config.GetHelpText()
	case "privacy":
		content = app.Config.GetPrivacyText()
	}

	content = FormatString(content, app.BasicMessageValues(ctx.EffectiveMessage))

	if isMedia {
		_, _, err = ctx.EffectiveMessage.EditCaption(bot, &gotgbot.EditMessageCaptionOpts{Caption: content, ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: markup}})
	} else {
		_, _, err = ctx.EffectiveMessage.EditText(bot, content, &gotgbot.EditMessageTextOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: markup}})
	}

	if err != nil {
		app.Log.Warn(err.Error(), zap.String("command", commandName))
	}

	return nil
}
