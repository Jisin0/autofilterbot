package core

import (
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// HandleJoinRequest handles join requests to force subsribe channels.
func HandleJoinRequest(bot *gotgbot.Bot, ctx *ext.Context) error {
	update := ctx.ChatJoinRequest

	// saves all join requests, not just from fsub channels

	err := _app.DB.SaveUserJoinRequest(update.From.Id, update.Chat.Id)
	if err != nil {
		_app.Log.Warn("handlejoinrequest: failed to save join request to db", zap.Int64("user_id", update.From.Id), zap.Int64("chat_id", update.Chat.Id), zap.Error(err))
	} else {
		_app.Log.Debug("handlejoinrequest: saved join request sucessfully", zap.Int64("user_id", update.From.Id), zap.Int64("chat_id", update.Chat.Id))
	}

	return nil
}
