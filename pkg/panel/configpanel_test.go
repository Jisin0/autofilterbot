package panel_test

import (
	"testing"

	"github.com/Jisin0/autofilterbot/pkg/panel"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/stretchr/testify/assert"
)

func mockCallbackFunc(ctx *panel.Context) (string, [][]gotgbot.InlineKeyboardButton, error) {
	return "foo", nil, nil
}

func mockCallbackQuery(data string) *ext.Context {
	return &ext.Context{
		Update: &gotgbot.Update{
			CallbackQuery: &gotgbot.CallbackQuery{
				Data: data,
			},
		},
	}
}

func countButtons(m [][]gotgbot.InlineKeyboardButton) int {
	if len(m) == 0 {
		return 0
	}

	var count int
	for _, row := range m {
		count += len(row)
	}

	return count
}

func TestPanel(t *testing.T) {
	assert := assert.New(t)

	p := panel.NewPanel()

	p.NewPage("pg1", "Page 1").WithCallbackFunc(mockCallbackFunc)
	pg2 := p.NewPage("pg2", "Page 2").WithContent("test")
	pg2.NewSubPage("sp1", "Sub Page 1").WithContent("test1")

	// Expected result is a panel with two buttons labelled "Page 1" & "Page 2" respectively.
	// Page 1 should return content "foo".
	// Page 2 should have content "test" and a keyboard towards the subpage "Sub Page 1"

	table := []struct {
		data        string // input callback data
		text        string // expected output text
		buttonCount int    // expected number of buttons including back/close buttons
		err         error  // expected error
	}{
		{
			data: "config:pg1|fakearg1_fakearg2",
			text: "foo",
			err:  nil,
		},
		{
			data:        "config:pg2",
			text:        "test",
			buttonCount: 2,
			err:         nil,
		},
		{
			data:        "config:pg2:sp1|",
			text:        "test1",
			buttonCount: 1,
			err:         nil,
		},
		{
			data: "config:pg2:sp3",
			err:  panel.PageNotFoundError{PageName: "sp3"},
		},
	}

	for _, item := range table {
		t.Run(item.data, func(t *testing.T) {
			c, m, e := panel.ProcessUpdate(p, mockCallbackQuery(item.data), nil)

			if item.text != "" {
				assert.Equal(item.text, c)
			}

			assert.Equal(item.buttonCount, countButtons(m))
			assert.Equal(item.err, e)
		})
	}
}
