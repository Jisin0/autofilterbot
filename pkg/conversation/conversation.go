/*
Package conversation is used to create interactive conversations with a user.
*/
package conversation

import (
	"context"
	"errors"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters"
)

var (
	ErrNilUpdate         = errors.New("update is nil")
	ErrUnsupportedUpdate = errors.New("update is not message or callback_query")
)

const (
	FiveMinutes = 5 * time.Minute
)

// Conversator provides easy methods to use conversation methods.
type Conversator struct {
	bot           *gotgbot.Bot
	chatId        int64
	userId        int64
	lastMessageId int64
}

// NewConversator creates a new Conversator from the given chatId, userId & messageId.
//
// NB: msgId is the id of the last message received from the user and is used to filter new messages only.
func NewConversator(bot *gotgbot.Bot, chatId, userId, msgId int64) *Conversator {
	return &Conversator{
		chatId:        chatId,
		userId:        userId,
		lastMessageId: msgId,
	}
}

// NewConversatorFromUpdate creates a new conversator from a given update, either message or callbbackquery.
//
// NOTE: Conversator will be nil if update is not valid.
func NewConversatorFromUpdate(bot *gotgbot.Bot, update *gotgbot.Update) *Conversator {
	if update == nil {
		return nil
	}

	switch {
	case update.CallbackQuery != nil:
		c := update.CallbackQuery

		return &Conversator{
			bot:           bot,
			chatId:        c.Message.GetChat().Id,
			userId:        c.From.Id,
			lastMessageId: c.Message.GetMessageId(),
		}
	case update.Message != nil:
		m := update.Message

		return &Conversator{
			bot:           bot,
			chatId:        m.Chat.Id,
			userId:        m.From.Id,
			lastMessageId: m.MessageId,
		}
	default:
		return nil
	}
}

// Listen listens for incoming messages in current conversation. returns a ErrListenTimeout if deadline exceded.
func (c *Conversator) Listen(ctx context.Context, d ...time.Duration) (*gotgbot.Message, error) {
	expiryDuration := FiveMinutes
	if len(d) != 0 {
		expiryDuration = d[0]
	}

	ctx, cancel := context.WithTimeout(ctx, expiryDuration)
	defer cancel()

	m, err := Listen(ctx, c.ListenFilter())
	if err == nil && m != nil {
		// update lastMessageId for next request
		c.lastMessageId = m.MessageId
	}

	return m, err
}

// Ask sends a message to the chat and waits for a reply. returns a ErrListenTimeout if deadline exceded:
//   - text: Text content of the message.
//   - opts: optional params for SendMessage.
func (c *Conversator) Ask(ctx context.Context, text string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
	if opts == nil {
		opts = &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML}
	} else if opts.ParseMode == "" {
		opts.ParseMode = gotgbot.ParseModeHTML
	}

	firstM, err := c.bot.SendMessage(c.chatId, text, opts)
	if err != nil {
		return nil, err
	}

	c.lastMessageId = firstM.MessageId

	return c.Listen(ctx)
}

// ListenFilter returns a filters.Message to match incoming messages in l.chatId from l.userId.
func (c *Conversator) ListenFilter() filters.Message {
	return listenFilter(c.chatId, c.userId, c.lastMessageId)
}

var (
	ErrListenTimeout = errors.New("Listen context deadline exceded")
)

// Listen listens for incoming message that matches the given filter.
//
// - ctx: Context with ideally a deadline set.
// - filter: The filter that the message must match.
//
// Returns either the message if it was received or a ListenTimeout error
func Listen(ctx context.Context, filter filters.Message) (*gotgbot.Message, error) {
	c := make(chan *gotgbot.Message, 1)

	ActiveListeners.Add(NewListener(filter, c))

	// listen for either a message or timeout
	for {
		select {
		case <-ctx.Done():
			return nil, ErrListenTimeout
		case m := <-c:
			return m, nil
		}
	}
}
