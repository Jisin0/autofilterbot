package core

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/pkg/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// NewBatch handles the /batch commmand.
func NewBatch(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !_app.AuthAdmin(ctx) {
		return nil
	}

	m := ctx.Message

	var (
		fromChatId, startId, endId int64
	)

	if replyM := m.ReplyToMessage; replyM != nil {
		if origin, ok := replyM.ForwardOrigin.(gotgbot.MessageOriginChannel); ok {
			fromChatId = origin.Chat.Id
			startId = origin.MessageId
		} else if link, err := functions.ParseMessageLink(replyM.Text); err == nil {
			if c, err := link.GetChat(bot); err == nil {
				fromChatId = c.Id
				startId = link.MessageId
			}
		}
	}

	split := strings.Fields(m.Text)
	if len(split) > 1 {
		if link, err := functions.ParseMessageLink(split[1]); err == nil {
			if startId != 0 {
				endId = link.MessageId
			} else {
				if c, err := link.GetChat(bot); err == nil {
					fromChatId = c.Id
					startId = link.MessageId
				}
			}
		}

		if len(split) > 2 && endId == 0 {
			if link, err := functions.ParseMessageLink(split[2]); err == nil {
				if startId != 0 {
					endId = link.MessageId
				} else {
					if c, err := link.GetChat(bot); err == nil {
						fromChatId = c.Id
						startId = link.MessageId
					}
				}
			}
		}
	}

	if startId == 0 {
		conv := conversation.NewConversatorFromUpdate(bot, ctx.Update)

		askM, err := conv.Ask("Please forward or send the post link of the first message in the batch:", nil)
		if err != nil {
			_app.Log.Debug("batch: conv exited with error", zap.Error(err))
			return nil
		}

		if origin, ok := askM.ForwardOrigin.(gotgbot.MessageOriginChannel); ok {
			fromChatId = origin.Chat.Id
			startId = origin.MessageId
		} else if link, err := functions.ParseMessageLink(askM.Text); err == nil {
			if c, err := link.GetChat(bot); err == nil {
				fromChatId = c.Id
				startId = link.MessageId
			}
		} else {
			askM.Reply(bot, "Message Is Not a Forwarded Channel Post or Message Link!", nil)
			return nil
		}
	}

	if endId == 0 {
		conv := conversation.NewConversatorFromUpdate(bot, ctx.Update)

		askM, err := conv.Ask("Please forward or send the post link of the last message in the batch:", nil)
		if err != nil {
			_app.Log.Debug("batch: conv exited with error", zap.Error(err))
			return nil
		}

		if origin, ok := askM.ForwardOrigin.(gotgbot.MessageOriginChannel); ok {
			endId = origin.MessageId
		} else if link, err := functions.ParseMessageLink(askM.Text); err == nil {
			endId = link.MessageId
		} else {
			askM.Reply(bot, "Message Is Not a Forwarded Channel Post or Message Link!", nil)
			return nil
		}
	}

	if startId > endId {
		m.Reply(bot, "First Message Cannot be After The Last :/", nil)
		return nil
	}

	if endId-startId > _app.Config.GetBatchSizeLimit() {
		m.Reply(bot, "Batch Too Large :/\n\nCreate a Smaller Batch or Update The Batch Size Limit From the Config Panel!", nil)
		_app.Log.Debug("batch: too large", zap.Int64("chat_id", fromChatId), zap.Int64("start", startId), zap.Int64("end", endId))
		return nil
	}

	data := &BatchURLData{
		ChatId:         fromChatId,
		StartMessageId: startId,
		EndMessageId:   endId,
	}
	url := fmt.Sprintf("https://t.me/%s?start=%s", bot.Username, data.Encode())

	text := fmt.Sprintf(`
<b>𝖬𝖾𝗌𝗌𝖺𝗀𝖾 𝖡𝖺𝗍𝖼𝗁 𝖧𝖺𝗌 𝖡𝖾𝖾𝗇 𝖢𝗋𝖾𝖺𝗍𝖾𝖽 𝖲𝗎𝖼𝖼𝖾𝗌𝗌𝖿𝗎𝗅𝗅𝗒 🎉</b>
<b>𝖳𝗋𝗒 𝖭𝗈𝗐:</b> <a href='%s'>ᴄʟɪᴄᴋ ʜᴇʀᴇ</a>
<b>𝖢𝗈𝗉𝗒:</b> <code>%s</code>
<blockquote>⚠️ 𝖳𝗁𝗂𝗌 𝖡𝖺𝗍𝖼𝗁 𝖶𝗂𝗅𝗅 𝖮𝗇𝗅𝗒 𝖶𝗈𝗋𝗄 𝖺𝗌 𝖫𝗈𝗇𝗀 𝖺𝗌 𝖳𝗁𝖾 𝖡𝗈𝗍 𝗂𝗌 𝖺 𝖬𝖾𝗆𝖻𝖾𝗋 𝗈𝖿 𝖳𝗁𝖾 𝖢𝗁𝖺𝗇𝗇𝖾𝗅 𝖺𝗇𝖽 𝖬𝖾𝗌𝗌𝖺𝗀𝖾𝗌 𝖠𝗋𝖾 𝖭𝗈𝗍 𝖣𝖾𝗅𝖾𝗍𝖾𝖽</blockquote>
<blockquote>ℹ️ 𝖢𝗁𝖺𝗇𝗀𝖾𝗌 𝗈𝗋 𝖤𝖽𝗂𝗍𝗌 𝗂𝗇 𝖮𝗋𝗂𝗀𝗂𝗇𝖺𝗅 𝖬𝖾𝗌𝗌𝖺𝗀𝖾𝗌 𝖶𝗂𝗅𝗅 𝖭𝗈𝗍 𝖻𝖾 𝖱𝖾𝖿𝗅𝖾𝖼𝗍𝖾𝖽 𝖨𝗆𝗆𝖾𝖽𝗂𝖺𝗍𝖾𝗅𝗒 𝖺𝗌 𝖬𝖾𝗌𝗌𝖺𝗀𝖾𝗌 𝖺𝗋𝖾 𝖢𝖺𝖼𝗁𝖾𝖽 𝖿𝗈𝗋 𝖤𝖿𝖿𝗂𝖼𝗂𝖾𝗇𝖼𝗒</blockquote>`, url, url)
	btn := [][]gotgbot.InlineKeyboardButton{
		{{Text: "𝖳𝗋𝗒 𝖭𝗈𝗐", Url: url}},
		{{Text: "𝖳𝖺𝗉 𝗍𝗈 𝖢𝗈𝗉𝗒", CopyText: &gotgbot.CopyTextButton{Text: url}}},
	}

	_, err := bot.SendMessage(m.Chat.Id, text, &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: btn},
		ParseMode:   gotgbot.ParseModeHTML,
	})
	if err != nil {
		_app.Log.Warn("batch: send success msg failed", zap.Error(err))
	}

	return nil
}

// BatchURLData is the url data from the start command.
type BatchURLData struct {
	ChatId         int64
	StartMessageId int64
	EndMessageId   int64
}

func (d *BatchURLData) Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("b_%d_%d_%d", d.ChatId, d.StartMessageId, d.EndMessageId)))
}

// BatchURLDataFromString converts a string, usually the start data, to batch url data.
func BatchURLDataFromString(s string) (*BatchURLData, error) {
	split := strings.Split(s, "_")
	if len(split) < 4 {
		return nil, errors.New("not enough arguments")
	}

	chatId, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return nil, err
	}

	startId, err := strconv.ParseInt(split[2], 10, 64)
	if err != nil {
		return nil, err
	}

	endId, err := strconv.ParseInt(split[3], 10, 64)
	if err != nil {
		return nil, err
	}

	return &BatchURLData{
		ChatId:         chatId,
		StartMessageId: startId,
		EndMessageId:   endId,
	}, nil
}
