package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/pkg/conversation"
	"github.com/Jisin0/autofilterbot/pkg/send"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// Broadcast handles the /broadcast command to copy msg to all bot users.
func Broadcast(bot *gotgbot.Bot, ctx *ext.Context) error {
	if !_app.AuthAdmin(ctx) {
		return nil
	}

	m := ctx.Message
	var (
		opts   send.SendOpts
		method send.SendMethod
	)

	if replyM := m.ReplyToMessage; replyM != nil {
		sendMethod, text, fileId, err := sendOptsFromMessage(replyM)
		if err != nil {
			m.Reply(bot, "<b>⛔ 𝖱𝖾𝗉𝗅𝗂𝖾𝖽 𝖬𝖾𝗌𝗌𝖺𝗀𝖾 𝖢𝗈𝗇𝗍𝖺𝗂𝗇𝗌 𝖴𝗇𝗌𝗎𝗉𝗉𝗈𝗋𝗍𝖾𝖽 𝖬𝖾𝖽𝗂𝖺!</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
			return nil
		}

		method = sendMethod
		opts.Text += text
		opts.FileId = fileId

		if replyM.ReplyMarkup != nil && len(replyM.ReplyMarkup.InlineKeyboard) != 0 {
			opts.Keyboard = append(opts.Keyboard, replyM.ReplyMarkup.InlineKeyboard...)
		}
	}

	split := strings.SplitN(m.OriginalHTML(), " ", 2)
	if len(split) > 1 {
		opts.Text += " " + split[1]
		if method == nil {
			method = send.SendMessage
		}
	}

	if method == nil {
		m, err := conversation.NewConversatorFromUpdate(bot, ctx.Update).Ask("<b>𝖯𝗅𝖾𝖺𝗌𝖾 𝖲𝖾𝗇𝖽 𝗍𝗁𝖾 𝖬𝖾𝗌𝗌𝖺𝗀𝖾 𝗍𝗈 𝖻𝖾 𝖡𝗋𝗈𝖺𝖽𝖼𝖺𝗌𝗍𝖾𝖽:</b>", nil)
		if err != nil {
			return nil
		}

		sendMethod, text, fileId, err := sendOptsFromMessage(m)
		if err != nil {
			m.Reply(bot, "<b>⛔ 𝖬𝖾𝗌𝗌𝖺𝗀𝖾 𝖢𝗈𝗇𝗍𝖺𝗂𝗇𝗌 𝖴𝗇𝗌𝗎𝗉𝗉𝗈𝗋𝗍𝖾𝖽 𝖬𝖾𝖽𝗂𝖺!</b>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
			return nil
		}

		method = sendMethod
		opts.Text += text
		opts.FileId = fileId

		if m.ReplyMarkup != nil && len(m.ReplyMarkup.InlineKeyboard) != 0 {
			opts.Keyboard = append(opts.Keyboard, m.ReplyMarkup.InlineKeyboard...)
		}
	}

	parsedText, keyboard, err := button.ParseFromText(opts.Text)
	if err != nil {
		m.Reply(bot, fmt.Sprintf("<b>𝖯𝖺𝗋𝗌𝗂𝗇𝗀 𝖡𝗎𝗍𝗍𝗈𝗇𝗌 𝖥𝖺𝗂𝗅𝖾𝖽 🙁</b>\nError: <code>%s</code>", err.Error()), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		_app.Log.Debug("broadcast: parse buttons failed", zap.Error(err), zap.String("text", opts.Text))
		return nil
	}
	opts.Text = parsedText
	opts.Keyboard = append(opts.Keyboard, button.UnwrapKeyboard(keyboard)...)

	users, err := _app.DB.GetAllUsers()
	if err != nil {
		m.Reply(bot, "Fetch Users From Database Failed :/", nil)
		_app.Log.Warn("broadcast: get all users failed", zap.Error(err))
		return nil
	}

	progressM, err := bot.SendMessage(m.Chat.Id, "Sᴛᴀʀᴛɪɴɢ Bʀᴏᴀᴅᴄᴀsᴛ...", nil)
	if err != nil {
		_app.Log.Warn("broadcast: send progress msg failed", zap.Error(err))
		return nil
	}

	p := newBroadcastProgress()

	for users.Next(context.Background()) {
		var u model.User

		err = users.Decode(&u)
		if err != nil {
			_app.Log.Warn("braodcast: decode user failed", zap.Error(err))
			continue
		}

		_, err = method(bot, u.UserId, &opts)
		if err != nil {
			p.failed++

			errStr := err.Error()
			switch {
			case strings.Contains(errStr, "blocked"):
				_app.DB.DeleteUser(u.UserId)
				p.blocked++
			case strings.Contains(errStr, "deleted"): //TODO: not sure what error msg for deleted acc is
				_app.DB.DeleteUser(u.UserId)
				p.deleted++
			case strings.Contains(errStr, "chat not found"):
				_app.DB.DeleteUser(u.UserId)
				fallthrough
			default:
				p.otherErr++
				_app.Log.Info("broadcast: failed to send", zap.Int64("chat_id", u.UserId), zap.Error(err))
			}
		} else {
			p.success++
		}

		p.total++

		if p.total%50 == 0 {
			_, _, err = progressM.EditText(
				bot,
				p.BuildMessage().String(),
				&gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML},
			)
		}
		if err != nil {
			_app.Log.Debug("broadcast: update progress failed", zap.Error(err))
		}

		if _app.Ctx.Err() != nil {
			progressM.EditText(
				bot,
				p.BuildMessage().WriteLn("<code>Broadcast Cancelled Due to Application Stopping</code>").String(),
				&gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML},
			)
			break
		}
	}

	_, _, err = progressM.EditText(
		bot,
		p.BuildMessage().WriteLn("<code>Broadcast Completed Successfully ✅</code>").String(),
		&gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML},
	)
	if err != nil {
		_app.Log.Warn("broadcast: update success msg failed", zap.Error(err))
	}

	return nil
}

