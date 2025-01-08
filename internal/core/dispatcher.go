package core

import (
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/Jisin0/autofilterbot/pkg/conversation"
	"github.com/Jisin0/autofilterbot/pkg/env"
	exthandlers "github.com/Jisin0/autofilterbot/pkg/filters"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/callbackquery"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
	"go.uber.org/zap"
)

const (
	autofilterHandlerGroup = iota + 1
	commandHandlerGroup
	callbackQueryGroup
	miscHandlerGroup
	middleWareGroup
)

// SetupDispatcher creates a new empty dispatcher with error and panic recovery setup.
func SetupDispatcher(log *zap.Logger) *ext.Dispatcher {
	d := ext.NewDispatcher(&ext.DispatcherOpts{
		// If an error is returned by a handler, log it and continue going.
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			logFields := []zap.Field{zap.Error(err)}
			logFields = addLogFieldsFromContext(ctx, logFields)

			log.Error("error occurred while handling update", logFields...)

			return ext.DispatcherActionNoop
		},

		Panic: func(b *gotgbot.Bot, ctx *ext.Context, r interface{}) {
			logFields := []zap.Field{zap.String("panic", fmt.Sprintf("%v\n%s", r, cleanedStack()))}
			logFields = addLogFieldsFromContext(ctx, logFields)

			log.Error("panic recovered", logFields...)
		},
	})

	d.AddHandlerToGroup(handlers.NewMessage(message.Supergroup, Autofilter), autofilterHandlerGroup)

	d.AddHandlerToGroup(exthandlers.NewCommands([]string{"start", "about", "help", "privacy"}, StaticCommands), commandHandlerGroup)
	d.AddHandlerToGroup(handlers.NewCommand("delete", DeleteFile).SetAllowChannel(true), commandHandlerGroup)
	d.AddHandlerToGroup(handlers.NewCommand("logs", Logs), commandHandlerGroup)

	d.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("cmd"), StaticCommands), callbackQueryGroup)
	d.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("close"), Close), callbackQueryGroup)
	d.AddHandlerToGroup(handlers.NewCallback(callbackquery.Prefix("navg"), Navigate), callbackQueryGroup)

	d.AddHandlerToGroup(handlers.NewMessage(exthandlers.ChatIds(env.Int64s("FILE_CHANNELS")), NewFile), miscHandlerGroup)

	d.AddHandlerToGroup(handlers.NewMessage(message.All, conversation.MessageHandler), middleWareGroup)

	return d
}

// logFieldsFromContext adds zap fields to logFields about specific update from ctx to help troubleshooting.
func addLogFieldsFromContext(ctx *ext.Context, logFields []zap.Field) []zap.Field {
	switch {
	case ctx.Message != nil:
		logFields = append(logFields,
			zap.Int64("chat_id", ctx.Message.Chat.Id),
			zap.Int64("message_id", ctx.Message.MessageId),
			zap.String("text", ctx.Message.Text),
		)
	case ctx.CallbackQuery != nil:
		logFields = append(logFields,
			zap.String("callback_query_id", ctx.CallbackQuery.Id),
			zap.String("data", ctx.CallbackQuery.Data),
			zap.Int64("chat_id", ctx.CallbackQuery.Message.GetChat().Id),
			zap.Int64("message_id", ctx.CallbackQuery.Message.GetMessageId()),
		)
	case ctx.InlineQuery != nil:
		logFields = append(logFields,
			zap.String("inline_query_id", ctx.InlineQuery.Id),
			zap.String("query", ctx.InlineQuery.Query),
			zap.Int64("from", ctx.InlineQuery.From.Id),
		)
	}

	return logFields
}

// cleanedStack returns stack trace with gotgbot library parts removed to prevent confusion.
// Copied from https://github.com/PaulSonOfLars/gotgbot/blob/v2/ext/dispatcher.go.
func cleanedStack() string {
	return strings.Join(strings.Split(string(debug.Stack()), "\n")[4:], "\n")
}
