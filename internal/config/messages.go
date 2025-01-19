package config

import (
	"fmt"

	"github.com/Jisin0/autofilterbot/internal/button"
	"github.com/Jisin0/autofilterbot/internal/model/message"
)

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
