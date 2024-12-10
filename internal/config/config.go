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
}
