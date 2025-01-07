/*
Package autofilter contains types and methods to help work with autofilter results.
*/
package autofilter

import (
	"context"
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/internal/format"
	"github.com/Jisin0/autofilterbot/internal/functions"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/pkg/shortener"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Files []model.File

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
			result = append(result, []gotgbot.InlineKeyboardButton{{Text: f.FileName, CallbackData: "fdetails|" + f.FileId}, {Text: size, Url: url}})
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

// ProcessQuery handles the query after sanitization and validation and returns the results.
func ProcessQuery(bot *gotgbot.Bot, ctx *ext.Context, db database.Database, query string) (*SearchResult, error) {

	return nil, nil
}

type FilesFromCursorOptions interface {
	GetMaxResults() int
	GetMaxPages() int
	GetMaxPerPage() int
}

// FilesFromCursor loops through a cursor and outputs an array of files.
func FilesFromCursor(ctx context.Context, c database.Cursor, opts FilesFromCursorOptions) ([]Files, error) {
	var (
		totalCount int
		finished   bool
		totalFiles = make([]Files, 0, opts.GetMaxResults())
	)

	for i := 0; i < opts.GetMaxPages(); i++ {
		row := make([]model.File, 0, opts.GetMaxPerPage())

		fmt.Println(i) //TODO: remove

		for j := 0; j < opts.GetMaxPerPage(); j++ {
			fmt.Println(j)
			if !c.Next(ctx) {
				finished = true
				break
			}

			var f model.File

			err := c.Decode(&f)
			if err != nil {
				return totalFiles, err
			}

			row = append(row, f)

			if finished {
				return totalFiles, err
			}
		}

		totalFiles = append(totalFiles, row)

		if totalCount >= opts.GetMaxResults() { // only checks after completing a page. should this behaviour be changed?
			break
		}
	}

	return totalFiles, nil
}
