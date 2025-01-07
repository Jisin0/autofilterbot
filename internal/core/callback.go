package core

import (
	"strconv"

	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// Close handles callback events from close buttons.
func Close(bot *gotgbot.Bot, ctx *ext.Context) error {
	c := ctx.CallbackQuery

	data := callbackdata.FromString(c.Data)
	if len(data.Args) == 0 { // no authorized users, anyone can use
		c.Message.Delete(bot, nil)
		return nil
	}

	for _, s := range data.Args {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			_app.Log.Warn("close: parse id failed", zap.Error(err), zap.String("id", s))
			continue
		}

		if id == c.From.Id {
			c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Message Has Been Deleted 🗑️"})
			c.Message.Delete(bot, nil)
			return nil
		}
	}

	_, err := c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text:      "You Can't Use This Button",
		ShowAlert: true,
	})
	if err != nil {
		_app.Log.Warn("close: answer query failed", zap.Error(err))
	}

	return nil
}
