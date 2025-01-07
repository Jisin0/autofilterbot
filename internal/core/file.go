package core

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/pkg/conversation"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

// NewFile handles message updates in any authorized file channels.
func NewFile(bot *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	if !functions.HasMedia(m) {
		return nil
	}

	file := functions.FileFromMessage(m)
	if file == nil {
		return nil
	}

	if file.FileName == "" {
		_app.Log.Debug("empty file name after sanitization", zap.String("file_id", file.FileName), zap.String("file_type", file.FileType))
		return nil
	}

	err := _app.DB.SaveFile(file)
	if err != nil {
		if _, ok := err.(database.FileAlreadyExistsError); ok {
			_app.Log.Debug("duplicate file skipped", zap.String("file_name", file.FileName))
			return nil
		}

		_app.Log.Warn("failed to save file", zap.Error(err))
	}

	return nil
}

// DeleteFile handles the /delete command to delete a file.
func DeleteFile(bot *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	if !_app.AuthAdmin(m) {
		return nil
	}

	var fileUniqueId string

	if f := functions.FileFromMessage(m.ReplyToMessage); f != nil {
		fileUniqueId = f.UniqueId
	} else {
		conv := conversation.NewConversatorFromUpdate(bot, ctx.Update)

		replyM, err := conv.Ask("Please send me the file you would like to delete:", nil)
		if err != nil {
			m.Reply(bot, fmt.Sprintf("An Error Occured: %v", err), nil)
			return nil
		}

		f := functions.FileFromMessage(replyM)
		if f == nil {
			replyM.Reply(bot, "No Media Found!", nil)
			return nil
		}

		fileUniqueId = f.UniqueId
	}

	err := _app.DB.DeleteFile(fileUniqueId)
	if err != nil {
		m.Reply(bot, fmt.Sprintf("Failed to Delete File: %v", err), nil)
		_app.Log.Warn("delete file failed", zap.String("file_id", fileUniqueId), zap.Error(err))
		return nil
	}

	m.Reply(bot, "𝖥𝗂𝗅𝖾 𝖶𝖺𝗌 𝖣𝖾𝗅𝖾𝗍𝖾𝖽 𝖲𝗎𝖼𝖼𝖾𝗌𝗌𝖿𝗎𝗅𝗅𝗒 🗑️", nil)

	return nil
}
