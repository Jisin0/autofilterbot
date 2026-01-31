package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Jisin0/autofilterbot/internal/fsub"
	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// All handles the callback from the "all" button in autofilter results.
func All(bot *gotgbot.Bot, ctx *ext.Context) error {
	c := ctx.CallbackQuery

	data := callbackdata.FromString(c.Data)
	if len(data.Args) < 2 {
		_app.Log.Warn("all: not enough args", zap.Strings("args", data.Args))
		c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Error: Not Enough Arguments", ShowAlert: true})

		return nil
	}

	pageIndex, err := strconv.Atoi(data.Args[1])
	if err != nil {
		_app.Log.Warn("all: parse index failed", zap.Error(err))
		c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Sorry An Error occurred :/", ShowAlert: true})

		return nil
	}

	uniqueId := data.Args[0]

	r, ok, err := _app.Cache.Autofilter.Get(uniqueId)
	if !ok {
		c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Search Result Has Expired!\nPlease Try Again...", ShowAlert: true})
		return nil
	}

	if err != nil {
		_app.Log.Warn("all: get result cache failed", zap.Error(err))
		c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Sorry An Error occurred :/", ShowAlert: true})

		return nil
	}

	if r.FromUser != c.From.Id {
		c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "You Can't Use This Button!", ShowAlert: true})
		return nil
	}

	if pageIndex >= len(r.Files) {
		_app.Log.Warn("all: page not found", zap.Int("length", len(r.Files)), zap.Int("index", pageIndex))
		c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{Text: "Result Page Not Found!", ShowAlert: true})

		return nil
	}

	ok, err = fsub.CheckFsub(_app, bot, ctx)
	if err != nil {
		if functions.IsChatNotFoundErr(err) { // user has not started bot or blocked
			// redirect to dm for a retry msg
			data := &RetryData{
				ChatId:    c.Message.GetChat().Id,
				MessageId: c.Message.GetMessageId(),
			}

			_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
				Url: fmt.Sprintf("t.me/%s?start=%s", bot.Username, data.Encode()),
			})
			if err != nil {
				_app.Log.Warn("all: retry answer failed", zap.Error(err))
			}

			return nil
		}

		_app.Log.Warn("all: check fsub failed", zap.Error(err))
	}

	if !ok {
		return nil
	}

	pageFiles := r.Files[pageIndex]

	sentMessages := make([]struct {
		chatId    int64
		messageId int64
	}, 0, len(pageFiles))

	var (
		warn    string
		delTime = _app.Config.GetFileAutoDelete()
	)
	if delTime != 0 {
		warn = fmt.Sprintf("<blockquote><b><i>⚠️ 𝖳𝗁𝗂𝗌 𝖥𝗂𝗅𝖾 𝖶𝗂𝗅𝗅 𝖻𝖾 𝖠𝗎𝗍𝗈𝗆𝖺𝗍𝗂𝖼𝖺𝗅𝗅𝗒 𝖣𝖾𝗅𝖾𝗍𝖾𝖽 𝗂𝗇 %d 𝖬𝗂𝗇𝗎𝗍𝖾𝗌. 𝖥𝗈𝗋𝗐𝖺𝗋𝖽 𝗂𝗍 𝗍𝗈 𝖠𝗇𝗈𝗍𝗁𝖾𝗋 𝖢𝗁𝖺𝗍 𝗈𝗋 𝖲𝖺𝗏𝖾𝖽 𝖬𝖾𝗌𝗌𝖺𝗀𝖾𝗌.</i></b></blockquote>", delTime)
	}

	for _, f := range pageFiles {
		msg, err := f.Send(bot, c.From.Id, &model.SendFileOpts{
			Caption: _app.FormatText(ctx, _app.Config.GetFileCaption(), map[string]any{
				"file_size": functions.FileSizeToString(f.FileSize),
				"file_name": f.FileName,
				"warn":      warn,
			}),
			Keyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "🗑️ ᴅᴇʟᴇᴛᴇ ғɪʟᴇ 🗑️", CallbackData: "close"}}},
		})
		if err != nil {
			if functions.IsChatNotFoundErr(err) { // user has not started bot or blocked
				// redirect to dm for a retry msg
				data := &RetryData{ //TODO: implement
					ChatId:    c.Message.GetChat().Id,
					MessageId: c.Message.GetMessageId(),
				}

				_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
					Url: fmt.Sprintf("t.me/%s?start=%s", bot.Username, data.Encode()),
				})
				if err != nil {
					_app.Log.Warn("all: retry answer failed", zap.Error(err))
				}

				return nil
			}

			_app.Log.Warn("all: send file failed", zap.Error(err), zap.String("file_id", f.FileId))

			continue
		}

		sentMessages = append(sentMessages, struct {
			chatId    int64
			messageId int64
		}{chatId: msg.Chat.Id, messageId: msg.MessageId})
	}

	_, err = c.Answer(bot, &gotgbot.AnswerCallbackQueryOpts{
		Text:      fmt.Sprintf("%d ғɪʟᴇs ʜᴀᴠᴇ ʙᴇᴇɴ sᴇɴᴛ ᴘʀɪᴠᴀᴛᴇʟʏ 🥳", len(sentMessages)),
		ShowAlert: true,
	})
	if err != nil {
		_app.Log.Warn("all: answer query failed", zap.Error(err))
	}

	if delTime != 0 {
		duration := time.Minute * time.Duration(delTime)

		for _, m := range sentMessages {
			err = _app.AutoDelete.Save(m.chatId, m.messageId, duration)
			if err != nil {
				_app.Log.Warn("all: save autodelete failed", zap.Error(err))
			}
		}
	}

	return nil
}

// RetryData is start data for a retry message, usually from an all or select option when user has not started the bot first.
type RetryData struct {
	// Chat to return to.
	ChatId int64
	// Id of message to return to.
	MessageId int64
}

// Encode encodes it to a base64 string.
func (d *RetryData) Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("r_%d_%d", d.ChatId, d.MessageId)))
}

// RetryDataFromString converts start data into RetryData structure.
func RetryDataFromString(s string) (*RetryData, error) {
	split := strings.Split(s, "_")
	if len(split) < 3 {
		return nil, errors.New("not enough arguments")
	}

	chatId, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return nil, err
	}

	msgId, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return nil, err
	}

	return &RetryData{ChatId: chatId, MessageId: msgId}, nil
}
