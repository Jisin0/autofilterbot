// Package config contains types for the bot's global configuration.
package config

import (
	"fmt"
	"runtime"

	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/model"
	"github.com/Jisin0/autofilterbot/internal/model/message"
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

	// File size is shown in seperate button if set
	SizeButton bool `json:"size_btn,omitempty" bson:"size_btn,omitempty"`

	Shortener shortener.Shortener

	// Time in minutes after which message should be deleted.
	AutodeleteTime int `json:"autodel_time,omitempty" bson:"autodel_time,omitempty"`
}

// GetStartMessage returns the custom start message if available or the default values.
// botUsername must be provided to create the add to group button.
func (c *Config) GetStartMessage(botUsername string) *message.Message {
	var (
		text    string
		buttons [][]button.InlineKeyboardButton
	)

	addToGroupUrl := fmt.Sprintf("https://t.me/%s?startgroup&admin=delete_messages+pin_messages+invite_users+ban_users+promote_members", botUsername)

	if c.StartText != "" {
		text = c.StartText
	} else {
		text = fmt.Sprintf(`<i><b>Hᴇʏ Tʜᴇʀᴇ <tg-spoiler>{mention}</tg-spoiler> 👅</b></i>
<blockquote><b>𝖠𝖼𝖼𝖾𝗌𝗌 𝖧𝗎𝗇𝖽𝗋𝖾𝖽𝗌 𝗈𝖿 𝖳𝗁𝗈𝗎𝗌𝖺𝗇𝖽𝗌 𝗈𝖿 𝖥𝗂𝗅𝖾𝗌 𝖤𝖺𝗌𝗂𝗅𝗒... <a href='%s'>𝖠𝖽𝖽 𝖬𝖾</a> 𝖳𝗈 𝖺 𝖦𝗋𝗈𝗎𝗉 𝗈𝗋 𝖳𝗒𝗉𝖾 𝖬𝗒 𝖴𝗌𝖾𝗋𝗇𝖺𝗆𝖾 𝗂𝗇 𝖠𝗇𝗒 𝖢𝗁𝖺𝗍 𝖳𝗈 𝖲𝗍𝖺𝗋𝗍 𝖲𝖾𝖺𝗋𝖼𝗁𝗂𝗇𝗀!!</b></blockquote>`, addToGroupUrl)
	}

	if len(c.StartButtons) != 0 {
		buttons = c.StartButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "ᴀʙᴏᴜᴛ", CallbackData: "cmd:about"}, {Text: "ʜᴇʟᴘ", CallbackData: "cmd:help"}},
			{{Text: "ᴘʀɪᴠᴀᴄʏ", CallbackData: "cmd:privacy"}, {Text: "sᴇᴀʀᴄʜ 🔎", SwitchInlineQueryCurrentChat: "", IsInline: true}},
			{{Text: "➕ ᴀᴅᴅ ᴍᴇ ᴛᴏ ʏᴏᴜʀ ɢʀᴏᴜᴘ  ➕", Url: addToGroupUrl}},
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
	runtime.Version()
	if c.AboutText != "" {
		text = c.AboutText
	} else {
		text = `
○ 𝖫𝖺𝗇𝗀𝗎𝖺𝗀𝖾  :  Go 1.22
○ 𝖫𝗂𝖻𝗋𝖺𝗋𝗒  :  <a href='https://github.com/PaulSonOfLars/gotgbot'>GoTgBot</a>
○ 𝖮𝖲  :  <code>{os}</code>
⛃ 𝖣𝖺𝗍𝖺𝖻𝖺𝗌𝖾  :  <code>{database}</code>
○ 𝖵𝖾𝗋𝗌𝗂𝗈𝗇  :  <code>0.1</code>
`
	}

	if len(c.AboutButtons) != 0 {
		buttons = c.AboutButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "« ʙᴀᴄᴋ", CallbackData: "cmd:start"}, {Text: "sᴛᴀᴛs", CallbackData: "cmd:stats"}},
			{{Text: "sᴏᴜʀᴄᴇ 🔗", Url: "https://github.com/Jisin0/autofilterbot"}},
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
<b>🖐️ 𝖧𝖾𝗋𝖾'𝗌 𝖳𝗐𝗈 𝖶𝖺𝗒𝗌 𝖸𝗈𝗎 𝖢𝖺𝗇 𝖴𝗌𝖾 𝖬𝖾. . .</b>

✈️ <b>𝖨𝗇𝗅𝗂𝗇𝖾</b> : <i>Just Start Typing my Username into any Chat and get Results On The Fly</i>
✍️ <b>𝖦𝗋𝗈𝗎𝗉</b> : <i>Add me to your Group Chat and Just Send the Name of the File you Want</i>

🤖 <b>User Commands:</b>
/start - check if I'm alive
/about - learn a bit about me
/help - get this message
/stats - some number crushing
/privacy - what data I steal
/uinfo - get user data stored
/id - if you know u know

🍷 <b>Exclusive Commands:</b>
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
			{{Text: "« ʙᴀᴄᴋ", CallbackData: "cmd:start"}, {Text: "ᴘʀɪᴠᴀᴄʏ", CallbackData: "cmd:privacy"}},
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
╭ ▸ 𝖴𝗌𝖾𝗋𝗌 : <code>{users}</code> 
├ ▸ 𝖥𝗂𝗅𝖾𝗌 : <code>{files}</code>
├ ▸ 𝖦𝗋𝗈𝗎𝗉𝗌 : <code>{groups}</code>
╰ ▸ 𝖴𝗉𝗍𝗂𝗆𝖾 : <code>{uptime}</code>
`
	}

	if len(c.StatsButtons) != 0 {
		buttons = c.StatsButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "« ʙᴀᴄᴋ", CallbackData: "cmd:about"}, button.CloseLocal()},
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
<blockquote expandable><b>Privacy Policy 📜</b>
<i>This bot stores the <b>publicly</b> visible data of users that is <b>required</b> for the bot to operate.

The following data of a user could be saved:
‣ Id
‣ Name
‣ Username
‣ Join Requests

ℹ️ Use the /uinfo command with your user id to view data stored about you.</i>
</blockquote>
`
	}

	if len(c.PrivacyButtons) != 0 {
		buttons = c.PrivacyButtons
	} else {
		buttons = [][]button.InlineKeyboardButton{
			{{Text: "« ʙᴀᴄᴋ", CallbackData: "cmd:help"}, button.CloseLocal()},
		}
	}

	return &message.Message{
		Text:     text,
		Keyboard: buttons,
	}
}

func (c *Config) GetShortener() shortener.Shortener {
	return c.Shortener
}

func (c *Config) GetAutodeleteTime() int {
	return c.AutodeleteTime
}
