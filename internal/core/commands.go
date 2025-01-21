/*
Basic static commands that don't require additional helpers
*/

package core

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/internal/model/message"
	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// StaticCommands handles all static text commands like about, help, privacy etc.
// Also handles callback queries in the format cmd:<command_name>
func StaticCommands(bot *gotgbot.Bot, ctx *ext.Context) error {
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

		if m, ok := c.Message.(gotgbot.Message); ok {
			isMedia = functions.HasMedia(&m)
		}
	} else {
		m := ctx.EffectiveMessage

		commandName = strings.ToLower(strings.Split(strings.ToLower(strings.Fields(ctx.EffectiveMessage.GetText())[0]), "@")[0][1:])
		isMedia = functions.HasMedia(m)
	}

	var (
		msg         *message.Message
		err         error
		extraValues map[string]any
	)

	switch commandName {
	case "start":
		msg = _app.Config.GetStartMessage(bot.Username)
	case "about":
		msg = _app.Config.GetAboutMessage()
		extraValues = map[string]any{
			"os":       runtime.GOOS,
			"database": _app.DB.GetName(),
		}
	case "help":
		msg = _app.Config.GetHelpMessage()
	case "privacy":
		msg = _app.Config.GetPrivacyMessage()
	case "stats": // failsafe
		return Stats(bot, ctx)
	default:
		msg = &message.Message{
			Text: fmt.Sprintf("Command %v Was Not Found!", commandName),
		}
	}

	msg.Format(_app.BasicMessageValues(ctx, extraValues))

	if isCallback {
		if isMedia {
			_, _, err = ctx.EffectiveMessage.EditCaption(bot, &gotgbot.EditMessageCaptionOpts{
				Caption:     msg.Text,
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: button.UnwrapKeyboard(msg.Keyboard)},
				ParseMode:   gotgbot.ParseModeHTML,
			})
		} else {
			_, _, err = ctx.EffectiveMessage.EditText(bot, msg.Text, &gotgbot.EditMessageTextOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: button.UnwrapKeyboard(msg.Keyboard)},
				ParseMode:   gotgbot.ParseModeHTML,
				LinkPreviewOptions: &gotgbot.LinkPreviewOptions{
					IsDisabled: true,
				},
			})
		}
	} else {
		_, err = msg.Send(bot, ctx.EffectiveChat.Id)
	}

	if err != nil {
		_app.Log.Warn(err.Error(), zap.String("command", commandName))
	}

	return nil
}

// Logs handles the /logs command.
func Logs(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !_app.AuthAdmin(ctx) {
		return nil
	}

	m := ctx.EffectiveMessage

	prg, _ := m.Reply(bot, "⏳ 𝖴𝗉𝗅𝗈𝖺𝖽𝗂𝗇𝗀 . . .", nil)

	f, err := os.Open("logs/app.log")
	if err != nil {
		_app.Log.Warn("open log file failed", zap.Error(err))
		return nil
	}

	_, err = bot.SendDocument(
		ctx.EffectiveChat.Id,
		gotgbot.InputFileByReader("app-log.json", f),
		&gotgbot.SendDocumentOpts{
			ReplyParameters: &gotgbot.ReplyParameters{
				MessageId: m.MessageId,
			},
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{button.Close(m.From.Id)}}},
		},
	)
	if err != nil {
		_app.Log.Warn("send log file failed", zap.Error(err))
	}

	if prg != nil {
		prg.Delete(bot, nil)
	}

	return nil
}

// Stats handles the stats command and callback query.
func Stats(bot *gotgbot.Bot, ctx *ext.Context) error {
	s, err := _app.DB.Stats()
	if err != nil {
		_app.Log.Warn("stats: get stats failed", zap.Error(err))
		return nil
	}

	m := _app.Config.GetStatsMessage().Format(_app.BasicMessageValues(ctx, map[string]any{
		"users":  s.Users,
		"files":  s.Files,
		"groups": s.Groups,
		"uptime": time.Since(_app.StartTime).Truncate(time.Second),
	}))

	switch {
	case ctx.Message != nil:
		_, err = m.Send(bot, ctx.EffectiveChat.Id)
	case ctx.CallbackQuery != nil:
		ctx.CallbackQuery.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Generating Stats ...", CacheTime: fiveHoursInSeconds})

		var isMedia bool

		if msg, ok := ctx.CallbackQuery.Message.(*gotgbot.Message); ok {
			isMedia = functions.HasMedia(msg)
		}

		if isMedia {
			_, _, err = ctx.EffectiveMessage.EditCaption(bot, &gotgbot.EditMessageCaptionOpts{
				Caption:     m.Text,
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: button.UnwrapKeyboard(m.Keyboard)},
				ParseMode:   gotgbot.ParseModeHTML,
			})
		} else {
			_, _, err = ctx.EffectiveMessage.EditText(bot, m.Text, &gotgbot.EditMessageTextOpts{
				ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: button.UnwrapKeyboard(m.Keyboard)},
				ParseMode:   gotgbot.ParseModeHTML,
				LinkPreviewOptions: &gotgbot.LinkPreviewOptions{
					IsDisabled: true,
				},
			})
		}
	}

	if err != nil {
		_app.Log.Warn("stats: send result failed", zap.Error(err))
	}

	return nil
}