type broadcastProgress struct {
	total    int
	success  int
	failed   int
	blocked  int
	deleted  int
	otherErr int
}

func newBroadcastProgress() *broadcastProgress {
	return &broadcastProgress{}
}

type broadcastProgressBuilder struct {
	strings.Builder
}

// WriteLn writes a string to the buffer after a new line.
func (b *broadcastProgressBuilder) WriteLn(s string) *broadcastProgressBuilder {
	b.WriteString("\n" + s)
	return b
}

func (p *broadcastProgress) BuildMessage() *broadcastProgressBuilder {
	var b broadcastProgressBuilder

	b.WriteString(fmt.Sprintf(`<b>𝖡𝗋𝗈𝖺𝖽𝖼𝖺𝗌𝗍 𝖯𝗋𝗈𝗀𝗋𝖾𝗌𝗌</b>
𝖳𝗈𝗍𝖺𝗅: %d
𝖲𝗎𝖼𝖼𝖾𝗌𝗌: %d
<blockquote>𝖥𝖺𝗂𝗅𝖾𝖽: %d
	𝖡𝗅𝗈𝖼𝗄𝖾𝖽: %d
	𝖣𝖾𝗅𝖾𝗍𝖾𝖽 %d
	𝖮𝗍𝗁𝖾𝗋: %d</blockquote>`, p.total, p.success, p.failed, p.blocked, p.deleted, p.otherErr))

	return &b
}

// sendOptsFromMessage gets message send message opts from given message.
//
// Error is only returned if message has no supported media or text.
func sendOptsFromMessage(m *gotgbot.Message) (method send.SendMethod, text, fileId string, err error) {
	switch {
	case m.Document != nil:
		method = send.SendDocument
		fileId = m.Document.FileId
	case m.Video != nil:
		method = send.SendVideo
		fileId = m.Video.FileId
	case m.Audio != nil:
		method = send.SendAudio
		fileId = m.Audio.FileId
	case m.Photo != nil:
		method = send.SendPhoto
		fileId = m.Photo[0].FileId
	case m.Animation != nil:
		method = send.SendAnimation
		fileId = m.Animation.FileId
	case m.Text != "":
		method = send.SendMessage
		text = m.OriginalHTML()
	default:
		err = errors.New("unsupported media type")
	}

	if m.Caption != "" {
		text = m.OriginalCaptionHTML()
	}

	return
}
