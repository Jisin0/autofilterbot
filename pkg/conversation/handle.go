package conversation

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

var ActiveListeners = NewListenerArray()

// MessageHandler filters message based on active listeners.
func MessageHandler(bot *gotgbot.Bot, update *ext.Context) error {
	// run in seperate goroutine to prevent blocking the thread
	go func() {
		m := update.Message

		l, ok := ActiveListeners.FindMatchAndDelete(m)
		if ok {
			l.messageChan <- update.Message
		}
	}()

	return nil
}

// listenFilter creates a filters.Message to listen match messages from a specific chatId from a user after the given messageId.
//
// - chatId: Target chat's id.
// - userId: Id of expected user.
// - lastMessageId: id of the last message in the chat
func listenFilter(chatId, userId, lastMessageId int64) filters.Message {
	return func(msg *gotgbot.Message) bool {
		return msg.From.Id == userId && msg.MessageId > lastMessageId && msg.Chat.Id == chatId
	}
}
