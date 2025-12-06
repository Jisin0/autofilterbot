package configpanel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Jisin0/autofilterbot/pkg/panel"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

// IntFieldOpts wraps optional values to IntField().
type IntFieldOpts struct {
	// Provide a range the value should fall within.
	Range *IntRange
	// Check within a list of possible values.
	PossibleValues []int
	// Description for the field.
	Description string
	// Middleware is called after a value is successully set, allowing any sync operations to be run.
	Middleware func(val int)
}

// IntField is a helper for cofiguring int values in the config panel.
func IntField(app AppPreview, fieldName string, opts IntFieldOpts) panel.CallbackFunc {
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
				app.GetLog().Warn("configpanel: int: not enough args for set", zap.Strings("args", data.Args))
				return s, nil, fmt.Errorf("not enough argmuents")
			}

			val, err := strconv.Atoi(data.Args[1])
			if err != nil {
				return "", nil, err
			}

			var isValid bool

			switch {
			case opts.Range != nil:
				isValid = opts.Range.Check(val)
			case len(opts.PossibleValues) != 0:
				isValid = containsInt(opts.PossibleValues, val)
			default:
				app.GetLog().Debug("configpanel: int: no range or possible values provided", zap.String("field", fieldName))
				isValid = true
			}

			if !isValid {
				return "", nil, fmt.Errorf("unknown value %d received for field %s", val, fieldName)
			}

			err = app.GetDB().UpdateConfig(ctx.Bot.Id, fieldName, val)
			if err != nil {
				return "", nil, err
			}

			if opts.Middleware != nil {
				opts.Middleware(val)
			}

			s = fmt.Sprintf("<i><b>✅ %s has been set to %d!</b></i>", ctx.Page.DisplayName, val)
		case OperationReset:
			err := app.GetDB().ResetConfig(ctx.Bot.Id, fieldName)
			if err != nil {
				return "", nil, err
			}

			s = fmt.Sprintf("<i><b>✅ %s has been Reset !</b></i>", ctx.Page.DisplayName)
		default:
			var s strings.Builder

			if opts.Description != "" {
				s.WriteString(fmt.Sprintf("ℹ️ %s\n\n", opts.Description))
			}

			s.WriteString(fmt.Sprintf("<i><b>🧮 Select One of the Values Below to Update %s to Given Number\n\n", ctx.Page.DisplayName))

			if v, ok := app.GetConfig().ToMap()[fieldName]; ok {
				if i, ok := v.(int); ok && i != 0 {
					s.WriteString(fmt.Sprintf("⭕ Current Value: %d\n\n", i))
				}
			}

			s.WriteString("🔁 Use the Reset Option to Set it to The Default Value:</b></i>")

			values := opts.PossibleValues
			if opts.Range != nil {
				values = opts.Range.Slice()
			}

			keyboard := make([][]gotgbot.InlineKeyboardButton, 0, len(values)/2)

			for i := 0; i < len(values); i += buttonsPerRow {
				end := i + buttonsPerRow
				if end > len(values) {
					end = len(values)
				}

				row := make([]gotgbot.InlineKeyboardButton, 0, end-i)

				for _, v := range values[i:end] {
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

// IntRange specifies a range to check values.
type IntRange struct {
	// Start of the range (inclusive).
	Start int
	// End of the range (inclusive).
	End int
	// Values within provided range that should be excluded.
	ExcludedValues []int
}

// Check reports wether the given value falls within the range.
func (r IntRange) Check(n int) bool {
	return n >= r.Start && n <= r.End && !containsInt(r.ExcludedValues, n)
}

// Slice returns the possible values as a slice.
func (r IntRange) Slice() []int {
	l := make([]int, 0)

	for i := r.Start; i <= r.End; i++ {
		l = append(l, i)
	}

	return l
}
