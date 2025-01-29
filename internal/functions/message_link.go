package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// MessageLink is parsed data from a message link.
type MessageLink struct {
	ChatId    int64
	MessageId int64
	Username  string
}

// GetChat gets a chat using it's id or username.
func (m *MessageLink) GetChat(bot *gotgbot.Bot) (*gotgbot.ChatFullInfo, error) {
	switch {
	case m.ChatId != 0:
		return bot.GetChat(m.ChatId, nil)
	case m.Username != "":
		return GetChatFromUsername(bot, m.Username)
	default:
		return nil, errors.New("both id and username are empty")
	}
}

// ParseMessageLink parses a message link in the format t.me/c/<id> or t.me/<username>.
func ParseMessageLink(s string) (*MessageLink, error) {
	split := strings.Split(s, "/")

	msgId, err := strconv.ParseInt(split[len(split)-1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse message id failed: %w", err)
	}

	var username string

	chatId, err := strconv.ParseInt(split[len(split)-2], 10, 64)
	if err != nil {
		username = split[len(split)-2]
	}

	return &MessageLink{
		ChatId:    chatId,
		MessageId: msgId,
		Username:  username,
	}, nil
}

// GetChatFromUsername constructs a getChat request using a username.
func GetChatFromUsername(bot *gotgbot.Bot, username string) (*gotgbot.ChatFullInfo, error) {
	r, err := bot.Request("getChat", map[string]string{"chat_id": "@" + username}, nil, nil)
	if err != nil {
		return nil, err
	}

	var c gotgbot.ChatFullInfo

	err = json.Unmarshal(r, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
