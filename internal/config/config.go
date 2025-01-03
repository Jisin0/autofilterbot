// Package config contains types for the bot's global configuration.
package config

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/internal/model/message"
)

// Config contains custom values saved for the bot using the config panel.
type Config struct {
	BotId int64 `json:"_id,omitempty" bson:"_id,omitempty" `
	// Force Subscribe Channels.
	FsubChannels []model.FsubChannel `json:"fsub,omitempty" bson:"fsub,omitempty"`

	// Autofilter result settings

	MaxResults int `json:"max_results,omitempty" bson:"max_results,omitempty"`
	MaxPerPage int `json:"max_per_page,omitempty" bson:"max_per_page,omitempty"`
	MaxPages   int `json:"max_pages,omitempty" bson:"max_pages,omitempty"`

	// Custom Start Message
	StartText    string                          `json:"start_text,omitempty" bson:"start_text,omitempty"`
	StartButtons [][]button.InlineKeyboardButton `json:"start_button,omitempty" bson:"start_button,omitempty"`
	// Custom About Message
	AboutText    string                          `json:"about_text,omitempty" bson:"about_text,omitempty"`
	AboutButtons [][]button.InlineKeyboardButton `json:"about_button,omitempty" bson:"about_button,omitempty"`
	// Custom Help Message
	HelpText    string                          `json:"help_text,omitempty" bson:"help_text,omitempty"`
	HelpButtons [][]button.InlineKeyboardButton `json:"help_button,omitempty" bson:"help_button,omitempty"`
	// Custom Stats Message
	StatsText   string                          `json:"stats_text,omitempty" bson:"stats_text,omitempty"`
	StatsButton [][]button.InlineKeyboardButton `json:"stats_button,omitempty" bson:"stats_button,omitempty"`
	// Custom Privacy Message
	PrivacyText    string                          `json:"privacy_text,omitempty" bson:"privacy_text,omitempty"`
	PrivacyButtons [][]button.InlineKeyboardButton `json:"privacy_button,omitempty" bson:"privacy_button,omitempty"`
}

// GetStartMessage returns the custom start message if available or the default values.
// botUsername must be provided to create the add to group button.
func (c *Config) GetStartMessage(botUsername string) *message.Message {
	var (
		text    string
		buttons [][]button.InlineKeyboardButton
	)

	if c.StartText != "" {
		text = c.StartText
	} else {
		text = `<i><b>Hey there {mention} 👋</b></i>

🔥 I'm an awesome media <b>search</b> bot that can filter through millions of <b>files</b> in seconds 🗃️

Add me to a group or type go inline to start using me 👇`
	}

	if len(c.StartButtons) != 0 {
		buttons = c.StartButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "➕ Add Me To Your Group  ➕", Url: fmt.Sprintf("https://t.me/%s?startgroup=true&admin=delete_messages+pin_messages+invite_users+ban_users+promote_members", botUsername)}},
			{{Text: "About", CallbackData: "cmd:about"}, {Text: "Help", CallbackData: "cmd:help"}},
			{{Text: "Search Inline 🔎", SwitchInlineQueryCurrentChat: "", IsInline: true}},
		}
	}

	return &message.Message{
		Text:     text,
		Keyboard: buttons,
	}
}

func (c *Config) GetAboutMessage() *message.Message {
	var (
		text    string
		buttons [][]button.InlineKeyboardButton
	)

	if c.AboutText != "" {
		text = c.AboutText
	} else {
		text = `
○ Language : Go
○ Library : gotgbot
○ Database : {database}
○ Version : 0.1
`
	}

	if len(c.AboutButtons) != 0 {
		buttons = c.AboutButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "Source", Url: "https://github.com/Jisin0/autofilterbot"}, {Text: "Stats", CallbackData: "cmd:stats"}},
		}
	}

	return &message.Message{
		Text:     text,
		Keyboard: buttons,
	}
}

func (c *Config) GetHelpMessage() *message.Message {
	var (
		text    string
		buttons [][]button.InlineKeyboardButton
	)

	if c.HelpText != "" {
		text = c.HelpText
	} else {
		text = `
🖐️ Here's Two Ways you can Use me. . .

○ <b>Inline</b>: Just Start Typing my Username into any Chat and get Results On The Fly ✈️
○ <b>Groups</b>: Add me to your Group Chat and Just Send the Name of the File you Want ✍️

🤖 Other Commands:
/start - check if I'm alive
/about - learn a bit about me
/help - get this message
/stats - some number crushing
/privacy - what data I steal
/uinfo - get user data stored
/id - if you know u know

🍷 Exclusive Commands:
/broadcast - spam users with ads
/settings - customize me
/batch - bunch up messages
/genlink - link to single file
/index - gather up files
/delete - assassinate a file
/deleteall - massacre matching files
`
	}

	if len(c.HelpButtons) != 0 {
		buttons = c.HelpButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "<- Back", CallbackData: "cmd:start"}, {Text: "Privacy", CallbackData: "cmd:privacy"}},
		}
	}

	return &message.Message{
		Text:     text,
		Keyboard: buttons,
	}
}

func (c *Config) GetStatsMessage() *message.Message {
	var (
		text    string
		buttons [][]button.InlineKeyboardButton
	)

	if c.StatsText != "" {
		text = c.StatsText
	} else {
		text = `
╭ ▸ Users : <code>{users}</code> 
├ ▸ Files : <code>{files}</code>
├ ▸ Groups : <code>{groups}</code>
╰ ▸ Uptime : <code>{uptime}</code>
`
	}

	if len(c.StatsButton) != 0 {
		buttons = c.StatsButton
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "<- Back", CallbackData: "cmd:about"}},
		}
	}

	return &message.Message{
		Text:     text,
		Keyboard: buttons,
	}
}

func (c *Config) GetPrivacyMessage() *message.Message {
	var (
		text    string
		buttons [][]button.InlineKeyboardButton
	)

	if c.PrivacyText != "" {
		text = c.PrivacyText
	} else {
		text = `
<blockquote expandable>Privacy Policy 📜
This bot stores the publicly visible data of users that is required for the bot to operate.

The following data of a user could be saved:
‣ Id
‣ Name
‣ Username
‣ Join Requests

ℹ️ Use the /uinfo command with your user id to view data stored about you.
</blockquote>
`
	}

	if len(c.PrivacyButtons) != 0 {
		buttons = c.PrivacyButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "<- Back", CallbackData: "cmd:help"}},
		}
	}

	return &message.Message{
		Text:     text,
		Keyboard: buttons,
	}
}

func (c *Config) GetMaxResults() int {
	if c.MaxResults != 0 {
		return c.MaxResults
	}

	return 50
}

func (c *Config) GetMaxPerPage() int {
	if c.MaxPerPage != 0 {
		return c.MaxPerPage
	}

	return 10
}

func (c *Config) GetMaxPages() int {
	if c.MaxResults != 0 {
		return c.MaxResults
	}

	return 5
}
