package button_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/stretchr/testify/assert"
)

func TestParseFromText(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		id    string
		input string

		outputText string
		outputBtn  [][]button.InlineKeyboardButton
		isErr      bool
	}{
		{
			id: "header",
			input: `some header text
			[Label0](url:google.com) [Label1](cmd:privacy)
			[Label2](inline:)
			[Label3](copy:foo)`,

			outputText: "some header text",
			outputBtn: [][]button.InlineKeyboardButton{
				{{Text: "Label0", Url: "google.com"}, {Text: "Label1", CallbackData: "cmd:privacy"}},
				{{Text: "Label2", IsInline: true, SwitchInlineQueryCurrentChat: ""}},
				{{Text: "Label3", CopyText: "foo"}},
			},
		},
		{
			id: "footer",
			input: `
			[foo bar](cmd:lorem)
			[mmmm](url:example.com)[Link](url:t.me)
			footer text hmmm`,

			outputText: "footer text hmmm",
			outputBtn: [][]button.InlineKeyboardButton{
				{{Text: "foo bar", CallbackData: "cmd:lorem"}},
				{{Text: "mmmm", Url: "example.com"}, {Text: "Link", Url: "t.me"}},
			},
		},
		{
			id: "notext",
			input: `
			
			[Button](inline:dayum)`,

			outputText: "",
			outputBtn:  [][]button.InlineKeyboardButton{{{Text: "Button", IsInline: true, SwitchInlineQueryCurrentChat: "dayum"}}},
		},
		{
			id:    "witherror",
			input: `blah blah blah [Lebel](url:)`,

			isErr: true,
		},
	}

	for _, item := range table {
		t.Run(item.id, func(t *testing.T) {
			txt, btn, err := button.ParseFromText(item.input)
			if err != nil {
				if !item.isErr {
					assert.Fail(err.Error(), item)
				}

				return
			}

			assert.Equal(item.outputText, txt)
			assert.Equal(item.outputBtn, btn)
		})
	}
}
