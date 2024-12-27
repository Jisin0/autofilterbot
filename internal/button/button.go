package button

import "github.com/PaulSonOfLars/gotgbot/v2"

// InlineKeyboardButton wraps gotgbot.InlineKeyboardButton with bson struct tags to add omitempty fields for optional fields.
type InlineKeyboardButton struct {
	Text         string `json:"text" bson:"text"`
	CallbackData string `json:"callback_data,omitempty" bson:"callback_data,omitempty"`
	Url          string `json:"url,omitempty" bson:"url,omitempty"`
	// Extra Custom field to remove pointer from value field.
	IsInline                     bool   `json:"is_empty,omitempty" bson:"is_empty,omitempty"`
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty" bson:"switch_inline_query_current_chat,omitempty"`
	CopyText                     string `json:"copy_text,omitempty" bson:"copy_text,omitempty"`
}

func NewInlineKeyboardButton(val gotgbot.InlineKeyboardButton) InlineKeyboardButton {
	b := InlineKeyboardButton{
		Text: val.Text,
	}

	switch {
	case val.CallbackData != "":
		b.CallbackData = val.CallbackData
	case val.Url != "":
		b.Url = val.Url
	case val.SwitchInlineQueryCurrentChat != nil:
		b.SwitchInlineQueryCurrentChat = *val.SwitchInlineQueryCurrentChat
	case val.CopyText != nil:
		b.CopyText = val.CopyText.Text
	}

	return b
}

// Unwrap chanages InlineKeyboardButton to gotgbot.InlineKeyboardButton.
func (val InlineKeyboardButton) Unwwrap() gotgbot.InlineKeyboardButton {
	b := gotgbot.InlineKeyboardButton{
		Text: val.Text,
	}

	switch {
	case val.CallbackData != "":
		b.CallbackData = val.CallbackData
	case val.Url != "":
		b.Url = val.Url
	case val.IsInline:
		b.SwitchInlineQueryCurrentChat = &val.SwitchInlineQueryCurrentChat
	case val.CopyText != "":
		b.CopyText = &gotgbot.CopyTextButton{Text: val.CopyText}
	}

	return b
}
