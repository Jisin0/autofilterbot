package core

import (
	"fmt"
	"strconv"

	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// Select handles the select button callback query.
func Select(bot *gotgbot.Bot, ctx *ext.Context) error {
	c := ctx.CallbackQuery

	data := callbackdata.FromString(c.Data)
	if len(data.Args) < 2 {
		_app.Log.Warn("select: not enough args", zap.Strings("args", data.Args))
		_, err := c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Error: Not Enough Arguments", ShowAlert: true})
		return err
	}

	pageIndex, err := strconv.Atoi(data.Args[1])
	if err != nil {
		_app.Log.Warn("select: parse index failed", zap.Error(err))
		_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Sorry An Error occurred :/", ShowAlert: true})
		return err
	}

	uniqueId := data.Args[0]

	r, ok, err := _app.Cache.Autofilter.Get(uniqueId)
	if !ok {
		_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Search Result Has Expired!\nPlease Try Again...", ShowAlert: true})
		return err
	}

	if err != nil {
		_app.Log.Warn("select: get result cache failed", zap.Error(err))
		_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Sorry An Error occurred :/", ShowAlert: true})
		return err
	}

	if r.FromUser != c.From.Id {
		_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "You Can't Use This Button!", ShowAlert: true})
		return err
	}

	if pageIndex >= len(r.Files) {
		_app.Log.Warn("select: page not found", zap.Int("length", len(r.Files)), zap.Int("index", pageIndex))
		_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Result Page Not Found!", ShowAlert: true})
		return err
	}

	if len(data.Args) > 2 { // if file uid in args
		r.SelectFile(pageIndex, data.Args[2])
	}

	var (
		pageFiles = r.Files[pageIndex]
		buttons   = make([][]gotgbot.InlineKeyboardButton, 0, len(pageFiles)+2)
	)

	buttons = append(buttons, selectHeaderRow(r.UniqueId, pageIndex))
	buttons = append(buttons, pageFiles.SelectMenu(uniqueId, pageIndex)...)
	buttons = append(buttons, selectFooterRow(r.UniqueId, pageIndex, len(r.Files)))

	_, _, err = c.Message.EditReplyMarkup(bot, &gotgbot.EditMessageReplyMarkupOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: buttons},
	})
	if err != nil {
		_app.Log.Warn("select: edit markup failed", zap.Error(err))
	}

	err = _app.Cache.Autofilter.Save(r)
	if err != nil {
		_app.Log.Warn("select: save result failed", zap.Error(err))
	}

	return nil
}

func selectHeaderRow(uniqueId string, pageIndex int) []gotgbot.InlineKeyboardButton {
	return []gotgbot.InlineKeyboardButton{{Text: "ᴇxɪᴛ", CallbackData: fmt.Sprintf("navg|%s_%d", uniqueId, pageIndex)}, {Text: "sᴇɴᴅ ➡️", CallbackData: fmt.Sprintf("sendsel|%s", uniqueId)}}
}

func selectFooterRow(uniqueId string, pageIndex, totalPages int) []gotgbot.InlineKeyboardButton {
	btns := make([]gotgbot.InlineKeyboardButton, 0, 3)
	if pageIndex != 0 {
		btns = append(btns, selectBackButton(uniqueId, pageIndex-1))
	}

	btns = append(btns, gotgbot.InlineKeyboardButton{Text: fmt.Sprintf("📑 𝗣𝗔𝗚𝗘 %d/%d", pageIndex+1, totalPages), CallbackData: "ignore"})

	if pageIndex+1 != totalPages {
		btns = append(btns, selectNextButton(uniqueId, pageIndex+1))
	}

	return btns
}

func selectBackButton(uniqueId string, pageIndex int) gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{Text: "« ʙᴀᴄᴋ", CallbackData: fmt.Sprintf("sel|%s_%d", uniqueId, pageIndex)}
}

func selectNextButton(uniqueId string, pageIndex int) gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{Text: "ɴᴇxᴛ »", CallbackData: fmt.Sprintf("sel|%s_%d", uniqueId, pageIndex)}
}
