package exthandlers

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

// NewAllUpdates creates a handler that accepts all updates.
func NewAllUpdates(r handlers.Response) ext.Handler {
	return AllUpdates{response: r}
}

// AllUpdates handles all updates.
type AllUpdates struct {
	response handlers.Response
}

func (_ AllUpdates) CheckUpdate(bot *gotgbot.Bot, ctx *ext.Context) bool {
	return true
}

func (h AllUpdates) HandleUpdate(bot *gotgbot.Bot, ctx *ext.Context) error {
	return h.response(bot, ctx)
}
func (_ AllUpdates) Name() string {
	return "All Updates"
}
