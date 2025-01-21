package autofilter

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/format"
	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/pkg/shortener"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type Files []File

// File wraps model.File with some search result specific data.
type File struct {
	model.File
	// Indicates whether the file was selected from the selection menu.
	IsSelected bool `json:"selected,omitempty"`
}

// Process returns a slice of buttons to be used in message markup.
func (files Files) Process(chatId int64, botUsername string, opts ProcessFilesOptions) [][]gotgbot.InlineKeyboardButton {
	return ProcessFiles(files, chatId, botUsername, opts)
}

type ProcessFilesOptions interface {
	GetButtonTemplate() string
	GetSizeButton() bool
	GetShortener() shortener.Shortener
}

// ProcessFiles changes files into a keboard slice to be used as markup in a message.
func ProcessFiles(files Files, chatId int64, botUsername string, opts ProcessFilesOptions) [][]gotgbot.InlineKeyboardButton {
	var (
		hasShortener = opts.GetShortener().ApiKey != ""
		result       = make([][]gotgbot.InlineKeyboardButton, 0, len(files))
	)

	for _, f := range files {
		url := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, URLData{
			FileUniqueId: f.UniqueId,
			ChatId:       chatId,
			HasShortener: hasShortener,
		}.Encode())
		size := functions.FileSizeToString(f.FileSize)

		if opts.GetSizeButton() {
			result = append(result, []gotgbot.InlineKeyboardButton{{Text: f.FileName, CallbackData: "fdetails|" + f.UniqueId}, {Text: size, Url: url}})
		} else {
			text := format.KeyValueFormat(opts.GetButtonTemplate(), map[string]string{
				"file_name": f.FileName,
				"file_size": size,
			})
			result = append(result, []gotgbot.InlineKeyboardButton{{Text: text, Url: url}})
		}
	}

	return result
}

// SelectMenu returns a keyboard with to select files from.
func (files Files) SelectMenu(uniqueId string, pageIndex int) [][]gotgbot.InlineKeyboardButton {
	keyboard := make([][]gotgbot.InlineKeyboardButton, 0, len(files))

	for _, f := range files {
		keyboard = append(keyboard, []gotgbot.InlineKeyboardButton{{
			Text:         fmt.Sprintf("%s[%s] %s", tick(f.IsSelected), functions.FileSizeToString(f.FileSize), f.FileName),
			CallbackData: fmt.Sprintf("sel|%s_%d_%s", uniqueId, pageIndex, f.UniqueId),
		}})
	}

	return keyboard
}

// tick returns a tick symbol if val is true.
func tick(val bool) string {
	if val {
		return "✅ "
	}

	return ""
}
