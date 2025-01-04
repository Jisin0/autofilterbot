package button

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Close creates a close button which allows an authorized user to delete the message to which it's attached.
func Close(userId ...int64) gotgbot.InlineKeyboardButton {
	data := callbackdata.New().AddPath("close")
	for _, u := range userId {
		data.AddArg(fmt.Sprint(u))
	}

	return gotgbot.InlineKeyboardButton{
		Text:         "ᴄʟᴏsᴇ ⛌",
		CallbackData: data.ToString(),
	}
}

// CloseLocal creates a close button of the local InlineKeyboardButton type.
func CloseLocal(userId ...int64) InlineKeyboardButton {
	data := callbackdata.New().AddPath("close")
	for _, u := range userId {
		data.AddArg(fmt.Sprint(u))
	}

	return InlineKeyboardButton{
		Text:         "ᴄʟᴏsᴇ ⛌",
		CallbackData: data.ToString(),
	}
}
