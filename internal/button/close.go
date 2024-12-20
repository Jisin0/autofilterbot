package button

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/pkg/configpanel/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Close creates a close button which allows an authorized user to delete the message to which it's attached.
func Close(userId ...int64) gotgbot.InlineKeyboardButton {
	data := callbackdata.New().AddPath("close")
	for _, u := range userId {
		data.AddArg(fmt.Sprint(u))
	}

	return gotgbot.InlineKeyboardButton{
		Text:         "X ᴄʟᴏsᴇ X",
		CallbackData: data.ToString(),
	}
}
