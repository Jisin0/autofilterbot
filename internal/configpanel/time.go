package configpanel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jisin0/autofilterbot/pkg/panel"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

const (
	buttonsPerRow = 3
)

// TimeField is a helper for modifying time fields stored in minutes as an integer.
func TimeField(app AppPreview, fieldName string, possibleValues []int) panel.CallbackFunc {
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
			if len(data.Args) < 2 {
				app.GetLog().Warn("panel: time: not enough args for set", zap.Strings("args", data.Args))
				return s, nil, fmt.Errorf("not enough argmuents")
			}

			val, err := strconv.Atoi(data.Args[1])
			if err != nil {
				return "", nil, err
			}

			if !containsInt(possibleValues, val) {
				return "", nil, fmt.Errorf("unknown value %d received for field %s", val, fieldName)
			}

			err = app.GetDB().UpdateConfig(ctx.Bot.Id, fieldName, val)
			if err != nil {
				return "", nil, err
			}

			s = fmt.Sprintf("<i><b>✅ %s has been set to %d minutes !</b></i>", ctx.Page.DisplayName, val)
		case OperationReset:
			err := app.GetDB().ResetConfig(ctx.Bot.Id, fieldName)
			if err != nil {
				return "", nil, err
			}

			s = fmt.Sprintf("<i><b>✅ %s has been Reset !</b></i>", ctx.Page.DisplayName)
		default:
			var s strings.Builder

			s.WriteString(fmt.Sprintf("<i><b>🧮 Select One of the Values Below to Update %s to Given Number in Minutes\n\n", ctx.Page.DisplayName))

			if v, ok := app.GetConfig().ToMap()[fieldName]; ok {
				if i, ok := v.(int); ok && i != 0 {
					s.WriteString(fmt.Sprintf("⭕ Current Value: %dMins\n\n", i))
				}
			}

			s.WriteString("🔁 Use the Reset Option to Set it to The Default Value:</b></i>")

			keyboard := make([][]gotgbot.InlineKeyboardButton, 0, len(possibleValues)/2)

			for i := 0; i < len(possibleValues); i += buttonsPerRow {
				end := i + buttonsPerRow
				if end > len(possibleValues) {
					end = len(possibleValues)
				}

				row := make([]gotgbot.InlineKeyboardButton, 0, end-i)

				for _, v := range possibleValues[i:end] {
					valStr := strconv.Itoa(v)
					row = append(row, gotgbot.InlineKeyboardButton{
						Text:         valStr,
						CallbackData: ctx.CallbackData.AddArg(OperationSet).AddArg(valStr).ToString(),
					})
				}
				keyboard = append(keyboard, row)
			}

			keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{{Text: "ʀᴇsᴇᴛ 🔁", CallbackData: ctx.CallbackData.AddArg(OperationReset).ToString()}})

			return s.String(), keyboard, nil
		}

		go app.RefreshConfig()

		return s, nil, nil
	}
}

func containsInt(s []int, val int) bool {
	for _, n := range s {
		if n == val {
			return true
		}
	}

	return false
}
