package cache_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Jisin0/autofilterbot/internal/cache"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		chatId int64
		start  int64
		end    int64
		data   []*gotgbot.Message
	}{
		{
			chatId: -100182309293,
			start:  69,
			end:    420,
			data: []*gotgbot.Message{
				{Text: "text msg", Entities: []gotgbot.MessageEntity{{Type: "url"}}},
				{Document: &gotgbot.Document{FileId: "AQPk==89SDHKsdjs7", FileName: "file name.mkv", FileSize: 1 << 30}},
			},
		},
		{
			chatId: 1024682446,
			start:  10,
			end:    15,
			data: []*gotgbot.Message{
				{Photo: []gotgbot.PhotoSize{{FileId: "QAQA==90jckSUDHkDSJ_hsIU3s"}}, Caption: "lorem ipsum"},
				{Sticker: &gotgbot.Sticker{FileId: "AQGW746=J__9SHjfdgr"}},
			},
		},
		{
			chatId: -202804923,
			start:  310928,
			end:    520833,
			data: []*gotgbot.Message{
				{Video: &gotgbot.Video{FileId: "1923YIUDHJHJ__SDHK", FileName: "foo bar.mkv", Duration: 1 << 20}},
				{Text: "just text"},
			},
		},
	}

	c := cache.NewBatch(time.Minute * 1)

	for _, item := range table {
		t.Run(fmt.Sprint(item.chatId, item.start, item.end), func(t *testing.T) {
			assert.NoError(c.Save(item.chatId, item.start, item.end, item.data))

			res, _, err := c.Get(item.chatId, item.start, item.end)
			assert.NoError(err)

			assert.Equal(item.data, res)
		})
	}
}
