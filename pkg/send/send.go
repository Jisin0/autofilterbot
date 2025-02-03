// Package send simplifies sending messages of various types.
package send

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

// SendOpts are options for sendMethod funcs.
type SendOpts struct {
	Text     string
	Keyboard [][]gotgbot.InlineKeyboardButton
	FileId   string
}

// SendMethod is type to be implemented by any sendMessage functions.
type SendMethod func(bot *gotgbot.Bot, chatId int64, opts *SendOpts) (*gotgbot.Message, error)

// Ensure SendDocument implements SendMethod.
var _ SendMethod = SendDocument

// SendDocument is helper for sending a document.
func SendDocument(bot *gotgbot.Bot, chatId int64, opts *SendOpts) (*gotgbot.Message, error) {
	if opts == nil {
		return nil, errors.New("options is nil")
	}

	if opts.FileId == "" {
		return nil, errors.New("file_id is required")
	}

	return bot.SendDocument(chatId, gotgbot.InputFileByID(opts.FileId), &gotgbot.SendDocumentOpts{
		Caption:   opts.Text,
		ParseMode: gotgbot.ParseModeHTML,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: opts.Keyboard,
		},
	})
}

// Ensure SendVideo implements SendMethod.
var _ SendMethod = SendVideo

// SendVideo is helper for sending a video.
func SendVideo(bot *gotgbot.Bot, chatId int64, opts *SendOpts) (*gotgbot.Message, error) {
	if opts == nil {
		return nil, errors.New("options is nil")
	}

	if opts.FileId == "" {
		return nil, errors.New("file_id is required")
	}

	return bot.SendVideo(chatId, gotgbot.InputFileByID(opts.FileId), &gotgbot.SendVideoOpts{
		Caption:   opts.Text,
		ParseMode: gotgbot.ParseModeHTML,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: opts.Keyboard,
		},
	})
}

// Ensure SendAudio implements SendMethod.
var _ SendMethod = SendAudio

// SendAudio is helper for sending an audio.
func SendAudio(bot *gotgbot.Bot, chatId int64, opts *SendOpts) (*gotgbot.Message, error) {
	if opts == nil {
		return nil, errors.New("options is nil")
	}

	if opts.FileId == "" {
		return nil, errors.New("file_id is required")
	}

	return bot.SendAudio(chatId, gotgbot.InputFileByID(opts.FileId), &gotgbot.SendAudioOpts{
		Caption:   opts.Text,
		ParseMode: gotgbot.ParseModeHTML,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: opts.Keyboard,
		},
	})
}

// Ensure SendPhoto implements SendMethod.
var _ SendMethod = SendPhoto

// SendPhoto is helper for sending a photo.
func SendPhoto(bot *gotgbot.Bot, chatId int64, opts *SendOpts) (*gotgbot.Message, error) {
	if opts == nil {
		return nil, errors.New("options is nil")
	}

	if opts.FileId == "" {
		return nil, errors.New("file_id is required")
	}

	return bot.SendPhoto(chatId, gotgbot.InputFileByID(opts.FileId), &gotgbot.SendPhotoOpts{
		Caption:   opts.Text,
		ParseMode: gotgbot.ParseModeHTML,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: opts.Keyboard,
		},
	})
}

// Ensure SendAnimation implements SendMethod.
var _ SendMethod = SendAnimation

// SendAnimation is helper for sending an animation.
func SendAnimation(bot *gotgbot.Bot, chatId int64, opts *SendOpts) (*gotgbot.Message, error) {
	if opts == nil {
		return nil, errors.New("options is nil")
	}

	if opts.FileId == "" {
		return nil, errors.New("file_id is required")
	}

	return bot.SendAnimation(chatId, gotgbot.InputFileByID(opts.FileId), &gotgbot.SendAnimationOpts{
		Caption:   opts.Text,
		ParseMode: gotgbot.ParseModeHTML,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: opts.Keyboard,
		},
	})
}

// Ensure SendMessage implements SendMethod.
var _ SendMethod = SendMessage

// sendMessgae is helper for sending a text message.
func SendMessage(bot *gotgbot.Bot, chatId int64, opts *SendOpts) (*gotgbot.Message, error) {
	if opts == nil {
		return nil, errors.New("options is nil")
	}

	if opts.Text == "" {
		return nil, errors.New("text is required")
	}

	return bot.SendMessage(chatId, opts.Text, &gotgbot.SendMessageOpts{
		ParseMode: gotgbot.ParseModeHTML,
		ReplyMarkup: gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: opts.Keyboard,
		},
	})
}
