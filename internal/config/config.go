// Package config contains types for the bot's global configuration.
package config

import (
	"github.com/Jisin0/autofilterbot/internal/models"
)

// Config contains the saved configs for the bot.
type Config struct {
	FsubChannels []models.FsubChannel
}
