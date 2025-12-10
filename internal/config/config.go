// Package config contains types for the bot's global configuration.
package config

import (
	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/pkg/shortener"
)

// Config contains custom values saved for the bot using the config panel.
type Config struct {
	BotId int64 `json:"_id" bson:"_id" `

	// Autofilter result settings

	MaxResults int `json:"max_results,omitempty" bson:"max_results,omitempty"`
	MaxPerPage int `json:"max_per_page,omitempty" bson:"max_per_page,omitempty"`
	MaxPages   int `json:"max_pages,omitempty" bson:"max_pages,omitempty"`

	// Custom Start Message
	StartText    string                          `json:"start_text,omitempty" bson:"start_text,omitempty"`
	StartButtons [][]button.InlineKeyboardButton `json:"start_buttons,omitempty" bson:"start_buttons,omitempty"`
	// Custom About Message
	AboutText    string                          `json:"about_text,omitempty" bson:"about_text,omitempty"`
	AboutButtons [][]button.InlineKeyboardButton `json:"about_buttons,omitempty" bson:"about_buttons,omitempty"`
	// Custom Help Message
	HelpText    string                          `json:"help_text,omitempty" bson:"help_text,omitempty"`
	HelpButtons [][]button.InlineKeyboardButton `json:"help_buttons,omitempty" bson:"help_buttons,omitempty"`
	// Custom Stats Message
	StatsText    string                          `json:"stats_text,omitempty" bson:"stats_text,omitempty"`
	StatsButtons [][]button.InlineKeyboardButton `json:"stats_buttons,omitempty" bson:"stats_buttons,omitempty"`
	// Custom Privacy Message
	PrivacyText    string                          `json:"privacy_text,omitempty" bson:"privacy_text,omitempty"`
	PrivacyButtons [][]button.InlineKeyboardButton `json:"privacy_buttons,omitempty" bson:"privacy_buttons,omitempty"`

	// Force Subscribe Channels.
	FsubChannels []model.FsubChannel `json:"fsub,omitempty" bson:"fsub,omitempty"`
	// Fsub message text.
	FsubText string `json:"fsub_text,omitempty" bson:"fsub_text,omitempty"`
	// Html formatted file caption.
	FileCaption string `json:"file_caption,omitempty" bson:"file_caption,omitempty"`
	// File autodelete time in minutes.
	FileAutoDelete int `json:"file_autodel,omitempty" bson:"file_autodel,omitempty"`

	// Template to use for autofilter result message
	ResultTemplate string `json:"af_template,omitempty" bson:"af_template,omitempty"`
	// Message sent when no results are available.
	NoResultText string `json:"no_result_text,omitempty" bson:"no_result_text,omitempty"`
	// Template to use for result buttons
	ButtonTemplate string `json:"btn_template,omitempty" bson:"btn_template,omitempty"`
	// File Details Calbback Template.
	FileDetailsTemplate string `json:"fdetails_template,omitempty" bson:"fdetails_template,omitempty"`

	// Maximum number of message in a single batch.
	BatchSizeLimit int64 `json:"batch_size,omitempty" bson:"batch_size,omitempty"`

	// File size is shown in separate button if set
	SizeButton bool `json:"size_btn,omitempty" bson:"size_btn,omitempty"`

	Shortener shortener.Shortener `json:"shortener,omitempty" bson:"shortener,omitempty"`

	// Time in minutes after which result message should be deleted.
	AutodeleteTime int `json:"autodel_time,omitempty" bson:"autodel_time,omitempty"`

	// Index of collection to use to store files.
	FileCollectionIndex int `json:"collection_index,omitempty" bson:"collection_index,omitempty"`
	// Indicates wether the updater should be run to update file collection periodically.
	FileCollectionUpdater bool `json:"collection_updater,omitempty" bson:"collection_updater,omitempty"`

	// cached value from ToMap, updated using UpdateMap
	cachedMap map[string]any
}

func (c *Config) GetShortener() shortener.Shortener {
	return c.Shortener
}

func (c *Config) GetAutodeleteTime() int {
	return c.AutodeleteTime
}

func (c *Config) GetFileDetailsTemplate() string {
	if c.FileDetailsTemplate != "" {
		return c.FileDetailsTemplate
	}

	return `Name: {file_name}
Size: {file_size}
Type: {file_type}
Uploaded: {date}`
}

func (c *Config) GetFsubChannels() []model.FsubChannel {
	return c.FsubChannels
}

func (c *Config) GetFsubText() string {
	if c.FsubText != "" {
		return c.FsubText
	}

	return `<i><b>👋 Hᴇʏ ᴛʜᴇʀᴇ {mention}</b></i>
<i>Pʟᴇᴀsᴇ ᴊᴏɪɴ ᴍʏ ᴄʜᴀɴɴᴇʟs ғɪʀsᴛ ᴛᴏ ɢᴇᴛ ʏᴏᴜʀ ғɪʟᴇ</i>

<i>Cʟɪᴄᴋ ᴛʜᴇ <b>JOIN</b> ʙᴜᴛᴛᴏɴ ʙᴇʟᴏᴡ ᴀɴᴅ ᴛʜᴇɴ <b>RETRY</b> ᴛᴏ ɢᴇᴛ ʏᴏᴜʀ ғɪʟᴇ 👇</i>`
}

func (c *Config) GetFileCaption() string {
	if c.FileCaption != "" {
		return c.FileCaption
	}

	return "<i>{file_name}</i>\n\n<b>📂 File Size</b>: <code>{file_size}</code>\n{warn}"
}

func (c *Config) GetFileAutoDelete() int {
	return c.FileAutoDelete
}

func (c *Config) GetBatchSizeLimit() int64 {
	if c.BatchSizeLimit != 0 {
		return c.BatchSizeLimit
	}

	return 50
}

func (c *Config) GetFileCollectionIndex() int {
	return c.FileCollectionIndex
}

func (c *Config) GetFileCollectiionUpdater() bool {
	return c.FileCollectionUpdater
}
