package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Jisin0/autofilterbot/internal/autofilter"
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
		_app.Log.Debug("newfile: empty file name after sanitization", zap.String("file_id", file.FileName), zap.String("file_type", file.FileType))
		return nil
	}

	err := _app.DB.SaveFile(file)
	if err != nil {
		if _, ok := err.(database.FileAlreadyExistsError); ok {
			_app.Log.Debug("newfile: duplicate file skipped", zap.String("file_name", file.FileName))
			return nil
		}

		_app.Log.Warn("newfile: failed to save file", zap.Error(err))
	}

	return nil
}

// DeleteFile handles the /delete command to delete a file.
func DeleteFile(bot *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	if !_app.AuthAdmin(ctx) {
		return nil
	}

	var fileUniqueId string

	if f := functions.FileFromMessage(m.ReplyToMessage); f != nil {
		fileUniqueId = f.UniqueId
	} else {
		conv := conversation.NewConversatorFromUpdate(bot, ctx.Update)

		replyM, err := conv.Ask(_app.Ctx, "Please send me the file you would like to delete:", nil)
		if err != nil {
			m.Reply(bot, fmt.Sprintf("An Error occurred: %v", err), nil)
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

const (
	delAllCountDangerous = 20 // if more than this many files are to be deleted, user must be prompted for confirmation
)

// DeleteAllFiles handles the /deleteall command to delete all matching files.
func DeleteAllFiles(bot *gotgbot.Bot, ctx *ext.Context) error {
	m := ctx.EffectiveMessage

	if !_app.AuthAdmin(ctx) {
		return nil
	}

	split := strings.SplitN(m.Text, " ", 2)
	if len(split) != 2 {
		m.Reply(bot, "<b>Improper Usage, Keyword Missing!</b>\n<blockquote>Format:\n /deleteall &lt keyword&gt</blockquote>\n<blockquote>Example:\n /deleteall 720p</blockquote>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		return nil
	}

	keyword := split[1]
	if keyword == "" { // shouldnt happen ideally
		m.Reply(bot, "Keyword Is Empty :/", nil)
		return nil
	}

	cursor, err := _app.DB.SearchFiles(keyword)
	if err != nil {
		m.Reply(bot, fmt.Sprintf("An Error occurred: %v", err), nil)
		_app.Log.Warn("delall: search files failed", zap.Error(err), zap.String("keyword", keyword))
		return nil
	}

	f, err := autofilter.FilesFromCursor(context.Background(), cursor, DeleteAllCursorOptions{})
	if err != nil {
		m.Reply(bot, fmt.Sprintf("An Error occurred: %v", err), nil)
		_app.Log.Warn("delall: files from cursor failed", zap.Error(err), zap.String("keyword", keyword))
		return nil
	}

	if len(f) == 0 || len(f[0]) == 0 {
		m.Reply(bot, fmt.Sprintf("I Couldn't Find Anything Matching %s", keyword), nil)
		return nil
	}

	files := f[0]
	if len(files) > delAllCountDangerous {
		conv := conversation.NewConversatorFromUpdate(bot, ctx.Update)

		msg, err := conv.Ask(_app.Ctx, fmt.Sprintf("<b>⛔ Dangerous Operation</b>\n<i>Are you sure you want to delete %d files? Send <code>yes</code> to confirm or /cancel to cancel:</i>", len(files)), &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
		if err != nil {
			_app.Log.Warn("delall: ask confirmation failed", zap.Error(err))
			return nil
		}

		if !strings.EqualFold(msg.Text, "yes") {
			msg.Reply(bot, "<i>Operation Cancelled !</i>", &gotgbot.SendMessageOpts{ParseMode: gotgbot.ParseModeHTML})
			return nil
		}
	}

	progM, err := m.Reply(bot, "🗑️ Deleting Files, Please Wait...", nil)
	if err != nil {
		_app.Log.Warn("delall: send progress msg failed", zap.Error(err))
		return nil
	}

	var (
		allErrors []error
		deleted   int
	)

	for _, file := range files {
		err := _app.DB.DeleteFile(file.UniqueId)
		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}

		deleted += 1
	}

	text := fmt.Sprintf("<i><b>✅ Deleted %d Files Successfully !</b></i>", deleted)

	if len(allErrors) != 0 {
		errs := errors.Join(allErrors...)
		text += fmt.Sprintf("\nErrors occurred: %v", errs)
		_app.Log.Info("delall: errors occurs", zap.Error(errs))
	}

	_, _, err = progM.EditText(bot, text, &gotgbot.EditMessageTextOpts{ParseMode: gotgbot.ParseModeHTML})
	if err != nil {
		_app.Log.Warn("delall: update progress msg failed", zap.Error(err))
	}

	return nil
}

// DeleteAllCursorOptions implements autofilter.FilesFromCursorOptions with optimal options for the operation.
type DeleteAllCursorOptions struct{}

// ensure it implements autofilter.FilesFromCursorOptions
var _ autofilter.FilesFromCursorOptions = DeleteAllCursorOptions{}

func (DeleteAllCursorOptions) GetMaxResults() int {
	return 100 // max 100 files can be deleted in a go
}

func (DeleteAllCursorOptions) GetMaxPages() int {
	return 1 // all files are in first page
}

func (DeleteAllCursorOptions) GetMaxPerPage() int {
	return 200 // just has to be more than max results
}
