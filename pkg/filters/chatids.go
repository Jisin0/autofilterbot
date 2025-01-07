package exthandlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

// ChatIds filters messages from any of the given chat ids.
func ChatIds(l []int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		for _, id := range l {
			if id == m.Chat.Id {
				return true
			}
		}

		return false
	}
}

// UserIds filters messages from any of the given user ids.
func UserIds(l []int64) filters.Message {
	return func(m *gotgbot.Message) bool {
		for _, id := range l {
			if id == m.From.Id {
				return true
			}
		}

		return false
	}
}
