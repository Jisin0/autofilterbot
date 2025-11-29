package model

import (
	"github.com/Jisin0/autofilterbot/pkg/callbackdata"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

// Index is data about an index operation.
type Index struct {
	// Unique id of operation.
	ID string `json:"_id" bson:"_id"`
	// Id of message to start with.
	StartMessageID int64 `json:"start" bson:"start"`
	// Id of last message to index.
	EndMessageID int64 `json:"end" bson:"end"`
	// Id of message currently being indexed. Start id incase of a message range.
	CurrentMessageID int64 `json:"current" bson:"current"`
	// Channel id in bot api format.
	ChannelID int64 `json:"channel" bson:"channel"`
	// Number of files successfully saved.
	Saved int `json:"saved,omitempty" bson:"saved,omitempty"`
	// Messages failed to save.
	Failed int `json:"failed,omitempty" bson:"failed,omitempty"`
	// Indicates whether the operation is paused.
	IsPaused bool `json:"is_paused,omitempty" bson:"is_paused,omitempty"`

	// Id of chat where the index operation was started
	ProgressMessageChatID int64 `json:"pmessage_chat,omitempty" bson:"pmessage_chat,omitempty"`
	//DEPRECATED: Id of the progress message
	// ProgressMessageID int64 `json:"pmessage_id,omitempty" bson:"pmessage_id,omitempty"`
}

const (
	IndexCharStart  = "s"
	IndexCharPause  = "p"
	IndexCharCancel = "c"
	IndexCharModify = "m"
)

// PauseButton returns a keyboard button that can be used to pase the operation.
func (o *Index) PauseButton() gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         "Pause ⏹️",
		CallbackData: callbackdata.New().AddPath("index").AddArg(o.ID).AddArg(IndexCharPause).ToString(),
	}
}

// StartButton returns a keyboard button that can be used to start the operation.
func (o *Index) StartButton() gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         "Start ⚡",
		CallbackData: callbackdata.New().AddPath("index").AddArg(o.ID).AddArg(IndexCharStart).ToString(),
	}
}

// ResumeButton does the same as the start button but sets the text to resume.
func (o *Index) ResumeButton() gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         "Resume ⏸️",
		CallbackData: callbackdata.New().AddPath("index").AddArg(o.ID).AddArg(IndexCharStart).ToString(),
	}
}

// CancelButton returns a button that aborts the index operation, completely erasing it.
func (o *Index) CancelButton() gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         "Cancel ❌",
		CallbackData: callbackdata.New().AddPath("index").AddArg(o.ID).AddArg(IndexCharCancel).ToString(),
	}
}

// ModifyButton returns a button that opens the panel to change index configuration.
func (o *Index) ModifyButton() gotgbot.InlineKeyboardButton {
	return gotgbot.InlineKeyboardButton{
		Text:         "Modify ⚙️",
		CallbackData: callbackdata.New().AddPath("index").AddArg(o.ID).AddArg(IndexCharModify).ToString(),
	}
}
