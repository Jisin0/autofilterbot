// Package config contains types for the bot's global configuration.
package config

import (
	"github.com/Jisin0/autofilterbot/internal/model"
)

// Config contains the saved configs for the bot.
type Config struct {
	FsubChannels []model.FsubChannel `json:"fsub,omitempty"`

	// Autofilter result settings

	MaxResults int `json:"max_results,omitempty"`
	MaxPerPage int `json:"max_per_page,omitempty"`
	MaxPages   int `json:"max_pages,omitempty"`

	// Custom text values.

	StartText string `json:"start_text,omitempty"`
	AboutText string `json:"about_text,omitempty"`
	HelpText  string `json:"help_text,omitempty"`
	StatsText string `json:"stats_text,omitempty"`
	PrivacyText string `json:"privacy_text,omitempty"`
}

func (c *Config) GetStartText() string {
	if c.StartText != "" {
		return c.StartText
	}

	return `
<i><b>Hey there {mention} 👋</b></i>

🔥 I'm an awesome media <b>search</b> bot that can filter through millions of <b>files</b> in seconds 🗃️

Add me to a group or type go inline to start using me 👇
`
}

func (c *Config) GetAboutText() string {
	if c.AboutText != "" {
		return c.AboutText
	}

	return `
○ Language : Go
○ Library : gotgbot
○ Database : {database}
○ Version : 0.1
`
}

func (c *Config) GetHelpText() string {
	if c.HelpText != "" {
		return c.HelpText
	}

	return `
🖐️ Here's Two Ways you can Use me. . .

○ <b>Inline</b>: Just Start Typing my Username into any Chat and get Results On The Fly ✈️
○ <b>Groups</b>: Add me to your Group Chat and Just Send the Name of the File you Want ✍️

🤖 Other Commands:
/start - check if I'm alive
/about - learn a bit about me
/help - get this message
/stats - some number crushing
/privacy - what data I steal
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

func (c *Config) GetStatsText() string {
	if c.StatsText != "" {
		return c.StatsText
	}

	return `
╭ ▸ Users : <code>{users}</code> 
├ ▸ Files : <code>{files}</code>
├ ▸ Groups : <code>{groups}</code>
╰ ▸ Uptime : <code>{uptime}</code>
`
}

func (c *Config) GetPrivacyText() string {
	if c.PrivacyText != "" {
		return c.PrivacyText
	}

	return `
<blockquote expandable>Privacy Policy
This bot stores the publicly visible data of users for marketing, analytics and core functioning purposes.

The following data of a user could be saved:
‣ Id
‣ Name
‣ Username
‣ Join Requests

ℹ️ Use the /info command with your user id to view data stored about you.
</blockquote>
`
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
