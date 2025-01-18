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
	// Force Subscribe Channels.
	FsubChannels []model.FsubChannel `json:"fsub,omitempty" bson:"fsub,omitempty"`

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

	// Template to use for autofilter result message
	ResultTemplate string `json:"af_template,omitempty" bson:"af_template,omitempty"`
	// Message sent when no results are available.
	NoResultText string `json:"no_result_text,omitempty" bson:"no_result_text,omitempty"`
	// Template to use for result buttons
	ButtonTemplate string `json:"btn_template,omitempty" bson:"btn_template,omitempty"`
	// File Details Calbback Template.
	FileDetailsTemplate string `json:"fdetails_template,omitempty" bson:"fdetails_template,omitempty"`

	// File size is shown in seperate button if set
	SizeButton bool `json:"size_btn,omitempty" bson:"size_btn,omitempty"`

	Shortener shortener.Shortener `json:"shortener,omitempty" bson:"shortener,omitempty"`

	// Time in minutes after which message should be deleted.
	AutodeleteTime int `json:"autodel_time,omitempty" bson:"autodel_time,omitempty"`

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
