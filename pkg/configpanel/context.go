package configpanel

import (
	"github.com/Jisin0/autofilterbot/internal/app"
	"github.com/Jisin0/autofilterbot/pkg/configpanel/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

// Context wraps the app, update and other additional data for callback functions to use.
type Context struct {
	// Application.
	App *app.App
	// Bot object.
	Bot *gotgbot.Bot
	// Full Update.
	Update *ext.Context
	// Query which propogated the request.
	CallbackQuery *gotgbot.CallbackQuery
	// CallbackData wraps the request path and arguments.
	CallbackData callbackdata.CallbackData
}
