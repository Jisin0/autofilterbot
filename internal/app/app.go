/*
Package app contains app type and helpers.
*/
package app

import (
	"time"

	"github.com/Jisin0/autofilterbot/internal/cache"
	"github.com/Jisin0/autofilterbot/internal/config"
	"github.com/Jisin0/autofilterbot/internal/database"
	"github.com/Jisin0/autofilterbot/pkg/autodelete"
	"github.com/Jisin0/autofilterbot/pkg/shortener"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

// App wraps various individual components of the app to orchestrate application processes.
type App struct {
	DB        database.Database
	Log       *zap.Logger
	StartTime time.Time
	Bot       *gotgbot.Bot
	Cache     *cache.Cache
	Config    *config.Config
	Admins    []int64

	AutoDelete *autodelete.Manager
	Shortener  *shortener.Shortener
}

func (a *App) GetDB() database.Database {
	return a.DB
}

func (a *App) GetLog() *zap.Logger {
	return a.Log
}

func (a *App) GetStartTime() time.Time {
	return a.StartTime
}

func (a *App) GetBot() *gotgbot.Bot {
	return a.Bot
}

func (a *App) GetCache() *cache.Cache {
	return a.Cache
}

func (a *App) GetConfig() *config.Config {
	return a.Config
}

func (a *App) GetAdmins() []int64 {
	return a.Admins
}

func (a *App) GetAutoDelete() *autodelete.Manager {
	return a.AutoDelete
}

func (a *App) GetShortener() *shortener.Shortener {
	return a.Shortener
}
