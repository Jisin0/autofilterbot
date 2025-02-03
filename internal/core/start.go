package core

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/format"
	"github.com/Jisin0/autofilterbot/internal/fsub"
	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

const (
	DataPrefixFile  = 'f'
	DataPrefixBatch = 'b'
	DataPrefixRetry = 'r'
)

// StartCommand handles the start command.
func StartCommand(bot *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.Message
	user := m.From

	go func() {
		err := _app.DB.SaveUser(user.Id)
		if err != nil {
			_app.Log.Warn("start: save user failed", zap.Error(err))
		}
	}()

	split := ctx.Args()
	if len(split) < 2 {
		return StaticCommands(bot, ctx)
	}

	bytes, err := base64.StdEncoding.DecodeString(split[1]) // any start data is expected to be base64 encoded
	if err != nil {
		_app.Log.Warn("start: decode data failed", zap.Error(err))
		return nil
	}

	data := string(bytes)
	switch data[0] {
	case DataPrefixFile:
		if f := _app.Config.GetFsubChannels(); len(f) != 0 {
			notJoined, err := fsub.GetNotMemberOrRequest(bot, _app.DB, f, user.Id)
			if err != nil {
				_app.Log.Warn("start: check fsub failed", zap.Error(err))
				return nil
			}

			if len(notJoined) != 0 {
				var btns [][]gotgbot.InlineKeyboardButton

				switch len(notJoined) {
				case 1:
					btns = [][]gotgbot.InlineKeyboardButton{{{Text: "ᴊᴏɪɴ ᴍʏ ᴄʜᴀɴɴᴇʟ", Url: notJoined[0].InviteLink}}}
				case 2:
					btns = [][]gotgbot.InlineKeyboardButton{
						{{Text: "ᴊᴏɪɴ ғɪʀsᴛ ᴄʜᴀɴɴᴇʟ", Url: notJoined[0].InviteLink}},
						{{Text: "ᴊᴏɪɴ sᴇᴄᴏɴᴅ ᴄʜᴀɴɴᴇʟ", Url: notJoined[1].InviteLink}},
					}
				default:
					btns = make([][]gotgbot.InlineKeyboardButton, 0, len(notJoined)+1)
					for i, c := range notJoined {
						btns = append(btns, []gotgbot.InlineKeyboardButton{{Text: fmt.Sprintf("ᴊᴏɪɴ ᴄʜᴀɴɴᴇʟ %d", i+1), Url: c.InviteLink}})
					}
				}

				btns = append(btns, []gotgbot.InlineKeyboardButton{button.Close(user.Id)}, []gotgbot.InlineKeyboardButton{{Text: "𝖱𝖾𝗍𝗋𝗒 🔃"}})

				_, err = m.Reply(bot,
					format.KeyValueFormat(_app.Config.GetFsubText(), _app.BasicMessageValues(ctx)),
					&gotgbot.SendMessageOpts{
						ParseMode:   gotgbot.ParseModeHTML,
						ReplyMarkup: gotgbot.InlineKeyboardMarkup{InlineKeyboard: btns},
					},
				)
				if err != nil {
					_app.Log.Warn("start: send fsub message failed", zap.Error(err))
				}

				return nil
			}
		}

		d, err := autofilter.URLDataFromString(data)
		if err != nil {
			_app.Log.Warn("start: parse sendfile start data failed", zap.Error(err))
			return nil
		}

		f, err := _app.DB.GetFile(d.FileUniqueId)
		if err != nil {
			_app.Log.Warn("start: get file failed", zap.Error(err))
			m.Reply(bot, "<i>📛 I Couldn't File the File You're Looking for, Please Report This to Admins :/</i>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
			return nil
		}

		var (
			warn    string
			delTime = _app.Config.GetFileAutoDelete()
		)
		if delTime != 0 {
			warn = fmt.Sprintf("<blockquote>⚠️ 𝖳𝗁𝗂𝗌 𝖥𝗂𝗅𝖾 𝖶𝗂𝗅𝗅 𝖻𝖾 𝖠𝗎𝗍𝗈𝗆𝖺𝗍𝗂𝖼𝖺𝗅𝗅𝗒 𝖣𝖾𝗅𝖾𝗍𝖾𝖽 𝗂𝗇 %d 𝖬𝗂𝗇𝗎𝗍𝖾𝗌. 𝖥𝗈𝗋𝗐𝖺𝗋𝖽 𝗂𝗍 𝗍𝗈 𝖠𝗇𝗈𝗍𝗁𝖾𝗋 𝖢𝗁𝖺𝗍 𝗈𝗋 𝖲𝖺𝗏𝖾𝖽 𝖬𝖾𝗌𝗌𝖺𝗀𝖾𝗌.</blockquote>", delTime)
		}

		msg, err := f.Send(bot, m.Chat.Id, &model.SendFileOpts{
			Caption: _app.FormatText(ctx, _app.Config.GetFileCaption(), map[string]any{
				"file_size": functions.FileSizeToString(f.FileSize),
				"file_name": f.FileName,
				"warn":      warn,
			}),
			Keyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "🗑️ ᴅᴇʟᴇᴛᴇ ғɪʟᴇ 🗑️", CallbackData: "close"}}},
		})
		if err != nil {
			_app.Log.Warn("start: send file failed", zap.Error(err), zap.String("file_id", f.FileId))
		}

		if delTime != 0 {
			err = _app.AutoDelete.SaveMessage(msg, time.Minute*time.Duration(delTime))
			if err != nil {
				_app.Log.Warn("start: insert auto delete failed", zap.Error(err))
			}
		}
	case DataPrefixRetry:
		d, err := RetryDataFromString(data)
		if err != nil {
			_app.Log.Warn("start: parse retry data failed", zap.Error(err), zap.String("data", data))
			return nil
		}

		url := fmt.Sprintf("https://t.me/c/%d/%d", functions.ChatIdToMtproto(d.ChatId), d.MessageId)
		text := fmt.Sprintf("<b>𝖯𝗅𝖾𝖺𝗌𝖾 𝖧𝖾𝖺𝖽 𝖡𝖺𝖼𝗄 𝗍𝗈 𝖳𝗁𝖾 𝖢𝗁𝖺𝗍 𝖺𝗇𝖽 𝖳𝗋𝗒 𝖠𝗀𝖺𝗂𝗇 <a href='%s'>» 𝖦𝗈 𝖡𝖺𝖼𝗄</a></b>", url)

		_, err = bot.SendMessage(m.Chat.Id, text, &gotgbot.SendMessageOpts{
			ReplyMarkup: gotgbot.InlineKeyboardMarkup{
				InlineKeyboard: [][]gotgbot.InlineKeyboardButton{{{Text: "« ɢᴏ ʙᴀᴄᴋ", Url: url}}},
			},
			ParseMode: gotgbot.ParseModeHTML,
		})
		if err != nil {
			_app.Log.Warn("start: send retry msg failed", zap.Error(err))
		}
	}

	return nil
}
