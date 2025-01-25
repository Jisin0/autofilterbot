package model

import (
	"fmt"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

const (
	FileTypeDocument = "document"
	FileTypeVideo    = "video"
	FileTypeAudio    = "audio"
	FileTypeVoice    = "voice"
)

// File is a single file stored in the database.
type File struct {
	// Unique id of the file used to save and fetch files from db.
	UniqueId string `json:"_id" bson:"_id"`
	// Id used to send and copy files
	FileId string `json:"file_id" bson:"file_id"`
	// Name of the file including extension
	FileName string `json:"file_name" bson:"file_name"`
	// Type of file either document/video/audio.
	FileType string `json:"file_type" bson:"file_type"`
	// Size of the file in bytes.
	FileSize int64 `json:"file_size" bson:"file_size"`
	// Unix timestamp of time when file was saved.
	Time int64 `json:"time,omitempty" bson:"time,omitempty"`
	// Id of the chat/channel where the file is posted.
	ChatId int64 `json:"chat_id,omitempty" bson:"chat_id,omitempty"`
	// Link to the original message containing the file.
	MessageLink string `json:"file_link,omitempty" bson:"file_link,omitempty"`
}

type SendFileOpts struct {
	Caption  string
	Keyboard [][]gotgbot.InlineKeyboardButton
}

// Send sends the file to chatId with given caption, markup and html parse mode.
func (f *File) Send(bot *gotgbot.Bot, chatId int64, opts *SendFileOpts) (*gotgbot.Message, error) {
	switch f.FileType {
	case FileTypeDocument:
		sendOpts := &gotgbot.SendDocumentOpts{ParseMode: gotgbot.ParseModeHTML}

		if opts != nil {
			if opts.Caption != "" {
				sendOpts.Caption = opts.Caption
			}

			if len(opts.Keyboard) != 0 {
				sendOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{InlineKeyboard: opts.Keyboard}
			}
		}

		return bot.SendDocument(chatId, gotgbot.InputFileByID(f.FileId), sendOpts)
	case FileTypeVideo:
		sendOpts := &gotgbot.SendVideoOpts{ParseMode: gotgbot.ParseModeHTML}

		if opts != nil {
			if opts.Caption != "" {
				sendOpts.Caption = opts.Caption
			}

			if len(opts.Keyboard) != 0 {
				sendOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{InlineKeyboard: opts.Keyboard}
			}
		}

		return bot.SendVideo(chatId, gotgbot.InputFileByID(f.FileId), sendOpts)
	case FileTypeAudio:
		sendOpts := &gotgbot.SendAudioOpts{ParseMode: gotgbot.ParseModeHTML}

		if opts != nil {
			if opts.Caption != "" {
				sendOpts.Caption = opts.Caption
			}

			if len(opts.Keyboard) != 0 {
				sendOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{InlineKeyboard: opts.Keyboard}
			}
		}

		return bot.SendAudio(chatId, gotgbot.InputFileByID(f.FileId), sendOpts)
	case FileTypeVoice:
		sendOpts := &gotgbot.SendVoiceOpts{ParseMode: gotgbot.ParseModeHTML}

		if opts != nil {
			if opts.Caption != "" {
				sendOpts.Caption = opts.Caption
			}

			if len(opts.Keyboard) != 0 {
				sendOpts.ReplyMarkup = gotgbot.InlineKeyboardMarkup{InlineKeyboard: opts.Keyboard}
			}
		}

		return bot.SendVoice(chatId, gotgbot.InputFileByID(f.FileId), sendOpts)
	default:
		return nil, fmt.Errorf("unsupported file type %s", f.FileType)
	}
}
