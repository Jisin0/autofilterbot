package cache

import (
	"fmt"
	"time"

	"github.com/Jisin0/autofilterbot/pkg/jsoncache"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Batch facilitates caching message batches.
type Batch struct {
	cache *jsoncache.Cache
}

// NewBatch creates a new batch storage cache.
func NewBatch(timeout time.Duration) *Batch {
	return &Batch{
		cache: jsoncache.NewCache(".batch", timeout),
	}
}

// Save stores a batch of message from a chat.
func (c *Batch) Save(chatId, startMessageId, endMessageId int64, data []*gotgbot.Message) error {
	// remove unnecessary data
	for _, m := range data {
		m.Chat = gotgbot.Chat{} // chat data is not required as msg is reconstructed
		m.From = nil
		m.Date = 0
		m.ReplyToMessage = nil
		m.ViaBot = nil
	}

	return c.cache.Save(fmt.Sprintf("%d-%d-%d", chatId, startMessageId, endMessageId), data)
}

// Get fetches a message batch from storage if available.
//
// If the cache file doesnt exist or was expired a nil error with ok set to false will be returned.
// In case of any other error ok will be true and error will be set.
func (c *Batch) Get(chatId, startMessageId, endMessageId int64) ([]*gotgbot.Message, bool, error) {
	var res []*gotgbot.Message

	err := c.cache.Load(fmt.Sprintf("%d-%d-%d", chatId, startMessageId, endMessageId), &res)
	if err != nil {
		if err == jsoncache.ErrFileNotFound || err == jsoncache.ErrCacheDataExpired {
			return nil, false, nil
		}

		return nil, true, err
	}

	return res, true, nil
}
