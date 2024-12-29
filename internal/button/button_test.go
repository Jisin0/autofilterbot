package button_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/stretchr/testify/assert"
)

func stringp(s string) *string {
	return &s
}

func TestUnwrapKeyboard(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input          [][]button.InlineKeyboardButton
		expectedOutput [][]gotgbot.InlineKeyboardButton
	}{
		{
			input: [][]button.InlineKeyboardButton{
				{{Text: "foo", CallbackData: "bar"}, {Text: "hello", Url: "google.com"}},
				{{Text: "test", IsInline: true, SwitchInlineQueryCurrentChat: "hmm"}},
				{{Text: "test1", CopyText: "brr"}},
			},

			expectedOutput: [][]gotgbot.InlineKeyboardButton{
				{{Text: "foo", CallbackData: "bar"}, {Text: "hello", Url: "google.com"}},
				{{Text: "test", SwitchInlineQueryCurrentChat: stringp("hmm")}},
				{{Text: "test1", CopyText: &gotgbot.CopyTextButton{Text: "brr"}}},
			},
		},
	}

	for _, item := range table {
		assert.Equal(item.expectedOutput, button.UnwrapKeyboard(item.input))
	}
}
