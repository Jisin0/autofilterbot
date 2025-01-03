package panel

import (
	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Context wraps the update and other additional data for callback functions to use.
type Context struct {
	// Bot object.
	Bot *gotgbot.Bot
	// Full Update.
	Update *ext.Context
	// Query which propogated the request.
	CallbackQuery *gotgbot.CallbackQuery
	// CallbackData wraps the request path and arguments.
	CallbackData callbackdata.CallbackData
}
