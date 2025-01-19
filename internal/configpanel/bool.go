package configpanel

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/pkg/panel"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// BoolField is a helper for modifying bool fields.
func BoolField(app AppPreview, fieldName string) panel.CallbackFunc {
	return func(ctx *panel.Context) (string, [][]gotgbot.InlineKeyboardButton, error) {
		var (
			op   string
			data = ctx.CallbackData
		)

		if len(data.Args) != 0 {
			op = data.Args[0]
		}

		var s string

		switch op {
		case OperationSet:
			err := app.GetDB().UpdateConfig(ctx.Bot.Id, fieldName, true)
			if err != nil {
				return "", nil, err
			}

			s = fmt.Sprintf("<i><b>✅ %s has been Enabled !</b></i>", ctx.Page.DisplayName)
		case OperationReset:
			err := app.GetDB().ResetConfig(ctx.Bot.Id, fieldName)
			if err != nil {
				return "", nil, err
			}

			s = fmt.Sprintf("<i><b>✅ %s has been Reset !</b></i>", ctx.Page.DisplayName)
		default:
			return fmt.Sprintf("<i>Use The Buttons Below to Enable/Disable %s</i>", ctx.Page.DisplayName),
				[][]gotgbot.InlineKeyboardButton{{{Text: "Enable", CallbackData: data.RemoveArgs().AddArg(OperationSet).ToString()}, {Text: "Disable", CallbackData: data.RemoveArgs().AddArg(OperationReset).ToString()}}},
				nil
		}

		go app.RefreshConfig() // is a goroutine a bit overkill here

		return s, nil, nil
	}
}
