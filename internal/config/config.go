// Package config contains types for the bot's global configuration.
package config

import (
	"github.com/Jisin0/autofilterbot/internal/model"
)

// Config contains the saved configs for the bot.
type Config struct {
	FsubChannels []model.FsubChannel
}
