package message

import (
	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Message wraps text and markup for a message making it easy to send.
type Message struct {
	// Text is the text or caption for media using html formatting.
	Text string `json:"text,omitempty"`
	// Keyboard is the inline keyboard for the message.
	Keyboard [][]button.InlineKeyboardButton
}

// Format formats a message from

// Send sends the message to the target chatId using html formatting by default.
func (m Message) Send(bot *gotgbot.Bot, chatId int64, opts ...*gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
	sendOpts := &gotgbot.SendMessageOpts{}

	if len(opts) != 0 {
		sendOpts = opts[0]
	}

	if sendOpts.ParseMode == "" {
		sendOpts.ParseMode = gotgbot.ParseModeHTML
	}

	if len(m.Keyboard) != 0 {
		sendOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{InlineKeyboard: button.UnwrapKeyboard(m.Keyboard)}
	}

	return bot.SendMessage(chatId, m.Text, sendOpts)
}
