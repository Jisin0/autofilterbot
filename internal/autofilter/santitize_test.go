package autofilter_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/stretchr/testify/assert"
)

func TestIsBadQuery(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		id       string
		text     string
		entities []gotgbot.MessageEntity
		output   bool
	}{
		{
			id:     "long",
			text:   "This is a very very very very very very long message.... Wanna know what else is long? ( ｡ •̀ ᴗ •́ ｡)",
			output: true,
		},
		{
			id:     "small",
			text:   "k",
			output: true,
		},
		{
			id:     "normal",
			text:   "An Average Sized Query",
			output: false,
		},
		{
			id:       "bad-entity",
			text:     "Good Text",
			entities: []gotgbot.MessageEntity{{Type: "text_mention"}},
			output:   true,
		},
	}

	for _, item := range table {
		t.Run(item.id, func(t *testing.T) {
			assert.Equal(item.output, autofilter.IsBadQuery(item.text, item.entities))
		})
	}
}
