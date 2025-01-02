/*
Basic static commands that don't require additional helpers
*/

package handlers

import (
	"fmt"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/app"
	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/model/message"
	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
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

	isCallback := ctx.CallbackQuery != nil
	if c := ctx.CallbackQuery; isCallback {
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
		msg *message.Message
		err error
	)

	switch commandName {
	case "start":
		msg = app.Config.GetStartMessage(bot.Username)
	case "about":
		msg = app.Config.GetAboutMessage()
	case "help":
		msg = app.Config.GetHelpMessage()
	case "privacy":
		msg = app.Config.GetPrivacyMessage()
	default:
		msg = &message.Message{
			Text: fmt.Sprintf("Commsnd %v Was Not Found!", commandName),
		}
	}

	msg.Format(app.BasicMessageValues(ctx.EffectiveMessage))

	if isCallback {
		if isMedia {
			_, _, err = ctx.EffectiveMessage.EditCaption(bot, &gotgbot.EditMessageCaptionOpts{Caption: msg.Text, ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: button.UnwrapKeyboard(msg.Keyboard)}})
		} else {
			_, _, err = ctx.EffectiveMessage.EditText(bot, msg.Text, &gotgbot.EditMessageTextOpts{ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: button.UnwrapKeyboard(msg.Keyboard)}})
		}
	} else {
		_, err = msg.Send(bot, ctx.EffectiveChat.Id)
	}

	if err != nil {
		app.Log.Warn(err.Error(), zap.String("command", commandName))
	}

	return nil
}
